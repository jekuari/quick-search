FROM golang AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN CG0_ENABLED=0 GOOS=linux go build -o main
RUN apt-get update && apt-get install -y ca-certificates

FROM ubuntu:latest

COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 

COPY html ./html

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
