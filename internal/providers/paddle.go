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

type PaddleProvider struct{}

func (p *PaddleProvider) Name() string {
	return "paddle"
}

func (p *PaddleProvider) ReSign(body []byte, originalHeaders http.Header, secret string) (http.Header, error) {
	// Paddle header: ts=1614556800;h1=abc123...
	// signed_payload = "{ts}:{body}"
	
	newTimestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(newTimestamp, 10)
	
	signedPayload := timestampStr + ":" + string(body)
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	newSignature := hex.EncodeToString(mac.Sum(nil))
	
	newHeaderValue := fmt.Sprintf("ts=%s;h1=%s", timestampStr, newSignature)
	
	newHeaders := make(http.Header)
	for k, v := range originalHeaders {
		newHeaders[k] = v
	}
	newHeaders.Set("Paddle-Signature", newHeaderValue)
	
	return newHeaders, nil
}
