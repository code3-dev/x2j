// Package android provides Android bindings for the X2J library
package android

import (
	"encoding/json"
	"pira/x2j/url"
)

// ParseV2RayURL converts a V2Ray URL to Xray JSON configuration
func ParseV2RayURL(v2rayURL string) (string, error) {
	config, err := url.ParseV2RayURL(v2rayURL)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}