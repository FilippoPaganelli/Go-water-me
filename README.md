## Build the bot

Build the bot for the target system (here a 64-bit Raspberry Pi):

```bash
GOOS=linux GOARCH=arm64 go build -o go-water-me
```

or simply 

```bash
make build
```

Make it executable:

```bash
chmod +x ./go-water-me
```

## Create a systemd service

Create the systemd service file to describe this service:

```bash
sudo nano /etc/systemd/system/go-water-me.service
```

Paste the following (edit the <fields>):

```bash
[Unit]
Description=Go Telegram Bot Service
After=network.target

[Service]
Type=simple
User=pagans
Group=pagans
ExecStart=<path-to-executable>
WorkingDirectory=<path-to-executable-folder>
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Now, register the service and start it right away:

```bash
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable go-water-me.service
sudo systemctl start go-water-me.service
```

You can always check the service's status by running:

```bash
sudo systemctl status go-water-me.service
```

More logsL

```bash
journalctl -u go-water-me.service -f
```

## Run without service

To just run the bot, type the command:

```bash
make run # must follow the build step
```

## Restart service

After code changes, simply re-build the bot and restart the service:

```bash
make build
make deploy
```

## Env variables

To run this code you should have set up a Telegram bot. The code expects then a `TELEGRAM_BOT_TOKEN` environment variable in a `.env` file in the repo's root folder.