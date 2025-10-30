#!/bin/bash

# Function to display menu
show_menu() {
    echo "===================="
    echo "X2J Build Script"
    echo "===================="
    echo "1) Build Android library"
    echo "2) Build iOS framework"
    echo "3) Build Windows DLL"
    echo "4) Build both Android and iOS"
    echo "5) Exit"
    echo "===================="
}

# Function to build Android library
build_android() {
    echo "Cleaning Android build directory..."
    rm -rf build/android
    mkdir -p build/android
    
    echo "Building Android library..."
    gomobile bind -v \
        -target=android \
        -androidapi 23 \
        -ldflags "-checklinkname=0" \
        -o build/android/x2j.aar \
        ./x2j/Android
    
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
    echo "Cleaning iOS build directory..."
    rm -rf build/ios
    mkdir -p build/ios
    
    echo "Building iOS framework..."
    gomobile bind -v \
        -target=ios \
        -ldflags "-checklinkname=0" \
        -o build/ios/IosX2J.xcframework \
        ./x2j/Ios
    
    if [ $? -eq 0 ]; then
        echo "✅ iOS framework built successfully at: build/ios/IosX2J.xcframework"
        return 0
    else
        echo "❌ Failed to build iOS framework"
        return 1
    fi
}

# Function to build Windows DLL
build_windows() {
    echo "Cleaning Windows build directory..."
    rm -rf build/windows
    mkdir -p build/windows

    echo "Building Windows DLL..."
    go build -v \
        -o build/windows/x2j.dll \
        -buildmode=c-shared \
        ./x2j/windows

    if [ $? -eq 0 ]; then
        echo "✅ Windows DLL built successfully at: build/windows/x2j.dll"
        return 0
    else
        echo "❌ Failed to build Windows DLL"
        return 1
    fi
}

# Initialize mobile bind
init_mobile_bind() {
    echo "Initializing mobile bind..."
    go get golang.org/x/mobile/bind
}

# Main script
if [ "$GITHUB_ACTIONS" = "true" ]; then
    # Auto mode for CI/CD: build both Android and iOS, then exit
    echo "CI detected: running Android and iOS builds, then exiting."
    init_mobile_bind
    echo "Starting Android build..."
    build_android
    android_result=$?
    echo ""
    echo "Starting iOS build..."
    build_ios
    ios_result=$?
    echo ""
    echo "===================="
    echo "Build Summary:"
    if [ $android_result -eq 0 ]; then
        echo "✅ Android: SUCCESS"
    else
        echo "❌ Android: FAILED"
    fi
    if [ $ios_result -eq 0 ]; then
        echo "✅ iOS: SUCCESS"
    else
        echo "❌ iOS: FAILED"
    fi
    echo "===================="
    if [ $android_result -ne 0 ] || [ $ios_result -ne 0 ]; then
        exit 1
    fi
    exit 0
fi

while true; do
    show_menu
    read -p "Please select an option (1-5): " choice
    
    case $choice in
        1)
            echo "Selected: Build Android library"
            init_mobile_bind
            build_android
            if [ $? -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        2)
            echo "Selected: Build iOS framework"
            init_mobile_bind
            build_ios
            if [ $? -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        3)
            echo "Selected: Build Windows DLL"
            build_windows
            if [ $? -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        4)
            echo "Selected: Build both Android and iOS"
            init_mobile_bind
            
            echo "Starting Android build..."
            build_android
            android_result=$?
            
            echo ""
            echo "Starting iOS build..."
            build_ios
            ios_result=$?
            
            echo ""
            echo "===================="
            echo "Build Summary:"
            if [ $android_result -eq 0 ]; then
                echo "✅ Android: SUCCESS"
            else
                echo "❌ Android: FAILED"
            fi
            
            if [ $ios_result -eq 0 ]; then
                echo "✅ iOS: SUCCESS"
            else
                echo "❌ iOS: FAILED"
            fi
            echo "===================="
            
            if [ $android_result -ne 0 ] || [ $ios_result -ne 0 ]; then
                exit 1
            fi
            echo ""
            ;;
        *)
            echo "Exiting..."
            exit 0
            ;;
    esac
done