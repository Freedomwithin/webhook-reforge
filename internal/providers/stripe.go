package providers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type StripeProvider struct{}

func (p *StripeProvider) Name() string {
	return "stripe"
}

func (p *StripeProvider) ReSign(body []byte, originalHeaders http.Header, secret string) (http.Header, error) {
	// Stripe header: t=1614556800,v1=abc123...
	// We only need the body to re-sign with a NEW timestamp.
	
	newTimestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(newTimestamp, 10)
	
	// signed_payload = "{timestamp}.{body}"
	signedPayload := timestampStr + "." + string(body)
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	newSignature := hex.EncodeToString(mac.Sum(nil))
	
	newHeaderValue := fmt.Sprintf("t=%s,v1=%s", timestampStr, newSignature)
	
	newHeaders := make(http.Header)
	for k, v := range originalHeaders {
		newHeaders[k] = v
	}
	newHeaders.Set("Stripe-Signature", newHeaderValue)
	
	return newHeaders, nil
}
