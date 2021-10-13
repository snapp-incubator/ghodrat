FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN pwd

RUN ls -lah

RUN go mod download && make generate

RUN go build -o /bin/app ./cmd/root.go

FROM alpine:latest

RUN apk add --no-cache libc6-compat 

WORKDIR /bin/

COPY --from=builder /bin/app .

LABEL org.opencontainers.image.source="https://github.com/snapp-incubator/ghodrat-%%COMMAND%%"

ENTRYPOINT ["/bin/app"]

CMD ["%%COMMAND%%", "--env=prod"]
