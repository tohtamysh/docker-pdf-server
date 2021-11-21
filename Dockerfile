FROM golang:1.17-alpine as build

WORKDIR /src

COPY * /src/

RUN go mod download

RUN CGO_ENABLED=0 go build -o server server.go

FROM chromedp/headless-shell:latest

RUN apt-get update && apt-get install dumb-init

COPY --from=build /src/server /server

EXPOSE 8000

ENTRYPOINT ["dumb-init", "--"]

CMD ["/server"]