package providers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStripeReSign(t *testing.T) {
	provider := &StripeProvider{}
	body := []byte(`{"id": "evt_123"}`)
	secret := "whsec_test"
	
	headers, err := provider.ReSign(body, make(http.Header), secret)
	if err != nil {
		t.Fatalf("ReSign failed: %v", err)
	}
	
	sigHeader := headers.Get("Stripe-Signature")
	if sigHeader == "" {
		t.Fatal("Stripe-Signature header missing")
	}
	
	// Verify the signature manually
	// Expected format: t=timestamp,v1=signature
	parts := make(map[string]string)
	for _, part := range strings.Split(sigHeader, ",") {
		kv := strings.Split(part, "=")
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}
	
	ts := parts["t"]
	sig := parts["v1"]
	
	if ts == "" || sig == "" {
		t.Fatalf("Malformed signature header: %s", sigHeader)
	}
	
	// Check if timestamp is recent (within 5 seconds)
	tsInt, _ := strconv.ParseInt(ts, 10, 64)
	if time.Now().Unix()-tsInt > 5 {
		t.Error("Timestamp is too old")
	}
	
	signedPayload := ts + "." + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	
	if sig != expectedSig {
		t.Errorf("Signature mismatch. Got %s, want %s", sig, expectedSig)
	}
}
