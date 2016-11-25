package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var publicPort = flag.String("public-port", "9999", "Public port on which to serve proxied content")
var internalPort = flag.String("internal-port", "9998", "Internal port on which your reverse SSH tunnel is listening")
var hostname = flag.String("hostname", "localhost", "Value of the Host http header for which your vhost is configured")
var scheme = flag.String("scheme", "http", "Use http or https locally")
var insecure = flag.Bool("insecure", false, "If your local server uses a self-signed cert, set this to true")

var requestLogger *log.Logger
var responseLogger *log.Logger

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	responseLog := resp.Status

	if 301 == resp.StatusCode || 302 == resp.StatusCode {
		responseLog = responseLog + " '" + resp.Request.URL.EscapedPath() + "'"
	}

	responseLogger.Println(responseLog)

	return resp, nil
}

func main() {
	requestLogger = log.New(os.Stdout, "[proxy]\t<-\t", log.LstdFlags)
	responseLogger = log.New(os.Stdout, "[proxy]\t->\t", log.LstdFlags)

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestLogger.Println("[" + r.RemoteAddr + "]\t" + r.Method + " " + r.URL.EscapedPath())

		director := func(req *http.Request) {
			req.URL.Scheme = *scheme
			req.URL.Host = "127.0.0.1:" + *internalPort
			req.Host = *hostname
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		}

		proxy := &httputil.ReverseProxy{Director: director, Transport: &transport{tr}}
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Listening for external requests on port " + *publicPort)
	fmt.Println("Forwarding traffic to " + *scheme + "://" + *hostname + "/ on port " + *internalPort)
	requestLogger.Fatal(http.ListenAndServe(":"+*publicPort, nil))
}
