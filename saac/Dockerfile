FROM  golang:1.16-alpine AS builder

WORKDIR /app

# Copy dependencies first to improve caching.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY saac/ . 

RUN go build -o saac


FROM alpine:latest
CMD ["/saac", "--access-node-addr", "localnet_access_2_1:9000"]
COPY --from=builder /app/saac /saac
