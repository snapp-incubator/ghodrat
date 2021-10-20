FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN go mod download  && go build -o /bin/app ./cmd/ghodrat/main.go

FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /bin/

COPY --from=builder /bin/app .

LABEL org.opencontainers.image.source="https://github.com/snapp-incubator/ghodrat-%%COMMAND%%"

ENTRYPOINT ["/bin/app"]

CMD ["%%COMMAND%%", "--env=prod"]
