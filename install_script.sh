#!/bin/bash

# install-requirements.sh
echo "Installing requirements for blog..."

# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install Go (check latest version at https://golang.org/dl/)
echo "Installing Go..."
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
rm go1.22.0.linux-arm64.tar.gz

# Add Go to PATH if not already there
if ! grep -q "/usr/local/go/bin" $HOME/.profile; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
fi
source $HOME/.profile

# Install Hugo
#echo "Installing Hugo..."
#sudo apt-get install hugo -y

# Create project directory
mkdir -p ~/blogg/api

# Install Go dependencies
echo "Installing Go dependencies..."
cd ~/blogg/api
go mod init blogg
go get -u github.com/joho/godotenv
go get -u golang.org/x/oauth2
go get -u github.com/gorilla/sessions
go get -u google.golang.org/api/oauth2/v2

# Create .env template
echo "Creating .env template..."
#cat > .env.template << EOL
#GOOGLE_CLIENT_ID=your_client_id_here
#GOOGLE_CLIENT_SECRET=your_client_secret_here
#SESSION_KEY=random_32_character_string_here
#EOL

echo "Installation complete!"
echo "Next steps:"
echo "1. Copy your blog files to ~/blogg"
echo "2. Copy your API files to ~/blogg/api"
echo "3. Configure .env with your Google credentials"
echo "4. Update Google Console with new redirect URI: http://your-raspberry-pi-ip:8000/auth/google/callback"