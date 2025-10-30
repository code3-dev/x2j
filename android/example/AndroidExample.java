package dev.pira.x2j.example;

import android.util.Log;
import androidx.appcompat.app.AppCompatActivity;
import android.os.Bundle;

// Import the generated Go library
// After building with gomobile, this will be the actual import:
import go.android.Android;

public class AndroidExample extends AppCompatActivity {
    private static final String TAG = "X2JExample";
    
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        
        // Example V2Ray URL (replace with a valid one)
        String v2rayURL = "vmess://eyJhZGQiOiIxMzkuMTY1LjE4Ni4xMzciLCJhaWQiOiI2NCIsImhvc3QiOiJ2MnVzMDEuaXN4Lnl0IiwiaWQiOiJjNWM1NmQ4NC0zYjQ5LTRhZTktYjNmYS00OWY2ZWM3YjVlMzQiLCJuZXQiOiJ3cyIsInBhdGgiOiJcL3JheSIsInBvcnQiOiI0NDMiLCJwcyI6Ilx1ZDgzY1x1ZGRlOVx1ZDgzY1x1ZGRlYSBXVyIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiJ0bHMiLCJ0eXBlIjoiIiwidiI6IjIifQ==";
        
        // Example 1: Basic usage with default settings
        convertV2RayURL(v2rayURL);
        
        // Example 2: With custom settings
        convertV2RayURLWithCustomSettings(v2rayURL);
    }
    
    private void convertV2RayURL(String v2rayURL) {
        try {
            // Basic conversion with default settings (port 1080, default DNS)
            String jsonConfig = Android.parseV2RayURL(v2rayURL);
            Log.d(TAG, "Generated JSON config: " + jsonConfig);
            
            // In a real implementation, you would use the result here
            // For example, save to file or pass to V2Ray core
            Log.d(TAG, "Converted V2Ray URL to JSON with default settings");
        } catch (Exception e) {
            Log.e(TAG, "Error converting V2Ray URL: " + e.getMessage());
        }
    }
    
    private void convertV2RayURLWithCustomSettings(String v2rayURL) {
        try {
            // Conversion with custom port and DNS
            String jsonConfig = Android.parseV2RayURLWithSettings(v2rayURL, 1081, "1.1.1.1, 8.8.8.8");
            Log.d(TAG, "Generated JSON config with custom settings: " + jsonConfig);
            
            Log.d(TAG, "Converted V2Ray URL to JSON with custom port (1081) and DNS (1.1.1.1, 8.8.8.8)");
        } catch (Exception e) {
            Log.e(TAG, "Error converting V2Ray URL with custom settings: " + e.getMessage());
        }
    }
}