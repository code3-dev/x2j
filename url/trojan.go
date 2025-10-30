package url

import (
	"fmt"
	"net/url"

	"pira/x2j/models"
)

// parseTrojanURL parses a Trojan URL and returns a V2Ray configuration
func parseTrojanURL(trojanURL string) (*models.V2RayConfig, error) {
	// Parse the URL
	parsedURL, err := url.Parse(trojanURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Trojan URL: %v", err)
	}
	
	// Create base config
	config := createBaseConfig()
	
	// Extract values
	password := parsedURL.User.Username()
	address := parsedURL.Hostname()
	port := 443
	if portStr := parsedURL.Port(); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
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
	security := getStringQueryParam(query, "security", "tls")
	sni := getStringQueryParam(query, "sni", "")
	fp := getStringQueryParam(query, "fp", "randomized")
	alpn := getStringQueryParam(query, "alpn", "")
	flow := getStringQueryParam(query, "flow", "")
	
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
		Protocol: "trojan",
		Settings: map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address": address,
					"method":  "chacha20-poly1305",
					"password": password,
					"port":    port,
					"level":   8,
					"flow":    flow,
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
	
	// Set remark in config (if needed for future use)
	_ = remark
	
	return config, nil
}