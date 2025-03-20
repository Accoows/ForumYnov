FROM golang:1.24.1-bookworm AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forumynov .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/forumynov .
COPY --from=builder /app/templates /app/templates

RUN chmod +x /app/forumynov

EXPOSE 8080

CMD [ "./forumynov" ]