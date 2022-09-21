FROM golang:1.19-buster as build

WORKDIR /app

COPY .git go.mod go.sum *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    GIT_COMMIT=$(git rev-parse --short HEAD) \
    go build \
        -ldflags="-w -s -extldflags '-static' -X main.PackageCommitHash=$GIT_COMMIT" \
        -a -o grafana-xmpp-webhook *.go


FROM scratch

WORKDIR /
COPY --from=build /app/grafana-xmpp-webhook /grafana-xmpp-webhook
EXPOSE 3033
USER 1000:1000

# We can't pass more arguments since it requires /bin/sh
CMD ["/grafana-xmpp-webhook"]
