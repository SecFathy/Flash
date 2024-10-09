#!/bin/bash

APP_NAME="flash"  # Your Go application name

# List of platforms to build for
PLATFORMS=("windows/amd64" "windows/386" "linux/amd64" "linux/arm64" "linux/arm" "darwin/amd64" "darwin/arm64")

# Create the output directory if it doesn't exist
OUTPUT_DIR="builds"
mkdir -p $OUTPUT_DIR

# Loop through each platform and build
for PLATFORM in "${PLATFORMS[@]}"
do
    # Split platform into OS and architecture
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)
    
    # Set output file name
    OUTPUT_NAME=$OUTPUT_DIR/${APP_NAME}-${OS}-${ARCH}
    
    # Add .exe extension for Windows builds
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi

    echo "Building for $OS/$ARCH..."
    
    # Compile the application for the target OS and architecture
    env GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_NAME

    if [ $? -ne 0 ]; then
        echo "An error occurred while building for $OS/$ARCH."
        exit 1
    fi
done

echo "Build finished for all platforms."
