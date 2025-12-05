FROM golang:1.25-alpine3.22 AS go-builder
ENV GOPATH=""

RUN apk add git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD ./ .
RUN go build -o ./solar-toolkit-daemon ./cmd/daemon
RUN go build -o ./solar-toolkit-gateway ./cmd/gateway

FROM alpine:3.21

COPY gateway/sql/migrations /app/migrations
COPY --from=go-builder /app/solar-toolkit-gateway /app/solar-toolkit-gateway
COPY --from=go-builder /app/solar-toolkit-daemon /app/solar-toolkit-daemon
COPY --from=go-builder /root/go/bin/migrate /bin/migrate

ENTRYPOINT ["/app/solar-toolkit-gateway"]
