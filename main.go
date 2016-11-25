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

func main() {
	log := log.New(os.Stdout, "[proxy]\t", log.LstdFlags)

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[" + r.RemoteAddr + "]\t" + r.Method + " " + r.URL.EscapedPath())

		director := func(req *http.Request) {
			req.URL.Scheme = *scheme
			req.URL.Host = "127.0.0.1:" + *internalPort
			req.Host = *hostname
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		}

		proxy := &httputil.ReverseProxy{Director: director, Transport: tr}
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Listening for external requests on port " + *publicPort)
	fmt.Println("Forwarding traffic to " + *scheme + "://" + *hostname + "/ on port " + *internalPort)
	log.Fatal(http.ListenAndServe(":"+*publicPort, nil))
}
