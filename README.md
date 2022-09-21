# Grafana XMPP Webhook

[![Publish Docker image](https://github.com/h5vx/grafana-xmpp-webhook/actions/workflows/docker-push.yaml/badge.svg)](https://github.com/h5vx/grafana-xmpp-webhook/actions/workflows/docker-push.yaml)

## Why do I need it?
You probably need this, if you want to send alerts from Grafana to XMPP MUC room

## Features
- Send alerts to XMPP MUC room (it's the only option for now)
- Alerts templating (using go-templates)
- Single static binary, ~4MB compressed image size

## Configuration
See example: [config.toml](./config.toml)

### docker-compose.yml
```yaml
version: '3.3'

networks:
  monitoring:
    driver: bridge

services:
  xmpp-webhook:
    image: h5vx/grafana-xmpp-webhook:latest
    container_name: xmpp-webhook
    volumes:
      - ./grafana-xmpp-webhook.toml:/config.toml
    restart: unless-stopped
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
    # ...
    networks:
      - monitoring
```

That way, you can [add contact point](https://grafana.com/docs/grafana/latest/alerting/contact-points/create-contact-point/) in Grafana as **Type:** Webhook, **Url:** `http://xmpp-webhook/alert`

## Contribution
Issues, feature requests and other contributions are welcome!
