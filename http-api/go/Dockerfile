FROM golang:1.23-alpine as builder

WORKDIR /go/src

RUN apk add --no-cache tzdata
ENV TZ=UTC

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go mod download

RUN go build -trimpath -ldflags="-s -w" -o /go/bin/api .

FROM scratch

WORKDIR /opt/bin

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /go/bin /opt/bin

ENV TZ=UTC
ENV USER=1000
EXPOSE 8080

CMD ["/opt/bin/api"]
