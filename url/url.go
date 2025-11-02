package url

import (
	"fmt"
	"strings"

	"pira/x2j/models"
)

// ParseV2RayURL parses a V2Ray URL and returns a V2Ray configuration
func ParseV2RayURL(v2rayURL string) (*models.V2RayConfig, error) {
	switch {
	case strings.HasPrefix(v2rayURL, "vmess://"):
		return parseVMessURL(v2rayURL)
	case strings.HasPrefix(v2rayURL, "vless://"):
		return parseVLessURL(v2rayURL)
	case strings.HasPrefix(v2rayURL, "ss://"):
		return parseShadowSocksURL(v2rayURL)
	case strings.HasPrefix(v2rayURL, "trojan://"):
		return parseTrojanURL(v2rayURL)
	default:
		return nil, fmt.Errorf("unsupported protocol in URL: %s", v2rayURL)
	}
}

// ConfigToV2RayURL converts a V2Ray configuration back to a share URL
func ConfigToV2RayURL(config *models.V2RayConfig) (string, error) {
	// Find the proxy outbound (first outbound is typically the proxy)
	if len(config.Outbounds) == 0 {
		return "", fmt.Errorf("no outbounds found in configuration")
	}

	proxyOutbound := config.Outbounds[0]
	
	switch proxyOutbound.Protocol {
	case "vmess":
		return configToVMessURL(config, &proxyOutbound)
	case "vless":
		return configToVLessURL(config, &proxyOutbound)
	case "shadowsocks":
		return configToShadowSocksURL(config, &proxyOutbound)
	case "trojan":
		return configToTrojanURL(config, &proxyOutbound)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", proxyOutbound.Protocol)
	}
}

// getStringQueryParam gets a string query parameter with a default value
func getStringQueryParam(query map[string][]string, key, defaultValue string) string {
	if values, ok := query[key]; ok && len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

// createBaseConfig creates the base V2Ray configuration structure
func createBaseConfig() *models.V2RayConfig {
	return &models.V2RayConfig{
		Log: models.LogConfig{
			Access:   "",
			Error:    "",
			LogLevel: "error",
			DNSLog:   false,
		},
		Inbounds: []models.InboundConfig{
			{
				Tag:      "in_proxy",
				Port:     1080,
				Protocol: "socks",
				Listen:   "127.0.0.1",
				Settings: map[string]interface{}{
					"auth":      "noauth",
					"udp":       true,
					"userLevel": 8,
				},
			},
		},
		Outbounds: []models.OutboundConfig{
			{
				Tag:      "direct",
				Protocol: "freedom",
				Settings: map[string]interface{}{},
			},
			{
				Tag:      "blackhole",
				Protocol: "blackhole",
				Settings: map[string]interface{}{},
			},
		},
		DNS: models.DNSConfig{
			Servers: []string{"1.1.1.1", "1.0.0.1", "8.8.8.8", "8.8.4.4"},
		},
		Routing: models.RoutingConfig{
			DomainStrategy: "UseIp",
			Rules:          []interface{}{},
			Balancers:      []interface{}{},
		},
	}
}

// populateTransportSettings populates transport settings based on parameters
func populateTransportSettings(streamSetting *models.StreamSettings, transport, headerType, host, path, seed, quicSecurity, key, mode, serviceName string) string {
	var sni string
	streamSetting.Network = transport
	
	switch transport {
	case "tcp":
		tcpSettings := map[string]interface{}{
			"header": map[string]interface{}{
				"type": "none",
			},
		}
		
		if headerType == "http" {
			header := tcpSettings["header"].(map[string]interface{})
			header["type"] = "http"
			
			request := map[string]interface{}{
				"path": []string{"/"},
				"headers": map[string]interface{}{
					"Host": []string{""},
					"User-Agent": []string{
						"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
						"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46",
					},
					"Accept-Encoding": []string{
						"gzip, deflate",
					},
					"Connection": []string{
						"keep-alive",
					},
					"Pragma": "no-cache",
				},
				"version": "1.1",
				"method":  "GET",
			}
			
			if path != "" {
				request["path"] = []string{path}
			}
			
			if host != "" {
				hostList := []string{host}
				request["headers"].(map[string]interface{})["Host"] = hostList
				if len(hostList) > 0 {
					sni = hostList[0]
				}
			}
			
			header["request"] = request
		}
		
		streamSetting.TCPSettings = tcpSettings
		
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":               1350,
			"tti":               50,
			"uplinkCapacity":    12,
			"downlinkCapacity":  100,
			"congestion":        false,
			"readBufferSize":    1,
			"writeBufferSize":   1,
			"header": map[string]interface{}{
				"type": headerType,
			},
		}
		
		if seed != "" {
			kcpSettings["seed"] = seed
		}
		
		streamSetting.KCPSettings = kcpSettings
		
	case "ws":
		wsSettings := map[string]interface{}{
			"path": "/",
			"headers": map[string]interface{}{
				"Host": "",
			},
		}
		
		if path != "" {
			wsSettings["path"] = path
		}
		
		if host != "" {
			wsSettings["headers"].(map[string]interface{})["Host"] = host
			sni = host
		}
		
		streamSetting.WSSettings = wsSettings
		
	case "h2", "http":
		streamSetting.Network = "h2"
		h2Settings := map[string]interface{}{
			"path": "/",
		}
		
		if host != "" {
			h2Settings["host"] = []string{host}
			hostList := []string{host}
			if len(hostList) > 0 {
				sni = hostList[0]
			}
		}
		
		if path != "" {
			h2Settings["path"] = path
		}
		
		streamSetting.HTTPSettings = h2Settings
		
	case "quic":
		quicSettings := map[string]interface{}{
			"security": quicSecurity,
			"key":      key,
			"header": map[string]interface{}{
				"type": headerType,
			},
		}
		
		streamSetting.QUICSettings = quicSettings
		
	case "grpc":
		grpcSettings := map[string]interface{}{
			"serviceName": serviceName,
			"multiMode":   mode == "multi",
		}
		
		streamSetting.GRPCSettings = grpcSettings
		if host != "" {
			sni = host
		}
		
	case "xhttp":
		xhttpSettings := map[string]interface{}{
			"host": host,
			"path": "/",
			"mode": "auto",
		}
		
		if path != "" {
			xhttpSettings["path"] = path
		}
		
		if mode != "" {
			xhttpSettings["mode"] = mode
		}
		
		streamSetting.XHTTPSettings = xhttpSettings
		if host != "" {
			sni = host
		}
		
	case "httpupgrade":
		httpupgradeSettings := map[string]interface{}{
			"host": "",
			"path": "",
		}
		
		if host != "" {
			httpupgradeSettings["host"] = host
			sni = host
		}
		
		if path != "" {
			httpupgradeSettings["path"] = path
		}
		
		streamSetting.HTTPUpgradeSettings = httpupgradeSettings
	}
	
	return sni
}

// populateTLSSettings populates TLS/reality settings
func populateTLSSettings(streamSetting *models.StreamSettings, streamSecurity string, allowInsecure bool, sni, fingerprint, alpns, publicKey, shortId, spiderX string) {
	streamSetting.Security = streamSecurity
	
	tlsSetting := map[string]interface{}{
		"allowInsecure": allowInsecure,
		"serverName":    sni,
	}
	
	if alpns != "" {
		tlsSetting["alpn"] = []string{alpns}
	}
	
	if fingerprint != "" {
		tlsSetting["fingerprint"] = fingerprint
	}
	
	if streamSecurity == "tls" {
		streamSetting.RealitySettings = nil
		streamSetting.TLSSettings = tlsSetting
	} else if streamSecurity == "reality" {
		tlsSetting["publicKey"] = publicKey
		tlsSetting["shortId"] = shortId
		tlsSetting["spiderX"] = spiderX
		
		streamSetting.TLSSettings = nil
		streamSetting.RealitySettings = tlsSetting
	}
}