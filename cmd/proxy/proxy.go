package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/fatih/color"

	masking "github.com/mziyabo/masking-proxy/cmd/masking"
	schema "github.com/mziyabo/masking-proxy/shared/schema"
)

var config schema.ProxyConfig

// TODO: load from file
// Loads proxy configuration
func loadConfig() {

	fmt.Println("Loading config")

	config.ApiServerUrl, _ = url.Parse("http://127.0.0.1:8001")
	config.Host = "127.0.0.1"
	config.Port = 3003
}

// TODO: add protocol, assuming http for now
// Start Proxy
func Start() {

	loadConfig()

	addr := strings.Join([]string{config.Host, fmt.Sprint(config.Port)}, ":")
	proxy := http.HandlerFunc(Handler)

	color.Cyan("Listening at: %s\n", addr)

	err := http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Panicf("Failed to listen on address: %s", addr)
	}

}

// Proxy http handler
func Handler(rw http.ResponseWriter, req *http.Request) {

	apiServerUrl := config.ApiServerUrl

	req.Host = apiServerUrl.Host
	req.URL.Host = apiServerUrl.Host
	req.URL.Scheme = apiServerUrl.Scheme
	req.RequestURI = ""

	addr, _, _ := net.SplitHostPort(req.RemoteAddr)

	req.Header.Set("X-Forwarded-For", addr)
	resp, err := http.DefaultClient.Do(req)

	fmt.Printf("ResponseWriter Type: %s\n", reflect.TypeOf(rw))

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
	defer resp.Body.Close()
	var b []byte

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
	default:
		b, _ = io.ReadAll(resp.Body)
	}

	if err != nil {
		log.Panic(err)
	}

	// Start masking work here
	// Do some conversions and parsing first...
	body := string(b)

	body = masking.Mask(body)
	back := []byte(body)

	go log.Println(rw.Header())

	// Bring it back
	rw.Write(back)
}
