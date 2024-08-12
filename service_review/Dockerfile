FROM golang:1.21.1

WORKDIR /usr/local/src

COPY ./ ./

RUN go mod tidy
RUN go build -o ./service_review_start ./service_review/cmd/service_review/main.go

CMD ["./service_review_start"]
