# X2J Android Example

This example demonstrates how to use the X2J library in an Android application to convert V2Ray URLs to JSON configurations.

## Prerequisites

- Android Studio 4.0 or higher
- Android SDK API level 35 or higher
- Go 1.16 or higher
- Gomobile tools

## Setup

1. Install Go and set up your Go environment:
   ```bash
   # Install Go (if not already installed)
   # Visit https://golang.org/dl/ and follow installation instructions
   ```

2. Install Gomobile tools:
   ```bash
   go install golang.org/x/mobile/cmd/gomobile@latest
   gomobile init
   ```

3. Build the X2J Android library:
   ```bash
   # From the project root directory
   ./mobile.sh
   # Select option 1 to build Android library
   ```

   This will generate `x2j.aar` in the `build/android` directory.

## Project Structure

- `AndroidExample.java` - Main activity demonstrating library usage
- `activity_main.xml` - Layout file for the example UI
- `x2j.aar` - The compiled Go library (generated during build)

## Integration Steps

1. Add the AAR dependency to your Android project:
   - Copy `x2j.aar` to your app's `libs` directory
   - Add the following to your app's `build.gradle`:
     ```gradle
     dependencies {
         implementation files('libs/x2j.aar')
     }
     ```

2. Import the Go package in your Java code:
   ```java
   import go.android.Android;
   ```

3. Use the library functions:
   ```java
   // Basic usage
   String jsonConfig = Android.ParseV2RayURL(v2rayURL);

   // With custom settings
   String jsonConfig = Android.ParseV2RayURLWithSettings(v2rayURL, port, dns);
   
   // With custom settings including remarks
   String jsonConfig = Android.ParseV2RayURLWithSettings(v2rayURL, port, dns, remarks);
   ```

## Example Features

The example application demonstrates:
- Converting V2Ray URLs with default settings
- Converting V2Ray URLs with custom port and DNS settings
- Converting V2Ray URLs with remarks/comments
- Proper error handling
- User-friendly UI for input and display
- Input validation

## Building and Running

1. Open the project in Android Studio
2. Sync project with Gradle files
3. Build and run on your device or emulator

## Error Handling

The example shows proper error handling for:
- Invalid V2Ray URLs
- Invalid port numbers
- Missing required fields
- Conversion failures

## Best Practices

- Always validate user input
- Handle all possible exceptions
- Provide clear feedback to users
- Use async operations for long-running tasks
- Follow Android UI guidelines

## License

This example is part of the X2J project and is licensed under the same terms as the main project.