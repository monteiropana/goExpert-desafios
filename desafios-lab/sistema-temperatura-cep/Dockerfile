
FROM golang:1.20
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server
ENTRYPOINT [ "/app/server" ] 

#FROM golang:1.20 AS build 
#FROM scratch 
#COPY --from=build app/server /server
#ENTRYPOINT [ "/server" ]