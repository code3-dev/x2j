# Android Integration Guide for X2J Library

This guide explains how to integrate the X2J Android library (built using mobile.sh) into your Android application.

## Library Build Process

The X2J Android library is built as an AAR file using the gomobile tool:
```bash
./mobile.sh
# Select option 1 for Android library
```

This creates `build/android/x2j.aar` which can be integrated into your Android project.

## Integration Steps

1. **Add the AAR to your project:**
   - Copy `x2j.aar` to your app's `libs` directory
   - Add the following to your app's `build.gradle`:
   ```gradle
   dependencies {
       implementation files('libs/x2j.aar')
   }
   ```

2. **Add required dependencies:**
   ```gradle
   dependencies {
       implementation 'com.google.code.gson:gson:2.8.9'
   }
   ```

3. **Sync your project** to import the library

## Usage in Android Code

After building the library with gomobile, you can use it in your Android application as follows:

```java
// Import the generated Go library
import go.android.Android;

// Basic conversion with default settings (port 1080, default DNS)
try {
    String jsonConfig = Android.parseV2RayURL(v2rayURL);
    // Use the JSON configuration
} catch (Exception e) {
    // Handle error
}

// Conversion with custom port and DNS
try {
    String jsonConfig = Android.parseV2RayURLWithSettings(v2rayURL, 1081, "1.1.1.1, 8.8.8.8");
    // Use the JSON configuration
} catch (Exception e) {
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