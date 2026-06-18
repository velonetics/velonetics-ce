FROM golang:1.25.11-bookworm AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /mock-webhook .

FROM alpine:3.23
COPY --from=builder /mock-webhook /usr/local/bin/mock-webhook
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/mock-webhook"]
