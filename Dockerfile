FROM golang:1.19-alpine3.15 AS builder

ADD . /app
RUN go mod download \
    && go build -o /build/server /app/cmd/app/main.go

FROM alpine:3.15
EXPOSE 80
COPY --from builder /build/ /
CMD ["/server", "--host", "0.0.0.0", "--port", "80"]