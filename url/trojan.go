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

	// Set remark in config
	config.Remarks = remark

	return config, nil
}

// configToTrojanURL converts a V2Ray configuration back to a Trojan URL
func configToTrojanURL(config *models.V2RayConfig, outbound *models.OutboundConfig) (string, error) {
	// Extract server info from outbound settings
	servers, ok := outbound.Settings["servers"].([]interface{})
	if !ok || len(servers) == 0 {
		return "", fmt.Errorf("invalid Trojan outbound settings")
	}
	
	serverInfo, ok := servers[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid server info in Trojan settings")
	}
	
	address, ok := serverInfo["address"].(string)
	if !ok {
		return "", fmt.Errorf("missing address in Trojan server info")
	}
	
	port, ok := serverInfo["port"].(float64)
	if !ok {
		return "", fmt.Errorf("missing port in Trojan server info")
	}
	
	password, ok := serverInfo["password"].(string)
	if !ok {
		return "", fmt.Errorf("missing password in Trojan server info")
	}
	
	// Create URL components
	u := &url.URL{
		Scheme:   "trojan",
		User:     url.User(password),
		Host:     fmt.Sprintf("%s:%d", address, int(port)),
		Fragment: config.Remarks,
	}
	
	// Query parameters
	query := make(url.Values)
	
	// Add flow if present
	if flow, ok := serverInfo["flow"].(string); ok && flow != "" {
		query.Set("flow", flow)
	}
	
	// Add transport settings if available
	if outbound.StreamSettings != nil {
		streamSettings := outbound.StreamSettings
		
		// Network type
		if streamSettings.Network != "" && streamSettings.Network != "tcp" {
			query.Set("type", streamSettings.Network)
		}
		
		// Security type
		if streamSettings.Security != "" && streamSettings.Security != "tls" {
			query.Set("security", streamSettings.Security)
		}
		
		// Extract network-specific settings
		switch streamSettings.Network {
		case "tcp":
			if tcpSettings, ok := streamSettings.TCPSettings.(map[string]interface{}); ok {
				if header, ok := tcpSettings["header"].(map[string]interface{}); ok {
					if headerType, ok := header["type"].(string); ok && headerType != "none" {
						query.Set("headerType", headerType)
					}
				}
			}
		case "kcp":
			if kcpSettings, ok := streamSettings.KCPSettings.(map[string]interface{}); ok {
				if header, ok := kcpSettings["header"].(map[string]interface{}); ok {
					if headerType, ok := header["type"].(string); ok && headerType != "none" {
						query.Set("headerType", headerType)
					}
				}
				
				if seed, ok := kcpSettings["seed"].(string); ok && seed != "" {
					query.Set("seed", seed)
				}
			}
		case "ws":
			if wsSettings, ok := streamSettings.WSSettings.(map[string]interface{}); ok {
				if path, ok := wsSettings["path"].(string); ok && path != "/" && path != "" {
					query.Set("path", path)
				}
				
				if headers, ok := wsSettings["headers"].(map[string]interface{}); ok {
					if host, ok := headers["Host"].(string); ok && host != "" {
						query.Set("host", host)
					}
				}
			}
		case "h2", "http":
			if httpSettings, ok := streamSettings.HTTPSettings.(map[string]interface{}); ok {
				if path, ok := httpSettings["path"].(string); ok && path != "/" && path != "" {
					query.Set("path", path)
				}
				
				if hosts, ok := httpSettings["host"].([]interface{}); ok && len(hosts) > 0 {
					if host, ok := hosts[0].(string); ok && host != "" {
						query.Set("host", host)
					}
				}
			}
		case "quic":
			if quicSettings, ok := streamSettings.QUICSettings.(map[string]interface{}); ok {
				if header, ok := quicSettings["header"].(map[string]interface{}); ok {
					if headerType, ok := header["type"].(string); ok && headerType != "none" {
						query.Set("headerType", headerType)
					}
				}
				
				if security, ok := quicSettings["security"].(string); ok && security != "" {
					query.Set("quicSecurity", security)
				}
				
				if key, ok := quicSettings["key"].(string); ok && key != "" {
					query.Set("key", key)
				}
			}
		case "grpc":
			if grpcSettings, ok := streamSettings.GRPCSettings.(map[string]interface{}); ok {
				if serviceName, ok := grpcSettings["serviceName"].(string); ok && serviceName != "" {
					query.Set("serviceName", serviceName)
				}
				
				if multiMode, ok := grpcSettings["multiMode"].(bool); ok && multiMode {
					query.Set("mode", "multi")
				}
			}
		}
		
		// TLS settings
		if streamSettings.Security != "" {
			query.Set("security", streamSettings.Security)
			
			if tlsSettings, ok := streamSettings.TLSSettings.(map[string]interface{}); ok {
				if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
					query.Set("sni", serverName)
				}
				
				if alpn, ok := tlsSettings["alpn"].([]interface{}); ok && len(alpn) > 0 {
					if alpnStr, ok := alpn[0].(string); ok && alpnStr != "" {
						query.Set("alpn", alpnStr)
					}
				}
				
				if fingerprint, ok := tlsSettings["fingerprint"].(string); ok && fingerprint != "randomized" && fingerprint != "" {
					query.Set("fp", fingerprint)
				}
			}
		}
	}
	
	u.RawQuery = query.Encode()
	return u.String(), nil
}