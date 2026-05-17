package replay

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Freedomwithin/webhookforge/internal/providers"
)

type ReplayEngine struct {
	File     string
	Target   string
	Secret   string
	Provider providers.Provider
}

func (e *ReplayEngine) Run() error {
	body, err := os.ReadFile(e.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Re-sign the payload
	// For replay, we don't have original headers, so we pass empty ones
	newHeaders, err := e.Provider.ReSign(body, make(http.Header), e.Secret)
	if err != nil {
		return fmt.Errorf("failed to re-sign: %w", err)
	}

	// Fire it
	req, err := http.NewRequest(http.MethodPost, e.Target, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range newHeaders {
		req.Header[k] = v
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fire request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Replayed %s to %s\n", e.File, e.Target)
	fmt.Printf("Response: %d %s\n", resp.StatusCode, string(respBody))

	return nil
}
