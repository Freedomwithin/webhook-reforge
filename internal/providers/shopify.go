package providers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

type ShopifyProvider struct{}

func (p *ShopifyProvider) Name() string {
	return "shopify"
}

func (p *ShopifyProvider) ReSign(body []byte, originalHeaders http.Header, secret string) (http.Header, error) {
	// Shopify signs the raw body with HMAC-SHA256 and base64 encodes the result
	// Header: X-Shopify-Hmac-SHA256
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	newSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	
	newHeaders := make(http.Header)
	for k, v := range originalHeaders {
		newHeaders[k] = v
	}
	newHeaders.Set("X-Shopify-Hmac-SHA256", newSignature)
	
	return newHeaders, nil
}
