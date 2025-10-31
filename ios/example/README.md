# X2J iOS Example

This example demonstrates how to use the X2J library in an iOS application to convert V2Ray URLs to JSON configurations. The example is built using SwiftUI for a modern, declarative UI implementation.

## Prerequisites

- Xcode 13.0 or higher
- iOS 14.0 or higher (for SwiftUI 2.0 features)
- Go 1.16 or higher
- Gomobile tools
- macOS (required for iOS development)

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

3. Build the X2J iOS framework:
   ```bash
   # From the project root directory
   ./mobile.sh
   # Select option 2 to build iOS framework
   ```

   This will generate `IosX2J.xcframework` in the `build/ios` directory.

## Project Structure

- `IOSExample.swift` - SwiftUI view demonstrating library usage
- `IosX2J.xcframework` - The compiled Go library (generated during build)

## Integration Steps

1. Add the XCFramework to your Xcode project:
   - Drag `IosX2J.xcframework` into your Xcode project
   - In your target's settings, under "General" > "Frameworks, Libraries, and Embedded Content":
     - Ensure `IosX2J.xcframework` is listed
     - Set "Embed & Sign" as the embedding option

2. Import the framework in your Swift code:
   ```swift
   import IosX2J
   ```

3. Use the library functions:
   ```swift
   // Basic usage
   let jsonConfig = try IosX2J.ParseV2RayURL(v2rayURL)

   // With custom settings
   let jsonConfig = try IosX2J.ParseV2RayURLWithSettings(v2rayURL, port, dns)
   
   // With custom settings including remarks
   let jsonConfig = try IosX2J.ParseV2RayURLWithSettings(v2rayURL, port, dns, remarks)
   ```

## Example Features

The example application demonstrates:
- Modern SwiftUI implementation
- Converting V2Ray URLs with default settings
- Converting V2Ray URLs with custom port and DNS settings
- Converting V2Ray URLs with remarks/comments
- Proper error handling and user feedback
- Input validation
- Responsive UI design

## Building and Running

1. Open the Xcode project
2. Select your target device or simulator
3. Build and run the application

## Error Handling

The example shows proper error handling for:
- Invalid V2Ray URLs
- Invalid port numbers
- Missing required fields
- Conversion failures
- Network-related errors

## SwiftUI Implementation Details

The example uses modern SwiftUI features:
- `@State` for managing view state
- Custom styling and layout
- Alert presentation
- Form validation
- Error presentation
- Preview support

## Best Practices

- Validate all user inputs
- Provide clear error messages
- Use Swift's error handling mechanisms
- Follow iOS Human Interface Guidelines
- Implement proper state management
- Use SwiftUI previews for development

## Troubleshooting

Common issues and solutions:
1. Framework not found:
   - Ensure the framework is properly embedded in your target
   - Clean and rebuild the project

2. Build errors:
   - Verify Xcode and iOS deployment target versions
   - Ensure all required dependencies are installed

3. Runtime errors:
   - Check debug console for detailed error messages
   - Verify input format and requirements

## License

This example is part of the X2J project and is licensed under the same terms as the main project.