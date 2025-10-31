# x2j - Xray to JSON Converter

A command-line tool written in Go that converts V2Ray share links (VMess, VLess, ShadowSocks, Trojan) to Xray JSON configuration files.

## Features

- Supports all major V2Ray protocols:
  - VMess
  - VLess
  - Shadowsocks
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
   ```bash
   git clone https://github.com/code3-dev/x2j.git
   cd x2j
   ```
3. Build the tool:
   ```bash
   go mod tidy
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

# Add remarks/comments to the configuration
./x2j -u <v2ray_link> -r "My Proxy Configuration"
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

# Save in json file with remarks
./x2j -u <v2ray link> -o out.json -r "Office Proxy"

# Save in json file with custom port, DNS and remarks
./x2j -u <v2ray link> -o out.json -p 10809 -d "1.1.1.1, 1.0.0.1" -r "Home Proxy"
```

Note: You can rename the executable file to anything for easier usage, for example: `./myapp -u <v2ray link> -o out.json`

### Examples

```bash
# VMess URL to console with custom port
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -p 10809

# VLess URL to file with custom port
./x2j -u "vless://00000000-0000-0000-0000-000000000000@127.0.0.1:8080?security=tls&type=ws&path=/&host=example.com#test" -o vless_config.json -p 10809

# ShadowSocks URL to file and console with custom port
./x2j -u "ss://YWVzLTEyOC1nY206cGFzc3dvcmQ@127LjAuMC4xOjgwODA#test" -o ss_config.json -p 10809 -c

# Trojan URL to file with custom port
./x2j -u "trojan://password@127.0.0.1:8080?security=tls&type=tcp#test" -o trojan_config.json -p 10809

# VMess URL with custom DNS servers
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "1.1.1.1, 1.0.0.1, 8.8.8.8" -c

# VMess URL with no DNS servers (use system DNS)
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "" -c

# VMess URL with single DNS server
./x2j -u "vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOjAsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwibmV0Ijoid3MiLCJwYXRoIjoiLyIsInBvcnQiOjgwODAsInBzIjoidGVzdCIsInRscyI6IiIsInYiOjIsImFpZCI6MCwidHlwZSI6IiJ9" -d "1.1.1.1" -c
```

## Mobile Platform Integration

For developers integrating x2j into mobile applications, check out our platform-specific integration guides:

- [Android Integration Guide](android/example) - Instructions for integrating the X2J Android library (AAR) into Android applications
- [iOS Integration Guide](ios/example) - Instructions for integrating the X2J iOS framework (XCFramework) into iOS applications

The mobile libraries are built using the [mobile.sh](mobile.sh) script which uses gomobile to create native libraries for each platform. The actual class and function names in the generated libraries follow gomobile's naming conventions based on the Go package and function names.

Note: 
- Go Mobile runs on the same architectures as Go, which currently means ARM, ARM64, 386 and amd64 devices and emulators. Notably, Android on MIPS devices is not yet supported.
- Target=ios requires the host machine to be running macOS.
- As of Go 1.5, only darwin/amd64 works on the iOS simulator.

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
- Windows (x86 and x64)
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

## License

This project is licensed under the [GPL v3 License](LICENSE).