package url

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"pira/x2j/models"
)

// parseVMessURL parses a VMess URL and returns a V2Ray configuration
func parseVMessURL(vmessURL string) (*models.V2RayConfig, error) {
	// Remove the scheme prefix
	raw := strings.TrimPrefix(vmessURL, "vmess://")
	
	// Add padding if needed
	if len(raw)%4 > 0 {
		raw += strings.Repeat("=", 4-len(raw)%4)
	}
	
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode VMess URL: %v", err)
	}
	
	// Parse JSON
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(decoded, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse VMess JSON: %v", err)
	}
	
	// Create base config
	config := createBaseConfig()
	
	// Extract values with defaults
	address := getStringValue(rawConfig, "add", "")
	port := getIntValue(rawConfig, "port", 443)
	id := getStringValue(rawConfig, "id", "")
	alterId := getIntValue(rawConfig, "aid", 0)
	security := getStringValue(rawConfig, "scy", "auto")
	level := 8
	
	// Transport settings
	net := getStringValue(rawConfig, "net", "tcp")
	headerType := getStringValue(rawConfig, "type", "")
	host := getStringValue(rawConfig, "host", "")
	path := getStringValue(rawConfig, "path", "")
	tls := getStringValue(rawConfig, "tls", "")
	alpn := getStringValue(rawConfig, "alpn", "")
	fp := getStringValue(rawConfig, "fp", "")
	
	// Create stream settings
	streamSetting := &models.StreamSettings{}
	sni := populateTransportSettings(streamSetting, net, headerType, host, path, path, host, path, headerType, path)
	
	// TLS settings
	populateTLSSettings(streamSetting, tls, true, sni, fp, alpn, "", "", "")
	
	// Create outbound configuration
	outbound := models.OutboundConfig{
		Tag:      "proxy",
		Protocol: "vmess",
		Settings: map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": address,
					"port":    port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       id,
							"alterId":  alterId,
							"security": security,
							"level":    level,
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
	
	return config, nil
}

// Helper functions for extracting values
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if str, ok := val.(string); ok {
			var result int
			if _, err := fmt.Sscanf(str, "%d", &result); err == nil {
				return result
			}
		}
	}
	return defaultValue
}