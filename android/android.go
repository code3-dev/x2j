// Package android provides Android bindings for the X2J library
package android

import (
	_ "golang.org/x/mobile/bind"
	"encoding/json"
	"strings"
	"pira/x2j/url"
)

// ParseV2RayURL converts a V2Ray URL to Xray JSON configuration
func ParseV2RayURL(v2rayURL string) (string, error) {
	return ParseV2RayURLWithSettings(v2rayURL, 1080, "")
}

// ParseV2RayURLWithSettings converts a V2Ray URL to Xray JSON configuration with custom port and DNS
func ParseV2RayURLWithSettings(v2rayURL string, port int, dns string) (string, error) {
	config, err := url.ParseV2RayURL(v2rayURL)
	if err != nil {
		return "", err
	}

	// Update the inbound port if a custom port is specified
	if port != 1080 {
		for i := range config.Inbounds {
			if config.Inbounds[i].Tag == "in_proxy" {
				config.Inbounds[i].Port = port
				break
			}
		}
	}

	// Update DNS servers if custom DNS is specified
	if dns != "" {
		if dns == "\"\"" || dns == "''" {
			// Empty DNS - clear the servers list
			config.DNS.Servers = []string{}
		} else {
			// Parse comma-separated DNS servers
			dnsServers := strings.Split(dns, ",")
			// Trim whitespace from each server
			for i, server := range dnsServers {
				dnsServers[i] = strings.TrimSpace(server)
			}
			config.DNS.Servers = dnsServers
		}
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}