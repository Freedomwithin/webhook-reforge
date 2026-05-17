# WebhookForge

A developer utility to re-sign and replay webhooks locally. Bypasses stale timestamp and signature validation issues during development.

## Phase 1 — MVP (Stripe Only)

### Features
- **Proxy Mode**: Acts as a middleware that re-signs incoming Stripe webhooks with a fresh timestamp and forwards them to your local application.
- **Replay Mode**: Takes a saved JSON payload, signs it with a current timestamp, and fires it at your target URL.

### Installation
```bash
go build -o webhookforge ./cmd/webhookforge
```

### Usage

#### Proxy Mode
```bash
./webhookforge proxy --port 9000 --target http://localhost:3000/webhooks --secret whsec_your_secret
```

#### Replay Mode
```bash
./webhookforge replay --file payment_intent.json --target http://localhost:3000/webhooks --secret whsec_your_secret
```

## License
MIT
