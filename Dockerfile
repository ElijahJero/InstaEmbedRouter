FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/InstagramEmbedResolver .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/InstagramEmbedResolver /app/InstagramEmbedResolver
COPY templates /app/templates
COPY resolvers.json /app/resolvers.json

EXPOSE 8080

ENV PROXY_PORT=8080

CMD ["/app/InstagramEmbedResolver"]
