FROM golang:1.18 as builder

WORKDIR /usr/src/proxy

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

COPY certs /usr/src/proxy/certs/
RUN apt-get update -y \
 && apt-get install -y golang-cfssl \
 && cfssl gencert -initca ./certs/ca/ca-csr.json | cfssljson -bare  ./certs/ca/ca \
 && cfssl gencert -ca=./certs/ca/ca.pem -ca-key=./certs/ca/ca-key.pem ./certs/masking-proxy/config.json | cfssljson -bare ./certs/masking-proxy/masking-proxy

RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -v -o masking-proxy ./cmd/proxy

FROM alpine:latest

WORKDIR /usr/src/proxy

COPY --from=builder /usr/src/proxy/masking-proxy ./
COPY --from=builder /usr/src/proxy/proxy.conf.json ./
COPY --from=builder /usr/src/proxy/certs ./certs

EXPOSE 3003

ENTRYPOINT ["./masking-proxy"]