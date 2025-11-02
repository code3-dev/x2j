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
	remarks := getStringValue(rawConfig, "ps", "") // "ps" is commonly used for remarks in VMess

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

	// Set remarks in config
	config.Remarks = remarks

	return config, nil
}

// configToVMessURL converts a V2Ray configuration back to a VMess URL
func configToVMessURL(config *models.V2RayConfig, outbound *models.OutboundConfig) (string, error) {
	// Extract server info from outbound settings
	settings, ok := outbound.Settings["vnext"].([]interface{})
	if !ok || len(settings) == 0 {
		return "", fmt.Errorf("invalid VMess outbound settings")
	}
	
	serverInfo, ok := settings[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid server info in VMess settings")
	}
	
	address, ok := serverInfo["address"].(string)
	if !ok {
		return "", fmt.Errorf("missing address in VMess server info")
	}
	
	port, ok := serverInfo["port"].(float64)
	if !ok {
		return "", fmt.Errorf("missing port in VMess server info")
	}
	
	users, ok := serverInfo["users"].([]interface{})
	if !ok || len(users) == 0 {
		return "", fmt.Errorf("missing users in VMess server info")
	}
	
	userInfo, ok := users[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid user info in VMess settings")
	}
	
	id, ok := userInfo["id"].(string)
	if !ok {
		return "", fmt.Errorf("missing id in VMess user info")
	}
	
	alterId := 0
	if aid, ok := userInfo["alterId"]; ok {
		if aidFloat, ok := aid.(float64); ok {
			alterId = int(aidFloat)
		}
	}
	
	security := "auto"
	if sec, ok := userInfo["security"].(string); ok {
		security = sec
	}
	
	// Create VMess JSON object
	vmessConfig := map[string]interface{}{
		"v":    "2",
		"ps":   config.Remarks,
		"add":  address,
		"port": int(port),
		"id":   id,
		"aid":  alterId,
		"scy":  security,
		"net":  "tcp",
		"type": "none",
		"host": "",
		"path": "",
		"tls":  "",
	}
	
	// Extract transport settings if available
	if outbound.StreamSettings != nil {
		streamSettings := outbound.StreamSettings
		
		// Network type
		if streamSettings.Network != "" {
			vmessConfig["net"] = streamSettings.Network
		}
		
		// TLS settings
		if streamSettings.Security != "" {
			vmessConfig["tls"] = streamSettings.Security
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
					vmessConfig["type"] = headerType
					
					if headerType == "http" {
						if request, ok := header["request"].(map[string]interface{}); ok {
							// Extract path
							if pathArr, ok := request["path"].([]interface{}); ok && len(pathArr) > 0 {
								if pathStr, ok := pathArr[0].(string); ok {
									vmessConfig["path"] = pathStr
								}
							}
							
							// Extract host
							if headers, ok := request["headers"].(map[string]interface{}); ok {
								if hostArr, ok := headers["Host"].([]interface{}); ok && len(hostArr) > 0 {
									if hostStr, ok := hostArr[0].(string); ok {
										vmessConfig["host"] = hostStr
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
					vmessConfig["type"] = headerType
				}
				
				if seed, ok := kcpSettings["seed"].(string); ok {
					vmessConfig["path"] = seed
				}
			}
		case "ws":
			if wsSettings, ok := streamSettings.WSSettings.(map[string]interface{}); ok {
				if path, ok := wsSettings["path"].(string); ok {
					vmessConfig["path"] = path
				}
				
				if headers, ok := wsSettings["headers"].(map[string]interface{}); ok {
					if host, ok := headers["Host"].(string); ok {
						vmessConfig["host"] = host
					}
				}
			}
		case "h2", "http":
			if httpSettings, ok := streamSettings.HTTPSettings.(map[string]interface{}); ok {
				if path, ok := httpSettings["path"].(string); ok {
					vmessConfig["path"] = path
				}
				
				if hosts, ok := httpSettings["host"].([]interface{}); ok && len(hosts) > 0 {
					if host, ok := hosts[0].(string); ok {
						vmessConfig["host"] = host
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
					vmessConfig["type"] = headerType
				}
				
				if security, ok := quicSettings["security"].(string); ok {
					vmessConfig["host"] = security
				}
				
				if key, ok := quicSettings["key"].(string); ok {
					vmessConfig["path"] = key
				}
			}
		case "grpc":
			if grpcSettings, ok := streamSettings.GRPCSettings.(map[string]interface{}); ok {
				if serviceName, ok := grpcSettings["serviceName"].(string); ok {
					vmessConfig["path"] = serviceName
				}
				
				if multiMode, ok := grpcSettings["multiMode"].(bool); ok && multiMode {
					vmessConfig["type"] = "multi"
				}
			}
		}
		
		// TLS settings
		if tlsSettings, ok := streamSettings.TLSSettings.(map[string]interface{}); ok {
			if serverName, ok := tlsSettings["serverName"].(string); ok && serverName != "" {
				vmessConfig["host"] = serverName
			}
			
			if alpn, ok := tlsSettings["alpn"].([]interface{}); ok && len(alpn) > 0 {
				if alpnStr, ok := alpn[0].(string); ok {
					vmessConfig["alpn"] = alpnStr
				}
			}
			
			if fingerprint, ok := tlsSettings["fingerprint"].(string); ok {
				vmessConfig["fp"] = fingerprint
			}
		}
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(vmessConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal VMess config: %v", err)
	}
	
	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	
	return "vmess://" + encoded, nil
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