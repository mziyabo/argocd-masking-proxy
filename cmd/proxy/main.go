/*
   Copyright [2022] [mziyabo]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	masking "github.com/mziyabo/masking-proxy/cmd/masking"
	"github.com/mziyabo/masking-proxy/cmd/shared"
)

var config shared.ProxyConfig

func main() {
	Start()
}

// Start begins listening and serving at configured address
func Start() {

	config = shared.Config

	proxyAddr := strings.Join([]string{config.Host, fmt.Sprint(config.Port)}, ":")
	proxy := http.HandlerFunc(Handler)

	log.Printf("Listening at: %s\n", proxyAddr)

	var listenErr error
	if config.TLSConfig.Enabled {
		listenErr = http.ListenAndServeTLS(proxyAddr, config.TLSConfig.Cert, config.TLSConfig.Key, proxy)
	} else {
		listenErr = http.ListenAndServe(proxyAddr, proxy)
	}

	if listenErr != nil {
		_ = fmt.Errorf("failed to listen on address: %s", proxyAddr)
		log.Panic(listenErr)
	}
}

// Proxy http handler
func Handler(rw http.ResponseWriter, req *http.Request) {

	req.Close = true
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	apiURL := config.ApiURL

	req.Host = apiURL.Host
	req.URL.Host = apiURL.Host
	req.URL.Scheme = apiURL.Scheme
	req.RequestURI = ""

	remoteAddr, _, _ := net.SplitHostPort(req.RemoteAddr)

	req.Header.Set("X-Forwarded-For", remoteAddr)
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", config.Token}, " "))

	client := http.Client{Timeout: 300 * time.Second}
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

	defer resp.Body.Close()

	rw.WriteHeader(resp.StatusCode)

	var data []byte
	reader := bytes.NewBuffer(data)

	_, readerErr := io.Copy(reader, resp.Body)
	data = reader.Bytes()

	if readerErr != nil {
		_ = fmt.Errorf("ReaderError: %s", readerErr)
	}

	// Start masking work here
	// Do some conversions and parsing first...
	body := masking.Mask(data)

	rw.Write(body)
}
