FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./mobilerecharge .

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/mobilerecharge .
COPY static/ ./static/
EXPOSE 8080
ENTRYPOINT ["./mobilerecharge"]