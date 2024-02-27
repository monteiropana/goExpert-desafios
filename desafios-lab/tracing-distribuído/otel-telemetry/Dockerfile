#FROM golang:latest
#WORKDIR /go/src
#ENV PATH="/go/bin:${PATH}"
#ENV GO111MODULE=on
#ENV CGO_ENABLED=1
#CMD ["tail", "-f", "/dev/null"]


FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
