package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Freedomwithin/webhook-reforge/internal/providers"
	"github.com/Freedomwithin/webhook-reforge/internal/proxy"
	"github.com/Freedomwithin/webhook-reforge/internal/replay"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func getProvider(name string) (providers.Provider, error) {
	switch name {
	case "stripe":
		return &providers.StripeProvider{}, nil
	case "paddle":
		return &providers.PaddleProvider{}, nil
	case "shopify":
		return &providers.ShopifyProvider{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// Replay handles re-signing and firing a saved JSON payload
func (a *App) Replay(file, target, secret, providerName string) string {
	provider, err := getProvider(providerName)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	r := &replay.ReplayEngine{
		File:     file,
		Target:   target,
		Secret:   secret,
		Provider: provider,
	}
	result, err := r.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return fmt.Sprintf("✅ Replay Successful | Status: %d | Response: %s", result.StatusCode, result.Body)
}

// ReplayPayload handles re-signing and firing a raw JSON string from the UI
func (a *App) ReplayPayload(jsonPayload, target, secret, providerName string) string {
	provider, err := getProvider(providerName)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "whr-*.json")
	if err != nil {
		return fmt.Sprintf("Error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(jsonPayload)
	tmpFile.Close()

	r := &replay.ReplayEngine{
		File:     tmpFile.Name(),
		Target:   target,
		Secret:   secret,
		Provider: provider,
	}
	result, err := r.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return fmt.Sprintf("✅ Replay Successful | Status: %d | Response: %s", result.StatusCode, result.Body)
}

// StartProxy handles starting the re-signing proxy
func (a *App) StartProxy(port int, target, secret, providerName string) string {
	provider, err := getProvider(providerName)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	p := &proxy.ProxyServer{
		Port:     port,
		Target:   target,
		Secret:   secret,
		Provider: provider,
	}

	go func() {
		if err := p.Start(); err != nil {
			fmt.Printf("Proxy error: %v\n", err)
		}
	}()

	return fmt.Sprintf("🚀 Proxy started on :%d", port)
}
