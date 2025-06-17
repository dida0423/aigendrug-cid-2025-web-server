FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

# Set default environment variables
ENV MAIN_DB_HOST=db \
    MAIN_DB_PORT=5432 \
    MAIN_DB_NAME=mydb \
    MAIN_DB_USER=diane \
    MAIN_DB_SCHEMA=aigendrug

CMD ["./main"]