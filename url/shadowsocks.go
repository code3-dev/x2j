package url

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"pira/x2j/models"
)

// parseShadowSocksURL parses a ShadowSocks URL and returns a V2Ray configuration
func parseShadowSocksURL(ssURL string) (*models.V2RayConfig, error) {
	// Parse the URL
	parsedURL, err := url.Parse(ssURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ShadowSocks URL: %v", err)
	}
	
	// Create base config
	config := createBaseConfig()
	
	// Extract values
	method := "none"
	password := ""
	address := parsedURL.Hostname()
	port := 443
	if portStr := parsedURL.Port(); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}
	
	// Parse user info (method:password encoded in base64)
	if parsedURL.User != nil {
		userInfo := parsedURL.User.String()
		// Add padding if needed
		if len(userInfo)%4 > 0 {
			userInfo += strings.Repeat("=", 4-len(userInfo)%4)
		}
		
		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(userInfo)
		if err == nil {
			// Split method and password
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) == 2 {
				method = parts[0]
				password = parts[1]
			}
		}
	}
	
	// Query parameters
	query := parsedURL.Query()
	remark := parsedURL.Fragment
	
	// Transport settings
	transport := getStringQueryParam(query, "type", "tcp")
	headerType := getStringQueryParam(query, "headerType", "")
	host := getStringQueryParam(query, "host", "")
	path := getStringQueryParam(query, "path", "")
	seed := getStringQueryParam(query, "seed", "")
	quicSecurity := getStringQueryParam(query, "quicSecurity", "")
	key := getStringQueryParam(query, "key", "")
	mode := getStringQueryParam(query, "mode", "")
	serviceName := getStringQueryParam(query, "serviceName", "")
	
	// TLS settings
	security := getStringQueryParam(query, "security", "")
	sni := getStringQueryParam(query, "sni", "")
	fp := getStringQueryParam(query, "fp", "")
	alpn := getStringQueryParam(query, "alpn", "")
	
	// Create stream settings
	streamSetting := &models.StreamSettings{}
	transportSni := populateTransportSettings(streamSetting, transport, headerType, host, path, seed, quicSecurity, key, mode, serviceName)
	
	// Use SNI from query parameters if provided, otherwise use transport SNI
	if sni == "" {
		sni = transportSni
	}
	
	// TLS settings
	populateTLSSettings(streamSetting, security, true, sni, fp, alpn, "", "", "")
	
	// Create outbound configuration
	outbound := models.OutboundConfig{
		Tag:      "proxy",
		Protocol: "shadowsocks",
		Settings: map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address": address,
					"method":  method,
					"password": password,
					"port":    port,
					"level":   8,
				},
			},
		},
		StreamSettings: streamSetting,
		Mux: &models.MuxConfig{
			Enabled:     false,
			Concurrency: 8,
		},
	}
	
	// Set the proxy outbound as the first outbound
	config.Outbounds = append([]models.OutboundConfig{outbound}, config.Outbounds...)

	// Set remark in config
	config.Remarks = remark

	return config, nil
}