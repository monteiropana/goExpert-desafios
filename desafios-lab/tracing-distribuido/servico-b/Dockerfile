# FROM golang:latest

# WORKDIR /app

# COPY go.mod ./

# COPY go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o /servico-b

# CMD chown -R www-data:www-data /app

# EXPOSE 8080

# ENTRYPOINT ["/servico-b"]

FROM golang:1.19 AS builder

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

RUN go build -o /main


FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /main /main

CMD chown -R www-data:www-data /app

EXPOSE 8080

ENTRYPOINT ["/main"]