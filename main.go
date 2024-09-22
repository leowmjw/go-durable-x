package main

import (
	"app/telegram" // Local ref using go.work ..
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Test out Telegram App + Bot APIs")
	appID, err := strconv.Atoi(os.Getenv("TELEGRAM_APPID"))
	if err != nil {
		panic(err)
	}
	appHash := os.Getenv("TELEGRAM_APPHASH")
	botToken := os.Getenv("TELEGRAM_BOT_KEY")
	fmt.Println("DUMP:", appID, appHash, botToken)
	//Run(appID, appHash, botToken)

	// Test as Bot ..
	//RunTGBot(botToken)
	// tets new .. as App
	//RunTGProto(appID, appHash)

	telegram.RunGoGram(appID, appHash)
}
