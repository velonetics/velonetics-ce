FROM golang:1.25.3-alpine AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /mock-graphql .

FROM alpine:3.23
COPY --from=builder /mock-graphql /usr/local/bin/mock-graphql
EXPOSE 4000
ENTRYPOINT ["/usr/local/bin/mock-graphql"]
