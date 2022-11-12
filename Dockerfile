FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build ./...
FROM scratch
COPY --from=builder /app/dartboard /app/dartboard
CMD ["/app/dartboard"]

