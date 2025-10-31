// Package android provides Android bindings for the X2J library
//
// This package is designed to be used with gomobile to create an Android library (.aar)
// that can convert V2Ray URLs to Xray JSON configurations. It provides both simple and
// advanced conversion functions with customizable settings.
package android

import (
	_ "golang.org/x/mobile/bind"
	"encoding/json"
	"fmt"
	"strings"
	"pira/x2j/url"
)

// ParseV2RayURL converts a V2Ray URL to Xray JSON configuration using default settings.
// It uses port 1080 and the system's default DNS servers.
//
// Parameters:
//   - v2rayURL: The V2Ray URL to convert (e.g., vmess://, vless://, ss://, trojan://)
//
// Returns:
//   - string: The Xray JSON configuration as a formatted string
//   - error: Any error that occurred during conversion
func ParseV2RayURL(v2rayURL string) (string, error) {
	if v2rayURL == "" {
		return "", fmt.Errorf("v2ray URL cannot be empty")
	}
	return ParseV2RayURLWithSettings(v2rayURL, 1080, "", "")
}

// ParseV2RayURLWithSettings converts a V2Ray URL to Xray JSON configuration with custom settings.
// It allows specifying a custom inbound port, DNS servers, and remarks.
//
// Parameters:
//   - v2rayURL: The V2Ray URL to convert (e.g., vmess://, vless://, ss://, trojan://)
//   - port: Custom port for the inbound proxy (default: 1080)
//   - dns: Comma-separated list of DNS servers (e.g., "1.1.1.1,8.8.8.8")
//     Use empty string ("") for default DNS, or '""' to clear DNS servers
//   - remarks: Remarks/Comments for the configuration (optional)
//
// Returns:
//   - string: The Xray JSON configuration as a formatted string
//   - error: Any error that occurred during conversion
func ParseV2RayURLWithSettings(v2rayURL string, port int, dns string, remarks string) (string, error) {
	if v2rayURL == "" {
		return "", fmt.Errorf("v2ray URL cannot be empty")
	}

	if port < 1 || port > 65535 {
		return "", fmt.Errorf("invalid port number: %d (must be between 1 and 65535)", port)
	}

	config, err := url.ParseV2RayURL(v2rayURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse V2Ray URL: %w", err)
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
			// Trim whitespace from each server and validate
			for i, server := range dnsServers {
				server = strings.TrimSpace(server)
				if server == "" {
					return "", fmt.Errorf("invalid DNS server: empty address at position %d", i+1)
				}
				dnsServers[i] = server
			}
			config.DNS.Servers = dnsServers
		}
	}

	// Set remarks if provided
	if remarks != "" {
		config.Remarks = remarks
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}

	return string(jsonData), nil
}