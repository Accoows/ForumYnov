FROM golang:1.24.1-bookworm AS builder

RUN apt-get update && apt-get install -y gcc libsqlite3-dev sqlite3

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o forumynov .
RUN ls -l /app

FROM golang:1.24.1-bookworm
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0

COPY --from=builder /app/forumynov .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/database ./database
COPY --from=builder /app/handlers ./handlers
COPY --from=builder /app/models ./models

RUN chmod +x /app/forumynov

EXPOSE 8080

CMD [ "/app/forumynov" ]