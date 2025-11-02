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

	// Set remark in config
	config.Remarks = remark

	return config, nil
}

// configToVLessURL converts a V2Ray configuration back to a VLess URL
func configToVLessURL(config *models.V2RayConfig, outbound *models.OutboundConfig) (string, error) {
	// Extract server info from outbound settings
	settings, ok := outbound.Settings["vnext"].([]interface{})
	if !ok || len(settings) == 0 {
		return "", fmt.Errorf("invalid VLess outbound settings")
	}
	
	serverInfo, ok := settings[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid server info in VLess settings")
	}
	
	address, ok := serverInfo["address"].(string)
	if !ok {
		return "", fmt.Errorf("missing address in VLess server info")
	}
	
	port, ok := serverInfo["port"].(float64)
	if !ok {
		return "", fmt.Errorf("missing port in VLess server info")
	}
	
	users, ok := serverInfo["users"].([]interface{})
	if !ok || len(users) == 0 {
		return "", fmt.Errorf("missing users in VLess server info")
	}
	
	userInfo, ok := users[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid user info in VLess settings")
	}
	
	uuid, ok := userInfo["id"].(string)
	if !ok {
		return "", fmt.Errorf("missing id in VLess user info")
	}
	
	// Create URL components
	u := &url.URL{
		Scheme:   "vless",
		User:     url.User(uuid),
		Host:     fmt.Sprintf("%s:%d", address, int(port)),
		Fragment: config.Remarks,
	}
	
	// Query parameters
	query := make(url.Values)
	
	// Add encryption if present
	if encryption, ok := userInfo["encryption"].(string); ok && encryption != "" && encryption != "none" {
		query.Set("encryption", encryption)
	}
	
	// Add flow if present
	if flow, ok := userInfo["flow"].(string); ok && flow != "" {
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
		if streamSettings.Security != "" {
			query.Set("security", streamSettings.Security)
		}
		
		// Extract network-specific settings
		switch streamSettings.Network {
		case "tcp":
			if tcpSettings, ok := streamSettings.TCPSettings.(map[string]interface{}); ok {
				if header, ok := tcpSettings["header"].(map[string]interface{}); ok {
					headerType := "" // Define the variable
					if headerTypeVal, ok := header["type"].(string); ok {
						headerType = headerTypeVal
					}
					if headerType != "none" {
						query.Set("headerType", headerType)
					}
					
					if headerType == "http" {
						if request, ok := header["request"].(map[string]interface{}); ok {
							// Extract path
							if pathArr, ok := request["path"].([]interface{}); ok && len(pathArr) > 0 {
								if pathStr, ok := pathArr[0].(string); ok && pathStr != "/" {
									query.Set("path", pathStr)
								}
							}
							
							// Extract host
							if headers, ok := request["headers"].(map[string]interface{}); ok {
								if hostArr, ok := headers["Host"].([]interface{}); ok && len(hostArr) > 0 {
									if hostStr, ok := hostArr[0].(string); ok && hostStr != "" {
										query.Set("host", hostStr)
									}
								}
							}
						}
					}
				}
			}
		case "kcp":
			if kcpSettings, ok := streamSettings.KCPSettings.(map[string]interface{}); ok {
				if header, ok := kcpSettings["header"].(map[string]interface{}); ok {
					headerType := "" // Define the variable
					if headerTypeVal, ok := header["type"].(string); ok {
						headerType = headerTypeVal
					}
					if headerType != "none" {
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
					headerType := "" // Define the variable
					if headerTypeVal, ok := header["type"].(string); ok {
						headerType = headerTypeVal
					}
					if headerType != "none" {
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
		case "xhttp":
			if xhttpSettings, ok := streamSettings.XHTTPSettings.(map[string]interface{}); ok {
				if host, ok := xhttpSettings["host"].(string); ok && host != "" {
					query.Set("host", host)
				}
				
				if path, ok := xhttpSettings["path"].(string); ok && path != "/" && path != "" {
					query.Set("path", path)
				}
				
				if mode, ok := xhttpSettings["mode"].(string); ok && mode != "auto" && mode != "" {
					query.Set("mode", mode)
				}
			}
		case "httpupgrade":
			if httpupgradeSettings, ok := streamSettings.HTTPUpgradeSettings.(map[string]interface{}); ok {
				if host, ok := httpupgradeSettings["host"].(string); ok && host != "" {
					query.Set("host", host)
				}
				
				if path, ok := httpupgradeSettings["path"].(string); ok && path != "/" && path != "" {
					query.Set("path", path)
				}
			}
		}
		
		// TLS/Reality settings
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
				
				if fingerprint, ok := tlsSettings["fingerprint"].(string); ok && fingerprint != "" {
					query.Set("fp", fingerprint)
				}
				
				if publicKey, ok := tlsSettings["publicKey"].(string); ok && publicKey != "" {
					query.Set("pbk", publicKey)
				}
				
				if shortId, ok := tlsSettings["shortId"].(string); ok && shortId != "" {
					query.Set("sid", shortId)
				}
				
				if spiderX, ok := tlsSettings["spiderX"].(string); ok && spiderX != "" {
					query.Set("spx", spiderX)
				}
			}
			
			if realitySettings, ok := streamSettings.RealitySettings.(map[string]interface{}); ok {
				if serverName, ok := realitySettings["serverName"].(string); ok && serverName != "" {
					query.Set("sni", serverName)
				}
				
				if publicKey, ok := realitySettings["publicKey"].(string); ok && publicKey != "" {
					query.Set("pbk", publicKey)
				}
				
				if shortId, ok := realitySettings["shortId"].(string); ok && shortId != "" {
					query.Set("sid", shortId)
				}
				
				if spiderX, ok := realitySettings["spiderX"].(string); ok && spiderX != "" {
					query.Set("spx", spiderX)
				}
				
				if fingerprint, ok := realitySettings["fingerprint"].(string); ok && fingerprint != "" {
					query.Set("fp", fingerprint)
				}
			}
		}
	}
	
	u.RawQuery = query.Encode()
	return u.String(), nil
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