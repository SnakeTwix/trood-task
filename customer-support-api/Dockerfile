FROM golang:1.24.1-alpine AS base
WORKDIR /app

FROM base AS prod

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -o ./bin ./main.go

CMD ["./bin"]
