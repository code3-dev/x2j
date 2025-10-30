# iOS Integration Guide for X2J Library

This guide explains how to integrate the X2J iOS library (built using mobile.sh) into your iOS application.

## Library Build Process

The X2J iOS library is built as an XCFramework using the gomobile tool:
```bash
./mobile.sh
# Select option 2 for iOS framework
```

This creates `build/ios/IosX2J.xcframework` which can be integrated into your iOS project.

## Integration Steps

1. **Add the framework to your project:**
   - Drag and drop `IosX2J.xcframework` into your Xcode project
   - Or add it through Xcode's "Frameworks, Libraries, and Embedded Content" section

2. **Import the module in your Swift files:**
   ```swift
   import IosX2J
   ```

## Usage in Swift Code

After building the library with gomobile, you can use it in your iOS application as follows:

```swift
import IosX2J

// Basic conversion with default settings (port 1080, default DNS)
do {
    let jsonConfig = try IosX2J.ParseV2RayURL(v2rayURL)
    // Use the JSON configuration
} catch {
    // Handle error
}

// Conversion with custom port and DNS
do {
    let jsonConfig = try IosX2J.ParseV2RayURLWithSettings(v2rayURL, 1081, "1.1.1.1, 8.8.8.8")
    // Use the JSON configuration
} catch {
    // Handle error
}
```

## Parameters

- `v2rayURL`: The V2Ray URL to convert (supports VMess, VLESS, Trojan, Shadowsocks)
- `port`: Custom port for inbound proxy (default: 1080)
- `dns`: Custom DNS servers as comma-separated string (default: "1.1.1.1, 1.0.0.1, 8.8.8.8, 8.8.4.4")

## Supported Protocols

- VMess
- VLESS
- Trojan
- Shadowsocks

## Example Output

The library will generate a properly formatted Xray JSON configuration that can be used directly with V2Ray/Xray core.