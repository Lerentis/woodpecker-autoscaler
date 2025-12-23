FROM golang:1.25 AS build

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o woodpecker-autoscaler ./cmd/

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build --chown=65534:65534 /app/woodpecker-autoscaler /usr/local/bin/woodpecker-autoscaler

USER nobody

ENTRYPOINT ["/usr/local/bin/woodpecker-autoscaler"]
