FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o eureak-search .

EXPOSE 8080

CMD ["/app/eureka-search", "-prefork=false", "-dev=false", "-port=:8080"]

