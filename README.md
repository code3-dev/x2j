# x2j - Xray to JSON Converter

A command-line tool written in Go that converts V2Ray share links (VMess, VLess, ShadowSocks, Trojan) to Xray JSON configuration files.

## Features

- Supports all major V2Ray protocols:
  - VMess
  - VLess
  - ShadowSocks
  - Trojan
- Cross-platform compatibility (Windows, macOS, Linux, Android, iOS)
- Multiple output options:
  - Console output
  - JSON file output
  - Combined console and file output
- Customizable inbound port with `-p` flag
- Customizable DNS servers with `-d` flag

## Installation

1. Install Go 1.25.0 or later
2. Clone this repository
3. Build the tool:
   ```bash
   go build -o x2j .
   ```

## Usage

### Basic usage

```bash
# Show JSON in console
./x2j -u <v2ray_link>

# Save to JSON file
./x2j -u <v2ray_link> -o out.json

# Save to JSON file and show in console
./x2j -u <v2ray_link> -o out.json -c

# Set custom inbound port (default is 1080)
./x2j -u <v2ray_link> -p 10809

# Set custom DNS servers
./x2j -u <v2ray_link> -d "1.1.1.1, 1.0.0.1"

# Clear DNS servers (use system DNS)
./x2j -u <v2ray_link> -d ""
```

### Command Examples

```bash
# Save in json file
./x2j -u <v2ray link> -o out.json

# Show json in console without save
./x2j -u <v2ray link>

# Save in json file and show json in console
./x2j -u <v2ray link> -o out.json -c

# Set custom inbound port
./x2j -u <v2ray link> -p 10809

# Set custom DNS servers
./x2j -u <v2ray link> -d "1.1.1.1, 1.0.0.1, 8.8.8.8"

# Clear DNS servers
./x2j -u <v2ray link> -d ""

# Save in json file with custom port
./x2j -u <v2ray link> -o out.json -p 10809

# Save in json file with custom port and show in console
./x2j -u <v2ray link> -o out.json -p 10809 -c

# Save in json file with custom DNS
./x2j -u <v2ray link> -o out.json -d "1.1.1.1, 1.0.0.1"

# Save in json file with custom port and DNS
./x2j -u <v2ray link> -o out.json -p 10809 -d "1.1.1.1, 1.0.0.1"

# Save in json file with absolute path
./x2j -u <v2ray link> -o C:\Users\Dell\Downloads\go\x2j\out.json

# Save in json file with relative path
./x2j -u <v2ray link> -o files/out.json

# Save in json file with parent directory path
./x2j -u <v2ray link> -o ../files/out.json
```

Note: You can rename the executable file to anything for easier usage, for example: `./myapp -u <v2ray link> -o out.json`

### Examples

```bash
# VMess URL to console with custom port
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -p 10809

# VLess URL to file with custom port
./x2j -u "vless://00000000-0000-0000-0000-000000000000@127.0.0.1:8080?security=tls&type=ws&path=/&host=example.com#test" -o vless_config.json -p 10809

# ShadowSocks URL to file and console with custom port
./x2j -u "ss://YWVzLTEyOC1nY206cGFzc3dvcmQ@127.0.0.1:8080#test" -o ss_config.json -p 10809 -c

# Trojan URL to file with custom port
./x2j -u "trojan://password@127.0.0.1:8080?security=tls&type=tcp#test" -o trojan_config.json -p 10809

# VMess URL with custom DNS servers
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "1.1.1.1, 1.0.0.1, 8.8.8.8" -c

# VMess URL with no DNS servers (use system DNS)
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "" -c

# VMess URL with single DNS server
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "1.1.1.1" -c
```

## Supported Protocols

- **VMess**: Full support for all VMess configurations
- **VLess**: Full support including TLS, Reality, and various transport protocols
- **ShadowSocks**: Support for all ShadowSocks methods
- **Trojan**: Full support including TLS configurations

## Supported Transport Protocols

- TCP
- WebSocket (WS)
- HTTP/2 (H2)
- mKCP
- QUIC
- gRPC
- XHTTP
- HTTPUpgrade

## Cross-platform Support

The tool can be compiled for:
- Windows (.exe - x86 and x64)
- Windows (.dll library - x86 and x64)
- macOS (Darwin)
- Linux
- Android (.aar library, .jar library)
- iOS (.xcframework)

## Building for Different Platforms

### Windows Executable
```
# For x64
GOOS=windows GOARCH=amd64 go build -o x2j.exe .

# For x86
GOOS=windows GOARCH=386 go build -o x2j.exe .
```

### Windows DLL
```
# For x64
GOOS=windows GOARCH=amd64 go build -o x2j.dll -buildmode=c-shared .

# For x86
GOOS=windows GOARCH=386 go build -o x2j.dll -buildmode=c-shared .
```

### macOS
```
GOOS=darwin GOARCH=amd64 go build -o x2j .
```

### Linux
```
GOOS=linux GOARCH=amd64 go build -o x2j .
```

### Android
```bash
# Make sure gomobile is installed
# Requires Android SDK to be installed
./mobile.sh
# Select option 1 for Android library (AAR and JAR)
```

### iOS
```bash
# Make sure gomobile is installed
# Requires Xcode to be installed
./mobile.sh
# Select option 2 for iOS framework
```

### Windows DLL (Alternative Method)
```bash
# Using the mobile.sh script
./mobile.sh
# Select option 3 for Windows DLL (will build for the host architecture)
```

## License

This project is licensed under the [GPL v3 License](LICENSE).