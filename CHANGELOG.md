# WebhookReforge Changelog

## v0.3.0 — Desktop UI
- Native desktop app via Wails v2
- Payload editor with JSON validation and formatting
- Live event log with status codes and response bodies
- Proxy mode with non-blocking goroutine
- Stripe, Paddle, Shopify provider switching in UI

## v0.2.0 — Multi-Provider CLI
- Paddle provider (ts=...;h1=... scheme, 5-second window)
- Shopify provider (raw body, base64 HMAC-SHA256)
- --provider flag on proxy and replay commands
- Extracted getProvider helper

## v0.1.0 — Stripe MVP
- Stripe webhook re-signing proxy
- File-based replay engine
- Single binary, zero external dependencies
