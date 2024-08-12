FROM golang:1.21.1

WORKDIR /usr/local/src

COPY ./ ./

RUN go mod tidy
RUN go build -o ./service_auth_start ./service_auth/cmd/service_auth/main.go

CMD ["./service_auth_start"]
