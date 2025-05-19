### Build the bot

Build the bot for the target system:

```bash
GOOS=linux GOARCH=arm64 GOARM=8 go build -o bot
```

Make it executable:

```bash
sudo chmod +x ./bot
```

### Create a systemd service file

Create the systemd service file to describe this service:

```bash
sudo nano /etc/systemd/system/bot.service
```

Paste the following:

```bash
[Unit]
Description=Go Telegram Bot Service
After=network.target

[Service]
Type=simple
User=pagans
Group=pagans
ExecStart=/home/pagans/Github/Go-water-me/bot
WorkingDirectory=/home/pagans/Github/Go-water-me/
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```