FROM golang:1.23 AS builder

WORKDIR /BookingSvc

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o booking-service ./cmd/main.go

#FROM alpine:latest



#COPY --from=builder /BookingSvc/booking-service /BookingSvc/booking-service

CMD ["./booking-service"]