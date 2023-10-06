FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o eureak .

EXPOSE 8080

CMD ["/app/eureka", "-prefork=false", "-dev=false", "-port=:8080"]

