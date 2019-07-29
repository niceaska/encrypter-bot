package main

import (
	"flag"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"fmt"
	"os"
	"time"
	"strconv"
	"encoding/base64"
	"crypto/sha512"
	"math/rand"
)

var (
	telegramBotToken string
)

func init() {
	flag.StringVar(&telegramBotToken, "token", "", "Telegram Bot Token")
	flag.Parse()

	if telegramBotToken == "" {
		log.Print("-token is required")
		os.Exit(1)
	}
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func GenRandPass(length int, err error) string {
	if (length <= 0 || err != nil) {
		length = 10;
	}
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
        "abcdefghijklmnopqrstuvwxyz" +
        digits + specials
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
        buf[i] = all[rand.Intn(len(all))]
	}
	for i := len(buf) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        buf[i], buf[j] = buf[j], buf[i]
	}
	buf[length / 2] = digits[rand.Intn(len(digits))]
	buf[length - 1] = specials[rand.Intn(len(specials))]
	str := string(buf)
	return str;
}

func main() {
	var reply string
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			command := update.Message.Command()
			arguments := update.Message.CommandArguments()
			switch command {
			case "sha512":
				data := []byte(arguments)
				reply = fmt.Sprintf("%x", sha512.Sum512(data))
			case "sha384":
				data := []byte(arguments)
				reply = fmt.Sprintf("%x", sha512.Sum384(data))
			case "pass":
				length, err := strconv.Atoi(arguments)
				reply = GenRandPass(length, err)
			default:
				reply = "Invalid command"
			}
		} else if update.Message != nil {
			reply = update.Message.Text
			if IsBase64(reply) {
				data, err := base64.StdEncoding.DecodeString(reply)
				if (err != nil) {
					log.Panic(err)
				}
				reply = string(data)
			} else {
				reply = string(base64.StdEncoding.EncodeToString([]byte(reply)))
			}
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
