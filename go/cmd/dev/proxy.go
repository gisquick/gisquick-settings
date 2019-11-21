package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Requires Go 1.12 to support reverse proxy for websocket connections

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	proxyURL, err := url.Parse("https://localhost")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	portPtr := flag.Int("port", 8002, "port number")
	flag.Parse()
	address := fmt.Sprintf(":%d", *portPtr)
	fmt.Printf("Starting proxy on http://localhost%s\n", address)
	log.Fatal(http.ListenAndServe(address, proxy))
}
