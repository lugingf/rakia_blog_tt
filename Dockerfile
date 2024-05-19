FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o blog_tt .


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/blog_tt .

EXPOSE 8080

CMD ["./blog_tt"]


