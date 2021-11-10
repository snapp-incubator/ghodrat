FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN go mod download  && \
    mkdir /bin/ghodrat && \
    go build -o /bin/ghodrat/app ./cmd/ghodrat/main.go && \
    cp ./static/audio.ogg /bin/ghodrat && \
    cp ./static/video.ivf /bin/ghodrat

FROM alpine:latest

RUN apk add --no-cache libc6-compat && mkdir /bin/ghodrat

COPY --from=builder /bin/ghodrat/audio.ogg /bin/ghodrat
COPY --from=builder /bin/ghodrat/video.ivf /bin/ghodrat

COPY --from=builder /bin/ghodrat/app /bin/ghodrat

LABEL org.opencontainers.image.source="https://github.com/snapp-incubator/ghodrat-%%COMMAND%%"

ENTRYPOINT ["/bin/ghodrat/app"]

CMD ["%%COMMAND%%"]
