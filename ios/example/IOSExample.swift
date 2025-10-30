import Foundation
// Import the generated Go library
// After building with gomobile, this will be the actual import:
import IosX2J

// iOS Example for using the X2J library
class IOSExample {
    private let tag = "X2JExample"
    
    func runExamples() {
        // Example V2Ray URL (replace with a valid one)
        let v2rayURL = "vmess://eyJhZGQiOiIxMzkuMTY1LjE4Ni4xMzciLCJhaWQiOiI2NCIsImhvc3QiOiJ2MnVzMDEuaXN4Lnl0IiwiaWQiOiJjNWM1NmQ4NC0zYjQ5LTRhZTktYjNmYS00OWY2ZWM3YjVlMzQiLCJuZXQiOiJ3cyIsInBhdGgiOiJcL3JheSIsInBvcnQiOiI0NDMiLCJwcyI6Ilx1ZDgzY1x1ZGRlOVx1ZDgzY1x1ZGRlYSBXVyIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiJ0bHMiLCJ0eXBlIjoiIiwidiI6IjIifQ=="
        
        // Example 1: Basic usage with default settings
        convertV2RayURL(v2rayURL: v2rayURL)
        
        // Example 2: With custom settings
        convertV2RayURLWithCustomSettings(v2rayURL: v2rayURL)
    }
    
    private func convertV2RayURL(v2rayURL: String) {
        do {
            // Basic conversion with default settings (port 1080, default DNS)
            let jsonConfig = try IosX2J.ParseV2RayURL(v2rayURL)
            print("\(tag): Generated JSON config: \(jsonConfig)")
            
            // In a real implementation, you would use the result here
            // For example, save to file or pass to V2Ray core
            print("\(tag): Converted V2Ray URL to JSON with default settings")
        } catch {
            print("\(tag): Error converting V2Ray URL: \(error.localizedDescription)")
        }
    }
    
    private func convertV2RayURLWithCustomSettings(v2rayURL: String) {
        do {
            // Conversion with custom port and DNS
            let jsonConfig = try IosX2J.ParseV2RayURLWithSettings(v2rayURL, 1081, "1.1.1.1, 8.8.8.8")
            print("\(tag): Generated JSON config with custom settings: \(jsonConfig)")
            
            print("\(tag): Converted V2Ray URL to JSON with custom port (1081) and DNS (1.1.1.1, 8.8.8.8)")
        } catch {
            print("\(tag): Error converting V2Ray URL with custom settings: \(error.localizedDescription)")
        }
    }
}