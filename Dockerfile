FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o web ./cmd/web/*.go

FROM alpine:latest
ENV APP_ENV=docker
ENV HOST=0.0.0.0
ENV PORT=8080
WORKDIR /root/
COPY --from=builder /app/web .
CMD [ "./web" ]
EXPOSE 8080