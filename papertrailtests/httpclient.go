package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Following proxy setup here https://jannewmarch.gitbooks.io/network-programming-with-go-golang-/http/proxy_handling.html
func NewHTTPClient() (*http.Client, error) {

	proxyURL, err := getProxyURL()
	if err != nil {
		// Access without proxy
		// logging.Log(logging.INFO, fmt.Sprintln(err))
		return &http.Client{Timeout: 30 * time.Second}, nil // Send err = nil so that we setup client without http_proxy
	}

	// transport := &http.Transport{
	// 	Proxy:           http.ProxyURL(proxyURL),
	// 	IdleConnTimeout: 120 * time.Second}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL)}
	return &http.Client{Transport: transport, Timeout: 30 * time.Second}, nil
}

func getProxyURL() (*url.URL, error) {

	// Get the proxy server setting
	httpProxy := os.Getenv("http_proxy")

	if httpProxy == "" {
		return nil, errors.New("No http_proxy env variable")
	}
	// logging.Log(logging.INFO, fmt.Sprintln("Using http_proxy = ", httpProxy))

	return url.Parse(httpProxy)
}

// // Following proxy setup here https://jannewmarch.gitbooks.io/network-programming-with-go-golang-/http/proxy_handling.html
// func newHTTPClient() (*http.Client, error) {
//
// 	if GproxyString == "" {
// 		client := &http.Client{Timeout: 30 * time.Second}
// 		return client, nil
// 	} else {
// 		proxyURL, err := url.Parse(GproxyString)
// 		if err != nil {
// 			logging.Log(logging.WARNING, fmt.Sprintf("Bad proxy URL: %s", proxyURL.String()))
// 		}
//
// 		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
// 		client := &http.Client{Transport: transport, Timeout: 30 * time.Second}
//
// 		return client, err
// 	}
// }
