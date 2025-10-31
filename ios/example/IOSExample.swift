import SwiftUI
import IosX2J

struct ContentView: View {
    @State private var v2rayURL: String = "vmess://eyJhZGQiOiJleGFtcGxlLmNvbSIsInBvcnQiOiI0NDMifQ=="
    @State private var port: String = "1080"
    @State private var dns: String = "1.1.1.1, 8.8.8.8"
    @State private var remarks: String = "Example Proxy Configuration"
    @State private var result: String = ""
    @State private var showAlert = false
    @State private var alertMessage = ""
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                Text("V2Ray URL")
                    .font(.headline)
                TextEditor(text: $v2rayURL)
                    .frame(height: 100)
                    .border(Color.gray.opacity(0.2))
                
                Text("Port")
                    .font(.headline)
                TextField("Enter port number", text: $port)
                    .keyboardType(.numberPad)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                Text("DNS Servers")
                    .font(.headline)
                TextField("Enter DNS servers (comma separated)", text: $dns)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                Text("Remarks (Optional)")
                    .font(.headline)
                TextField("Enter remarks/comments for this configuration", text: $remarks)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                Button(action: {
                    convertWithDefaultSettings()
                }) {
                    Text("Convert with Default Settings")
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                }
                
                Button(action: {
                    convertWithCustomSettings()
                }) {
                    Text("Convert with Custom Settings")
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.green)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                }
                
                if !result.isEmpty {
                    Text("Result")
                        .font(.headline)
                    Text(result)
                        .padding()
                        .background(Color.gray.opacity(0.1))
                        .cornerRadius(8)
                }
            }
            .padding()
        }
        .alert(isPresented: $showAlert) {
            Alert(title: Text("Error"),
                  message: Text(alertMessage),
                  dismissButton: .default(Text("OK")))
        }
    }
    
    private func convertWithDefaultSettings() {
        guard !v2rayURL.isEmpty else {
            showError("Please enter a V2Ray URL")
            return
        }
        
        do {
            let jsonConfig = try IosX2J.ParseV2RayURL(v2rayURL)
            result = jsonConfig
            print("Successfully converted V2Ray URL to JSON with default settings")
        } catch {
            handleError("Error converting V2Ray URL", error: error)
        }
    }
    
    private func convertWithCustomSettings() {
        guard !v2rayURL.isEmpty else {
            showError("Please enter a V2Ray URL")
            return
        }
        
        guard !port.isEmpty, let portNumber = Int32(port) else {
            showError("Please enter a valid port number")
            return
        }
        
        guard !dns.isEmpty else {
            showError("Please enter DNS servers")
            return
        }
        
        do {
            let jsonConfig = try IosX2J.ParseV2RayURLWithSettings(v2rayURL, portNumber, dns, remarks)
            result = jsonConfig
            print("Successfully converted V2Ray URL to JSON with custom settings")
        } catch {
            handleError("Error converting V2Ray URL with custom settings", error: error)
        }
    }
    
    private func handleError(_ message: String, error: Error) {
        print("\(message): \(error.localizedDescription)")
        showError("\(message): \(error.localizedDescription)")
    }
    
    private func showError(_ message: String) {
        alertMessage = message
        showAlert = true
    }
}

#if DEBUG
struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
#endif