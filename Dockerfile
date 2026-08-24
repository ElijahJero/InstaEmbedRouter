FROM golang:1.24

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /usr/local/bin/InstagramEmbedResolver .

EXPOSE 8080

ENV PROXY_PORT=8080
ENV PROMETHEUS_URL=

CMD ["/usr/local/bin/InstagramEmbedResolver"]
