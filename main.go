package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"pira/x2j/models"
	"pira/x2j/url"
)

func main() {
	// Define command-line flags
	v2rayURL := flag.String("u", "", "V2Ray URL to convert")
	urlBase64 := flag.String("ub64", "", "Base64-encoded V2Ray URL to convert")
	jsonFile := flag.String("j", "", "JSON file to convert to URL")
	jsonBase64 := flag.String("jb64", "", "Base64-encoded JSON to convert to URL")
	output := flag.String("o", "", "Output JSON file path (optional)")
	console := flag.Bool("c", false, "Show JSON in console")
	port := flag.Int("p", 1080, "Custom port for inbound proxy (default: 1080)")
	dns := flag.String("d", "", "Custom DNS servers (comma-separated, e.g. \"1.1.1.1, 1.0.0.1\")")
	remarks := flag.String("r", "", "Remarks/Comments for the configuration (optional)")

	flag.Parse()

	// Check if either URL, base64 URL, JSON file, or base64 JSON is provided
	if *v2rayURL == "" && *urlBase64 == "" && *jsonFile == "" && *jsonBase64 == "" {
		fmt.Println("Error: Either V2Ray URL, base64-encoded URL, JSON file, or base64-encoded JSON is required")
		fmt.Println("Usage: ./x2j -u <v2ray_url> [-o output.json] [-c] [-p port] [-d dns_servers] [-r remarks]")
		fmt.Println("   Or: ./x2j -ub64 <base64_url> [-o output.json] [-c] [-p port] [-d dns_servers] [-r remarks]")
		fmt.Println("   Or: ./x2j -j <json_file> [-c]")
		fmt.Println("   Or: ./x2j -jb64 <base64_json> [-c]")
		os.Exit(1)
	}

	// Handle base64 URL to JSON conversion
	if *urlBase64 != "" {
		// Decode base64 URL
		urlData, err := base64.StdEncoding.DecodeString(*urlBase64)
		if err != nil {
			fmt.Printf("Error decoding base64 URL: %v\n", err)
			os.Exit(1)
		}
		
		decodedURL := string(urlData)

		// Parse the V2Ray URL and convert to Xray JSON
		config, err := url.ParseV2RayURL(decodedURL)
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
		return
	}

	// Handle base64 JSON to URL conversion
	if *jsonBase64 != "" {
		// Decode base64 JSON
		jsonData, err := base64.StdEncoding.DecodeString(*jsonBase64)
		if err != nil {
			fmt.Printf("Error decoding base64 JSON: %v\n", err)
			os.Exit(1)
		}

		// Parse JSON into V2RayConfig
		var config models.V2RayConfig
		if err := json.Unmarshal(jsonData, &config); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}

		// Convert config to URL
		v2rayURL, err := url.ConfigToV2RayURL(&config)
		if err != nil {
			fmt.Printf("Error converting config to URL: %v\n", err)
			os.Exit(1)
		}

		// Output the URL
		if *console {
			fmt.Println(v2rayURL)
		} else {
			fmt.Println(v2rayURL)
		}
		return
	}

	// Handle JSON file to URL conversion
	if *jsonFile != "" {
		// Read JSON file
		jsonData, err := os.ReadFile(*jsonFile)
		if err != nil {
			fmt.Printf("Error reading JSON file: %v\n", err)
			os.Exit(1)
		}

		// Parse JSON into V2RayConfig
		var config models.V2RayConfig
		if err := json.Unmarshal(jsonData, &config); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}

		// Convert config to URL
		v2rayURL, err := url.ConfigToV2RayURL(&config)
		if err != nil {
			fmt.Printf("Error converting config to URL: %v\n", err)
			os.Exit(1)
		}

		// Output the URL
		if *console {
			fmt.Println(v2rayURL)
		} else {
			fmt.Println(v2rayURL)
		}
		return
	}

	// Handle URL to JSON conversion (existing functionality)
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