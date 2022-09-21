all: grafana-xmpp-webhook

grafana-xmpp-webhook: *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags='-w -s -extldflags "-static" -X main.PackageCommitHash=$(shell git rev-parse --short HEAD)' \
	-a -o grafana-xmpp-webhook *.go
