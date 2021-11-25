FROM golang:1.17-alpine AS prod

WORKDIR /sisikyo
COPY . .

RUN go build -o /server server/cmd/main.go
EXPOSE 8080
ENV GIN_MODE=release
CMD [ "/server", "--port=8080", "--api-timeout=5s", "--cache=12h", "-o-id=$OAUTH_ID", "-o-secret=$OAUTH_SECRET", "--db-driver=$DB_DRIVER", "--db-source=$DB_SOURCE", "--db-timeout=5s" ]

FROM golang:1.17-alpine AS dev

WORKDIR /sisikyo
COPY . .

RUN go build -tags debug -o /server server/cmd/main.go
EXPOSE 8080
ENV GIN_MODE=debug
CMD [ "/server", "--port=8080", "--api-timeout=5s", "--cache=12h", "-o-id=$OAUTH_ID", "-o-secret=$OAUTH_SECRET", "--db-driver=$DB_DRIVER", "--db-source=$DB_SOURCE", "--db-timeout=5s" ]

