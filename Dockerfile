FROM golang:1.18-alpine as builder

ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo building for "$TARGETPLATFORM"

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY ./ ./
RUN GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build -o qnlive ./app/live/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates && apk update && apk add tzdata

WORKDIR /root/
COPY --from=builder /workspace/qnlive .
COPY deploy/wait .
COPY deploy/qnlive-entrypoint.sh .
COPY deploy/qnlive.yaml /etc/qnlive.yaml
RUN chmod +x ./qnlive-entrypoint.sh && chmod +x ./wait

CMD ./wait && ./qnlive-entrypoint.sh