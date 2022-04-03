package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	color "github.com/fatih/color"

	masking "github.com/mziyabo/masking-proxy/cmd/masking"
	"github.com/mziyabo/masking-proxy/shared"
)

var config shared.ProxyConfig

func init() {
	config = shared.Config
}

// TODO: add protocol, assuming http for now
// Start Proxy
func Start() {

	addr := strings.Join([]string{config.Host, fmt.Sprint(config.Port)}, ":")
	proxy := http.HandlerFunc(Handler)

	var listenErr error

	color.Cyan("Listening at: %s\n", addr)

	// TODO: fix tls
	if config.TLSConfig.Enabled {
		listenErr = http.ListenAndServeTLS(addr, config.TLSConfig.Cert, config.TLSConfig.Key, proxy)
	} else {
		listenErr = http.ListenAndServe(addr, proxy)
	}

	if listenErr != nil {
		_ = fmt.Errorf("Failed to listen on address: %s", addr)
		log.Panic(listenErr)
	}
}

// Proxy http handler
func Handler(rw http.ResponseWriter, req *http.Request) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	apiServerUrl := config.ApiServerUrl

	req.Host = apiServerUrl.Host
	req.URL.Host = apiServerUrl.Host
	req.URL.Scheme = apiServerUrl.Scheme
	req.RequestURI = ""

	// TODO: should be in config.go!!
	// Read token in from /var/run
	dat, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		req.Header.Add("Authorization", ("Bearer " + string(dat)))
	}

	addr, _, _ := net.SplitHostPort(req.RemoteAddr)

	req.Header.Set("X-Forwarded-For", addr)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}

	// Return response content
	rw.WriteHeader(resp.StatusCode)
	rw.Header().Del("Content-Length")

	defer resp.Body.Close()
	var data []byte

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	// Start masking work here
	// Do some conversions and parsing first...
	body := masking.Mask(data)

	// TODO: fix logging levels
	go log.Println(rw.Header())
	go log.Println(string(body))
	go log.Println("---")

	// Bring it back
	rw.Write(body)
}
