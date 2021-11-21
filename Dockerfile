FROM chromedp/headless-shell:latest

RUN apt-get update && apt-get install dumb-init

COPY server /app/server

ENTRYPOINT ["dumb-init", "--"]

CMD ["/app/server"]