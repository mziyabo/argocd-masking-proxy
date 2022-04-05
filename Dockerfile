FROM golang:1.18 as builder

WORKDIR /usr/src/proxy

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

COPY certs /usr/src/proxy/certs/

RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -v -o . ./...

FROM alpine:latest

WORKDIR /usr/src/proxy

COPY --from=builder /usr/src/proxy/masking-proxy ./
COPY --from=builder /usr/src/proxy/proxy.conf.json ./
COPY --from=builder /usr/src/proxy/certs ./certs

EXPOSE 3003

ENTRYPOINT ["./masking-proxy"]