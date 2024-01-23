FROM golang:1.21 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux  go build -o /app/main .


FROM scratch

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/config.yaml /app/config.yaml
CMD ["/app/main"]
