Kerberos.io

# Build

docker build -t blomron/telegram-sidecar .

# Environment variables

CHAT_IDS: List of id's where to send the messages to space seperated (example: "12345678 87654321")
BOT_ID: Telegram bot id (example: "123456")
TOKEN: Telegram token (example: "ABC-DEF1234ghIkl-zyx57W2v1u123ew11")

# Run

docker run -it -p 8090:8090 blomron/telegram-sidecar:latest /bin/sh

# Deploy 

helm secrets upgrade --install kerberos-io-chart ./kerberos-io-chart -f kerberos-io-chart/secrets.yaml
