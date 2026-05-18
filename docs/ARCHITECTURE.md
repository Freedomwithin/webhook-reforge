## WebhookForge — Build Outline

---

### What It Actually Does (Technical Core)

Stripe's signature header looks like this:
```
Stripe-Signature: t=1614556800,v1=abc123...
```
The `t` is a Unix timestamp. Stripe rejects anything older than 5 minutes. Paddle's tolerance is ~5 seconds. When you save a real webhook payload and try to replay it later, the timestamp is stale, the signature no longer matches, and your validation code rejects it.

WebhookForge intercepts the payload, swaps `t` to the current time, recomputes `HMAC-SHA256(new_t.payload, secret)`, reconstructs the header, and forwards it to your local app. Your app sees a valid, fresh webhook every time.

---

### Stack

**Go.** Single binary, no runtime, ships as a download, fast HTTP stdlib, crypto/hmac built in. This is the correct choice — don't overthink it.

---

### Phase 1 — MVP (Target: 1 Weekend)

**Scope: Stripe only. CLI only. Works.**

**What to build:**

A Go binary with two modes:

```bash
# Mode 1: Proxy — sits between ngrok/saved payload sender and your app
webhookforge proxy --port 9000 --target http://localhost:3000/webhooks --secret whsec_xxx

# Mode 2: Replay — feed it a saved JSON file, it re-signs and fires it
webhookforge replay --file ./saved_events/payment_intent.json --target http://localhost:3000/webhooks --secret whsec_xxx
```

**Core logic for Stripe re-signing:**
1. Receive POST with body and `Stripe-Signature` header
2. Parse header, extract existing payload
3. Get current Unix timestamp `t_new`
4. Compute `signed_payload = "{t_new}.{raw_body}"`
5. Compute `sig = HMAC-SHA256(signed_payload, secret)`
6. Reconstruct header: `t={t_new},v1={sig}`
7. Forward POST to target with new header

**Config file** (`webhookforge.yaml`):
```yaml
providers:
  stripe:
    secret: whsec_xxx
    target: http://localhost:3000/webhooks/stripe
    listen_port: 9000
```

**Deliverable:** Single Go binary, config file, README with setup in under 5 minutes, MIT licensed, on GitHub.

---

### Phase 2 — Provider Expansion + TUI (2-3 Weeks)

Add Paddle and Shopify. Each has a different signing scheme:

| Provider | Header | Signs | Timestamp |
|----------|--------|-------|-----------|
| Stripe | `Stripe-Signature` | `t.payload` | Yes, 5 min window |
| Paddle | `Paddle-Signature` | `ts:payload` | Yes, ~5 sec window |
| Shopify | `X-Shopify-Hmac-SHA256` | raw body only | No timestamp |
| GitHub | `X-Hub-Signature-256` | raw body only | No timestamp |

Each becomes a **provider adapter** — an interface with `Sign(payload []byte, secret string) string` and `Verify(...)`. Clean, extensible.

Also add a **local web UI on localhost:9001** (just plain HTML served by the binary, no external dependency):
- Payload editor — paste or load a saved event, edit fields, fire it
- Event log — history of replayed events with status codes from your app
- One-click re-fire from history

---

### Phase 3 — Monetization Layer (Desktop App)

Wrap the Go daemon with **Wails** (Go-native desktop framework, produces a real desktop app with a webview). This becomes the paid tier.

Paid features:
- System tray daemon — runs in background, always ready
- Visual payload editor with field highlighting
- Saved event library — organize and name your test scenarios
- Bulk replay — run a sequence of events in order (e.g., simulate a full checkout flow)
- Provider secret manager (encrypted local storage)

**Pricing:** $49 flat license. No subscription. One-time. This is the right call for a developer utility — subscriptions create friction, flat licenses get expensed immediately.

Free tier (CLI + basic proxy) stays open source forever. That's your distribution engine.

---

### File Structure

```
webhookforge/
├── cmd/
│   └── webhookforge/
│       └── main.go          # CLI entry point
├── internal/
│   ├── proxy/
│   │   └── proxy.go         # HTTP proxy server
│   ├── providers/
│   │   ├── provider.go      # Interface definition
│   │   ├── stripe.go        # Stripe adapter
│   │   ├── paddle.go        # Paddle adapter
│   │   └── shopify.go       # Shopify adapter
│   ├── replay/
│   │   └── replay.go        # File-based replay engine
│   └── config/
│       └── config.go        # YAML config loader
├── web/                     # Phase 2 local UI (embedded)
├── webhookforge.yaml.example
├── README.md
└── go.mod
```

---

### Distribution Plan

- GitHub release with pre-built binaries for Linux, Mac, Windows
- Homebrew tap: `brew install webhookforge`
- Post to: r/golang, r/webdev, r/SaaS, Hacker News Show HN, relevant Discord servers (Stripe devs, Indie Hackers)
- The Show HN post is important — "I built a local webhook re-signing proxy because Paddle's 5-second tolerance window is brutal" is a concrete, relatable hook

---

### Execution Order

1. Get the Stripe re-signing logic working in isolation first — just a function that takes a payload and secret and returns a valid header. Test it against your own Stripe test keys.
2. Wrap it in the proxy server.
3. Add the replay-from-file mode.
4. Write the README.
5. Ship v0.1 to GitHub.
6. Then expand providers and build the UI.
