package url

import (
	"encoding/json"
	"fmt"
	"net/url"

	"pira/x2j/models"
)

// parseVLessURL parses a VLess URL and returns a V2Ray configuration
func parseVLessURL(vlessURL string) (*models.V2RayConfig, error) {
	// Parse the URL
	parsedURL, err := url.Parse(vlessURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VLess URL: %v", err)
	}
	
	// Create base config
	config := createBaseConfig()
	
	// Extract values
	uuid := parsedURL.User.Username()
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
	security := getStringQueryParam(query, "security", "")
	sni := getStringQueryParam(query, "sni", "")
	fp := getStringQueryParam(query, "fp", "")
	alpn := getStringQueryParam(query, "alpn", "")
	pbk := getStringQueryParam(query, "pbk", "")
	sid := getStringQueryParam(query, "sid", "")
	spx := getStringQueryParam(query, "spx", "")
	encryption := getStringQueryParam(query, "encryption", "none")
	flow := getStringQueryParam(query, "flow", "")
	
	// Create stream settings
	streamSetting := &models.StreamSettings{}
	transportSni := populateTransportSettings(streamSetting, transport, headerType, host, path, seed, quicSecurity, key, mode, serviceName)
	
	// Use SNI from query parameters if provided, otherwise use transport SNI
	if sni == "" {
		sni = transportSni
	}
	
	// TLS settings
	populateTLSSettings(streamSetting, security, true, sni, fp, alpn, pbk, sid, spx)
	
	// Handle xhttp specific settings
	if transport == "xhttp" {
		handleXhttpSettings(streamSetting, query)
	}
	
	// Create outbound configuration
	outbound := models.OutboundConfig{
		Tag:      "proxy",
		Protocol: "vless",
		Settings: map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": address,
					"port":    port,
					"users": []interface{}{
						map[string]interface{}{
							"id":         uuid,
							"encryption": encryption,
							"flow":       flow,
							"level":      8,
						},
					},
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

// handleXhttpSettings handles xhttp specific settings
func handleXhttpSettings(streamSetting *models.StreamSettings, query map[string][]string) {
	extra := getStringQueryParam(query, "extra", "")
	
	xhttpSettings := map[string]interface{}{
		"host": getStringQueryParam(query, "host", ""),
		"path": getStringQueryParam(query, "path", "/"),
		"mode": getStringQueryParam(query, "mode", "auto"),
	}
	
	// Add extra settings if available
	if extra != "" {
		var extraData map[string]interface{}
		if err := json.Unmarshal([]byte(extra), &extraData); err == nil {
			xhttpSettings["extra"] = extraData
		}
	}
	
	streamSetting.XHTTPSettings = xhttpSettings
}

// getStringQueryParam gets a string query parameter with a default value
func getStringQueryParam(query map[string][]string, key, defaultValue string) string {
	if values, ok := query[key]; ok && len(values) > 0 {
		return values[0]
	}
	return defaultValue
}