# WebhookReforge

A high-fidelity developer utility to re-sign and replay webhooks locally. Bypasses stale timestamp and signature validation issues during development.

![WebhookReforge UI](https://raw.githubusercontent.com/Freedomwithin/webhook-reforge/main/docs/screenshot.png)

## Provider Support

| Provider | Scheme | Tolerance |
|----------|--------|-----------|
| Stripe | `t=....,v1=....` HMAC-SHA256 hex | 5 minutes |
| Paddle | `ts=...;h1=....` HMAC-SHA256 hex | ~5 seconds |
| Shopify | `X-Shopify-Hmac-SHA256` HMAC-SHA256 base64 | No timestamp |

## Why WebhookReforge?

Many webhook providers (Stripe, Paddle, Shopify) use HMAC signatures to ensure payload integrity. These signatures often include a Unix timestamp to prevent replay attacks.
- **Paddle** has a strict **5-second** tolerance window.
- **Stripe** has a **5-minute** window.

This makes replaying saved payloads for testing nearly impossible—the timestamp becomes stale, the signature fails validation, and your app rejects the event.

**WebhookReforge** solves this by:
1. Intercepting or loading a webhook payload.
2. Swapping the stale timestamp with the current time.
3. Re-computing the HMAC-SHA256 signature using your developer secret.
4. Forwarding the re-signed request to your local application.

Your code sees a valid, fresh webhook every time.

## Features

- **Multi-Provider Support**: Built-in adapters for **Stripe**, **Paddle**, and **Shopify**.
- **Proxy Mode**: Sits between your webhook source (ngrok, etc.) and your app, re-signing events on the fly.
- **Replay Mode**: Fire saved JSON payloads from a file or the UI with perfect signatures.
- **Native Desktop UI**: A professional, cross-platform GUI built with Wails v2 and Vanilla JS.
- **Local-First**: Zero telemetry, zero cloud dependencies. Your secrets stay on your machine.

## Installation

Download the latest binary for your platform from the [Releases](https://github.com/Freedomwithin/webhook-reforge/releases) page.

### Building from Source

#### Prerequisites (Linux Mint/Ubuntu)
```bash
sudo apt update
sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.1-dev
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
```

#### Build
```bash
wails build -tags webkit2_41
```

## Usage

### GUI Mode
Simply run the binary to launch the Desktop UI. Paste your payload, select your provider, and hit **Fire**.

### CLI Mode

#### Proxy Mode
```bash
./webhook-reforge proxy --port 9000 --target http://localhost:3000/webhooks --secret whsec_xxx --provider stripe
```

#### Replay Mode
```bash
./webhook-reforge replay --file payment.json --target http://localhost:3000/webhooks --secret whsec_xxx --provider paddle
```

## License
MIT
