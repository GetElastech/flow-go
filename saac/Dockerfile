FROM  golang:1.16-alpine AS builder

WORKDIR /app

# Copy dependencies first to improve caching.
COPY go.mod .
COPY go.sum .

COPY saac/ . 

RUN go build -o saac


FROM alpine:latest
CMD ["/saac"]
COPY --from=builder /app/saac /saac
