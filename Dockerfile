FROM golang AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN CG0_ENABLED=0 GOOS=linux go build -o main

FROM ubuntu:latest

COPY --from=builder /app/main .
COPY html ./html

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
