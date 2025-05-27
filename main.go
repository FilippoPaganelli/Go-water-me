package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"tinygo.org/x/bluetooth"
)

var (
	adapter                 = bluetooth.DefaultAdapter
	soilMoistureServiceUUID = "e969c779-776f-4979-8eb4-d6250e8ea79b"
	// soilMoistureCharacteristicUUID = "4f6b5586-709d-4b06-94fd-8cbea7c32c28"
)

const ARDUINO_DEVICE_NAME string = "Go-water-me (peripheral)"

func main() {
	// ----- Bluetooth

	// Enable BLE interface
	must("enable BLE stack", adapter.Enable())
	println("INFO: BLE adapter enabled")

	// Start scanning
	found := make(chan bluetooth.ScanResult)
	go scanForArduino(found)
	arduino := <-found

	// Try to connect
	peripheral, err := adapter.Connect(arduino.Address, bluetooth.ConnectionParams{})
	if err != nil {
		println("ERROR: ", err.Error())
		return
	}

	println("INFO: connected to ", peripheral.Address.String())

	services, err := peripheral.DiscoverServices(soilMoistureServiceUUID)
	if err != nil {
		println("ERROR: ", err.Error())
		return
	}
	println(services)

	// ----- Telegram bot

	loadEnv()
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

func scanForArduino(found chan bluetooth.ScanResult) {
	println("INFO: Scanning for BLE peripherals...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == ARDUINO_DEVICE_NAME {

			err := adapter.StopScan()
			if nil != err {
				panic(err)
			}
			println("INFO: Found peripheral!", device.LocalName())

			found <- device
		}
	})
	must("start scan", err)
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
