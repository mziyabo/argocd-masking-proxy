FROM golang:1.18

WORKDIR /usr/src/proxy

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

EXPOSE 3003

ENTRYPOINT ["masking-proxy"]