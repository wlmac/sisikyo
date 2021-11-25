FROM golang:1.17-alpine

WORKDIR /sisikyo
COPY . .

RUN go build -o /server server/cmd/main.go
EXPOSE 8080
CMD [ "/server", "--port=8080", "--api-timeout=5s", "--cache=12h" ]

