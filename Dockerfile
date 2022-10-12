FROM golang:1.18.4
WORKDIR /live

COPY go.mod go.sum /live/
RUN go mod download

COPY ./ /live/
RUN GOOS=linux GOARCH=amd64 go build -o qnlive ./app/live/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates && apk update && apk add tzdata

WORKDIR /root/
COPY --from=0 /live/qnlive .

CMD ["./qnlive", "-f", "/etc/qnlive/qnlive.yaml"]