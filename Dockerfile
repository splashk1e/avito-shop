FROM golang:1.22.6
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main ./cmd/main.go -secret asdjqhuiqw43289zxcMASas -config ./config/config.yaml
CMD ["/app/main"]
