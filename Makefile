build:
	podman build --platform linux/arm64 -t telegram-bot .

up:
	podman run -d --name telegram-bot --restart=unless-stopped telegram-bot

down:
	podman stop telegram-bot