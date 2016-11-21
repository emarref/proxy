package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
)

var publicPort = flag.String("public-port", "9999", "Public port on which to serve proxied content")
var internalPort = flag.String("internal-port", "9998", "Internal port on which your reverse SSH tunnel is listening")
var hostname = flag.String("hostname", "localhost", "Value of the Host http header for which your vhost is configured")
var scheme = flag.String("scheme", "http", "Use http or https locally")
var insecure = flag.Bool("insecure", false, "If your local server uses a self-signed cert, set this to true")

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	log.Fatal(http.ListenAndServe(":"+*publicPort, nil))
}
