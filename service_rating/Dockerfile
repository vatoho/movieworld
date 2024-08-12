FROM golang:1.21.1

WORKDIR /usr/local/src

COPY ./ ./

RUN go mod tidy
RUN go build -o ./service_rating_start ./service_rating/cmd/service_rating/main.go

CMD ["./service_rating_start"]
