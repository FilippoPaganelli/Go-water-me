build:
	GOOS=linux GOARCH=arm64 go build -o go-water-me
	chmod +x go-water-me

deploy:
	sudo systemctl restart go-water-me.service

run:
	./go-water-me

logs:
	journalctl -u go-water-me.service -f
