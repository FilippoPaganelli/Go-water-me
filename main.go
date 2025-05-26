package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	loadEnv()

	// ----- Bluetooth

	const ARDUINO_DEVICE_NAME string = "Go-water-me (peripheral)"
	var arduino bluetooth.ScanResult

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())
	println("INFO: BLE adapter enabled")

	// Start scanning.
	println("INFO: Scanning for BLE peripherals...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == ARDUINO_DEVICE_NAME {
			arduino = device

			err := adapter.StopScan()
			if nil != err {
				panic(err)
			}
		}
	})
	must("start scan", err)

	println("INFO: Found Arduino peripheral!", fmt.Sprintf("\"%s\"", arduino.LocalName()))

	var serviceData = arduino.ServiceData()
	println(serviceData)

	// ----- Telegram bot

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "foo", bot.MatchTypeCommand, fooHandler)
	b.Start(ctx)
}

func loadEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		println("Error loading .env file")
		panic(err)
	}

	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func must(action string, err error) {
	if err != nil {
		panic("ERROR: Failed to " + action + ": " + err.Error())
	}
}

func fooHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Caught *foo*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Say message with `/foo` anywhere, to do something\\!",
		ParseMode: models.ParseModeMarkdown,
	})
}
