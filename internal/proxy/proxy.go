package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Freedomwithin/webhook-reforge/internal/providers"
)

type ProxyServer struct {
	Port     int
	Target   string
	Secret   string
	Provider providers.Provider
}

func (s *ProxyServer) Start() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		// Re-sign the payload
		newHeaders, err := s.Provider.ReSign(body, r.Header, s.Secret)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to re-sign: %v", err), http.StatusInternalServerError)
			return
		}

		// Forward to target
		targetURL, err := url.Parse(s.Target)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxyReq, err := http.NewRequest(http.MethodPost, targetURL.String(), bytes.NewReader(body))
		if err != nil {
			http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for k, v := range newHeaders {
			if k == "Content-Length" {
				continue
			}
			proxyReq.Header[k] = v
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Return target's response
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	fmt.Printf("🚀 WebhookForge Proxy listening on :%d -> %s\n", s.Port, s.Target)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), handler)
}
