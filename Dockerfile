FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .                

RUN go build -o favorite-service-app ./cmd/favoriteApp/main.go

EXPOSE 8080

CMD ["./favorite-service-app"]