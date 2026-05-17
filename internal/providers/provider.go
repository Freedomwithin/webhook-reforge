package providers

import "net/http"

// Provider defines the interface for webhook providers
type Provider interface {
	Name() string
	// ReSign takes an existing request, re-signs it with the current timestamp and secret,
	// and returns the updated headers.
	ReSign(body []byte, originalHeaders http.Header, secret string) (http.Header, error)
}
