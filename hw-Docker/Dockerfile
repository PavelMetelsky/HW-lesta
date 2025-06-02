FROM golang:1.21.5 as builder

WORKDIR /go-app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd /go-app/cmd
COPY internal /go-app/internal
COPY docs /go-app/docs

RUN cd /go-app/cmd/demo/ && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/web ./main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

ENV API_PORT="${API_PORT}" \
    DB_URL="${DB_URL}"

COPY --from=builder /bin/web /bin/web

RUN chmod +x /bin/web

CMD ["/bin/sh", "-c", "/bin/web"]
