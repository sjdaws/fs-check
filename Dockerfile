FROM golang:1.20 as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine

COPY --from=builder /app/fs-check /app/fs-check
RUN /sbin/apk add --no-cache findutils

CMD ["/app/fs-check"]
