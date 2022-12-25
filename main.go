package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	apiBaseUrl string
}

func NewProxy(apiBaseUrl string) *Proxy { return &Proxy{apiBaseUrl: apiBaseUrl} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	req, err = http.NewRequest(r.Method, p.apiBaseUrl+r.RequestURI, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

}

func main() {
	proxyUrl := flag.String("proxy", "https://github.com/", "proxy url")
	localPort := flag.String("local", "12345", "local port")
	flag.Parse()

	proxy := NewProxy(*proxyUrl)
	fmt.Println("==============================")
	fmt.Println("Proxy Server started")
	err := http.ListenAndServe(":"+*localPort, proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
