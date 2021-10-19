FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN go mod download  && go build -o /bin/app ./cmd/ghodrat/main.go

FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder ./static/audio.ogg /bin/ghodrat

COPY --from=builder /bin/app /bin/ghodrat

LABEL org.opencontainers.image.source="https://github.com/snapp-incubator/ghodrat-%%COMMAND%%"

ENTRYPOINT ["/bin/app"]

CMD ["%%COMMAND%%", "--env=prod"]
