#!/bin/bash

# Exit on any error
set -e

# Function to display menu
show_menu() {
    echo "===================="
    echo "X2J Build Script"
    echo "===================="
    echo "1) Build Android library"
    echo "2) Build iOS framework"
    echo "3) Build both Android and iOS"
    echo "4) Exit"
    echo "===================="
}

# Function to check if mobile packages exist
check_mobile_packages() {
    if [ ! -d "android" ] && [ ! -d "ios" ]; then
        return 1
    fi
    return 0
}

# Function to check Go version
check_go_version() {
    required_version="1.16"
    current_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | grep -oE '[0-9]+\.[0-9]+')
    
    if [ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" != "$required_version" ]; then
        echo "❌ Error: Go version $required_version or higher is required (current: $current_version)"
        return 1
    fi
    return 0
}

# Function to build Android library
build_android() {
    # Check if android directory exists
    if [ ! -d "android" ]; then
        echo "❌ Android directory not found. Skipping Android build."
        return 1
    fi
    
    echo "Cleaning Android build directory..."
    rm -rf build/android
    mkdir -p build/android
    
    echo "Building Android library..."
    # Add -x flag for verbose output and better debugging
    gomobile bind -v -x \
        -target=android \
        -androidapi 35 \
        -ldflags "-w -s" \
        -o build/android/x2j.aar \
        ./android
    
    if [ $? -eq 0 ]; then
        echo "✅ Android library built successfully at: build/android/x2j.aar"
        return 0
    else
        echo "❌ Failed to build Android library"
        return 1
    fi
}

# Function to build iOS framework
build_ios() {
    # Check if running on macOS
    if [ "$(uname)" != "Darwin" ]; then
        echo "❌ iOS builds are only supported on macOS"
        return 1
    }
    
    # Check if ios directory exists
    if [ ! -d "ios" ]; then
        echo "❌ iOS directory not found. Skipping iOS build."
        return 1
    }
    
    # Check if Xcode is installed
    if ! command -v xcodebuild &> /dev/null; then
        echo "❌ Xcode is required for iOS builds but was not found"
        return 1
    fi
    
    echo "Cleaning iOS build directory..."
    rm -rf build/ios
    mkdir -p build/ios
    
    echo "Building iOS framework..."
    # Add -x flag for verbose output and better debugging
    gomobile bind -v -x \
        -target=ios \
        -ldflags "-w -s" \
        -o build/ios/IosX2J.xcframework \
        ./ios
    
    if [ $? -eq 0 ]; then
        echo "✅ iOS framework built successfully at: build/ios/IosX2J.xcframework"
        return 0
    else
        echo "❌ Failed to build iOS framework"
        return 1
    fi
}

# Initialize mobile bind
init_mobile_bind() {
    echo "Initializing mobile bind..."
    
    # Check Go version first
    if ! check_go_version; then
        exit 1
    fi
    
    # Ensure modules are up to date
    echo "Updating Go modules..."
    go mod download
    go mod tidy
    
    # Install/update gomobile
    echo "Installing/updating gomobile..."
    go install golang.org/x/mobile/cmd/gomobile@latest
    
    # Initialize gomobile
    echo "Initializing gomobile..."
    gomobile init
    
    # Verify gomobile installation
    if ! command -v gomobile &> /dev/null; then
        echo "❌ Failed to install gomobile"
        exit 1
    fi
    
    echo "✅ Mobile bind initialized successfully"
}

# Setup Android SDK for CI environments
setup_android_sdk() {
    if [ "$GITHUB_ACTIONS" = "true" ]; then
        echo "Setting up Android SDK for GitHub Actions..."
        # No setup needed here as it's handled in GitHub Actions workflow
        echo "Android SDK setup is handled in GitHub Actions workflow"
    fi
}

# Main script
if [ "$GITHUB_ACTIONS" = "true" ]; then
    # Auto mode for CI/CD: build both Android and iOS, then exit
    echo "CI detected: running builds, then exiting."
    
    # Initialize mobile bind first to ensure gomobile is available
    init_mobile_bind
    
    echo "Starting Android build..."
    build_android
    android_result=$?
    echo ""
    
    # Only attempt iOS build on macOS
    if [ "$(uname)" = "Darwin" ]; then
        echo "Starting iOS build..."
        build_ios
        ios_result=$?
    else
        echo "Skipping iOS build (not on macOS)"
        ios_result=0
    fi
    
    echo ""
    echo "===================="
    echo "Build Summary:"
    if [ $android_result -eq 0 ]; then
        echo "✅ Android: SUCCESS"
    else
        echo "❌ Android: FAILED"
    fi
    
    if [ "$(uname)" = "Darwin" ]; then
        if [ $ios_result -eq 0 ]; then
            echo "✅ iOS: SUCCESS"
        else
            echo "❌ iOS: FAILED"
        fi
    else
        echo "⏭️ iOS: SKIPPED (not on macOS)"
    fi
    echo "===================="
    
    if [ $android_result -ne 0 ] || [ $ios_result -ne 0 ]; then
        exit 1
    fi
    exit 0
fi

# Interactive mode
while true; do
    show_menu
    read -p "Please select an option (1-4): " choice
    
    case $choice in
        1)
            echo "Selected: Build Android library"
            if ! check_mobile_packages; then
                echo "No mobile packages found. Skipping Android build."
            else
                init_mobile_bind
                build_android
            fi
            if [ $? -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        2)
            echo "Selected: Build iOS framework"
            if [ "$(uname)" != "Darwin" ]; then
                echo "❌ iOS builds are only supported on macOS"
                exit 1
            fi
            if ! check_mobile_packages; then
                echo "No mobile packages found. Skipping iOS build."
            else
                init_mobile_bind
                build_ios
            fi
            if [ $? -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        3)
            echo "Selected: Build both Android and iOS"
            if ! check_mobile_packages; then
                echo "No mobile packages found. Skipping mobile builds."
            else
                init_mobile_bind
                
                echo "Starting Android build..."
                build_android
                android_result=$?
                
                echo ""
                if [ "$(uname)" = "Darwin" ]; then
                    echo "Starting iOS build..."
                    build_ios
                    ios_result=$?
                else
                    echo "Skipping iOS build (not on macOS)"
                    ios_result=0
                fi
                
                echo ""
                echo "===================="
                echo "Build Summary:"
                if [ $android_result -eq 0 ]; then
                    echo "✅ Android: SUCCESS"
                else
                    echo "❌ Android: FAILED"
                fi
                
                if [ "$(uname)" = "Darwin" ]; then
                    if [ $ios_result -eq 0 ]; then
                        echo "✅ iOS: SUCCESS"
                    else
                        echo "❌ iOS: FAILED"
                    fi
                else
                    echo "⏭️ iOS: SKIPPED (not on macOS)"
                fi
                echo "===================="
                
                if [ $android_result -ne 0 ] || [ $ios_result -ne 0 ]; then
                    exit 1
                fi
            fi
            echo ""
            ;;
        4)
            echo "Exiting..."
            exit 0
            ;;
        *)
            echo "Invalid option. Please try again."
            ;;
    esac
done