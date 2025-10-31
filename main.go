package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"pira/x2j/url"
)

func main() {
	// Define command-line flags
	v2rayURL := flag.String("u", "", "V2Ray URL to convert")
	output := flag.String("o", "", "Output JSON file path (optional)")
	console := flag.Bool("c", false, "Show JSON in console")
	port := flag.Int("p", 1080, "Custom port for inbound proxy (default: 1080)")
	dns := flag.String("d", "", "Custom DNS servers (comma-separated, e.g. \"1.1.1.1, 1.0.0.1\")")
	remarks := flag.String("r", "", "Remarks/Comments for the configuration (optional)")

	flag.Parse()

	// Check if URL is provided
	if *v2rayURL == "" {
		fmt.Println("Error: V2Ray URL is required")
		fmt.Println("Usage: ./x2j -u <v2ray_url> [-o output.json] [-c] [-p port] [-d dns_servers] [-r remarks]")
		os.Exit(1)
	}

	// Parse the V2Ray URL and convert to Xray JSON
	config, err := url.ParseV2RayURL(*v2rayURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	// Update the inbound port if a custom port is specified
	if *port != 1080 {
		for i := range config.Inbounds {
			if config.Inbounds[i].Tag == "in_proxy" {
				config.Inbounds[i].Port = *port
				break
			}
		}
	}

	// Update DNS servers if custom DNS is specified
	if *dns != "" {
		if *dns == "\"\"" || *dns == "''" {
			// Empty DNS - clear the servers list
			config.DNS.Servers = []string{}
		} else {
			// Parse comma-separated DNS servers
			dnsServers := strings.Split(*dns, ",")
			// Trim whitespace from each server
			for i, server := range dnsServers {
				dnsServers[i] = strings.TrimSpace(server)
			}
			config.DNS.Servers = dnsServers
		}
	}

	// Set remarks if provided
	if *remarks != "" {
		config.Remarks = *remarks
	}

	// Format the JSON output
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	// Output to console if requested or no output file specified
	if *console || *output == "" {
		fmt.Println(string(jsonData))
	}

	// Save to file if output path is specified
	if *output != "" {
		err = os.WriteFile(*output, jsonData, 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration saved to %s\n", *output)
	}
}