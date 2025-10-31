package dev.pira.x2j.example;

import android.os.Bundle;
import android.util.Log;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;
import androidx.appcompat.app.AppCompatActivity;

// Import the generated Go library
import go.android.Android;

public class AndroidExample extends AppCompatActivity {
    private static final String TAG = "X2JExample";
    
    private EditText urlInput;
    private EditText portInput;
    private EditText dnsInput;
    private EditText remarksInput;
    private TextView resultText;
    private Button convertButton;
    private Button convertWithSettingsButton;
    
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        
        // Initialize UI components
        initializeViews();
        setupListeners();
        
        // Set default values
        urlInput.setText("vmess://eyJhZGQiOiJleGFtcGxlLmNvbSIsInBvcnQiOiI0NDMifQ=="); // Example URL
        portInput.setText("1080");
        dnsInput.setText("1.1.1.1, 8.8.8.8");
        remarksInput.setText("Example Proxy Configuration");
    }
    
    private void initializeViews() {
        urlInput = findViewById(R.id.urlInput);
        portInput = findViewById(R.id.portInput);
        dnsInput = findViewById(R.id.dnsInput);
        remarksInput = findViewById(R.id.remarksInput);
        resultText = findViewById(R.id.resultText);
        convertButton = findViewById(R.id.convertButton);
        convertWithSettingsButton = findViewById(R.id.convertWithSettingsButton);
    }
    
    private void setupListeners() {
        convertButton.setOnClickListener(v -> {
            String url = urlInput.getText().toString().trim();
            if (url.isEmpty()) {
                showError("Please enter a V2Ray URL");
                return;
            }
            convertV2RayURL(url);
        });
        
        convertWithSettingsButton.setOnClickListener(v -> {
            String url = urlInput.getText().toString().trim();
            String portStr = portInput.getText().toString().trim();
            String dns = dnsInput.getText().toString().trim();
            String remarks = remarksInput.getText().toString().trim();
            
            if (url.isEmpty() || portStr.isEmpty() || dns.isEmpty()) {
                showError("Please fill in all fields");
                return;
            }
            
            try {
                int port = Integer.parseInt(portStr);
                convertV2RayURLWithCustomSettings(url, port, dns, remarks);
            } catch (NumberFormatException e) {
                showError("Invalid port number");
            }
        });
    }
    
    private void convertV2RayURL(String v2rayURL) {
        try {
            // Basic conversion with default settings
            String jsonConfig = Android.ParseV2RayURL(v2rayURL);
            displayResult(jsonConfig);
            Log.d(TAG, "Successfully converted V2Ray URL to JSON with default settings");
        } catch (Exception e) {
            handleError("Error converting V2Ray URL", e);
        }
    }
    
    private void convertV2RayURLWithCustomSettings(String v2rayURL, int port, String dns, String remarks) {
        try {
            // Conversion with custom settings
            String jsonConfig = Android.ParseV2RayURLWithSettings(v2rayURL, port, dns, remarks);
            displayResult(jsonConfig);
            Log.d(TAG, "Successfully converted V2Ray URL to JSON with custom settings");
        } catch (Exception e) {
            handleError("Error converting V2Ray URL with custom settings", e);
        }
    }
    
    private void displayResult(String jsonConfig) {
        // Update UI with the result
        resultText.setText(jsonConfig);
        showToast("Conversion successful!");
    }
    
    private void handleError(String message, Exception e) {
        Log.e(TAG, message + ": " + e.getMessage());
        showError(message);
    }
    
    private void showError(String message) {
        showToast("Error: " + message);
    }
    
    private void showToast(String message) {
        Toast.makeText(this, message, Toast.LENGTH_SHORT).show();
    }
}