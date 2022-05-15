#   Copyright [2022] [mziyabo]
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

BINARY_NAME=masking-proxy

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

TAG=passthrough
IMAGE_NAME= masking-proxy-internal

build:
	go build -o ./bin/${BINARY_NAME} ./cmd/proxy

# TODO: use docker to run this- as certs are being generated in Dockerfile
build_and_run:
	go build -o ./bin/${BINARY_NAME} ./cmd/proxy
	cp -r ./certs ./bin
	cp ./proxy.conf.json ./bin
	./bin/${BINARY_NAME}

container:
	 docker build -t $(IMAGE_NAME):$(TAG) .

clean:
	go clean
	rm -rf ./bin/