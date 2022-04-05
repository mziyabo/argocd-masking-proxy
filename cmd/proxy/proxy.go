package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	req.Close = true

	req.Host = apiServerUrl.Host
	req.URL.Host = apiServerUrl.Host
	req.URL.Scheme = apiServerUrl.Scheme
	req.RequestURI = ""

	// DEBUG
	fmt.Println(req.URL)
	fmt.Println(req.Context())

	// TODO: should be in config.go!!!
	// Read token in from /var/run
	dat, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		req.Header.Add("Authorization", ("Bearer " + string(dat)))
	}

	addr, _, _ := net.SplitHostPort(req.RemoteAddr)

	client := http.Client{Timeout: 300 * time.Second}

	req.Header.Set("X-Forwarded-For", addr)
	resp, err := client.Do(req)

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

	if resp.Header.Get("Content-Length") != "" {
		contentLength, strconvErr := strconv.Atoi(resp.Header.Get("Content-Length"))
		fmt.Println(contentLength)
		if strconvErr != nil {
			fmt.Println(strconvErr)
		}

	}
	var data []byte
	reader := bytes.NewBuffer(data)
	var readerErr error

	count, readerErr := io.Copy(reader, resp.Body)
	data = reader.Bytes()

	if readerErr != nil {
		log.Println(resp.Status)
		log.Println(resp.Request.URL)
		log.Println(readerErr)
	}

	// Start masking work here
	// Do some conversions and parsing first...
	body := masking.Mask(data)
	//rw.Header().Add("Content-Length", fmt.Sprint(len(data)))

	log.Printf("DATA-PRE: %s, DATA-POST: %s", fmt.Sprint(count), fmt.Sprint(len(body)))
	//l := len(data)

	// TODO: fix logging levels
	log.Println(rw.Header())

	// Bring it back
	rw.Write(body)
}
