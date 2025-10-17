#!/bin/bash

# Install dependencies
sudo apt update
sudo apt install -y redis-tools

# Create data directory
mkdir -p data

# Set permissions
chmod +x bin/*