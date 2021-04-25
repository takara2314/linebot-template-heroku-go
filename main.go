package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot *linebot.Client
)

const (
	// チャネルシークレットを設定
	YOUR_CHANNEL_SECRET = "XXXXX"
	// チャネルアクセストークンを設定
	YOUR_CHANNEL_ACCESS_TOKEN = "XXXXX"
)

func main() {
	var err error

	bot, err = linebot.New(
		YOUR_CHANNEL_SECRET,
		YOUR_CHANNEL_ACCESS_TOKEN,
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	router := gin.Default()
	// https://example.herokuapp.com/callback にアクセスされたら、
	// callbackPOST という変数を呼び出す
	router.POST("/callback", callbackPOST)

	// ポート番号を環境変数から取得
	var port string = os.Getenv("PORT")
	router.Run(":" + port)
}

// https://example.herokuapp.com/callback にアクセスされたら呼び出される関数
func callbackPOST(c *gin.Context) {
	// LINEから送られてきた情報の前処理
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		// ここにメッセージの内容による処理を書いていこう

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージの種類が「テキスト」なら
			case *linebot.TextMessage:
				var responseMessage string = ""

				// message.Text という変数にメッセージの内容が入っている
				switch message.Text {
				case "おはようございます":
					responseMessage = "Good morning!"
				case "こんにちは":
					responseMessage = "Good afternoon!"
				case "こんばんは":
					responseMessage = "Good evening!"
				default:
					responseMessage = "その言葉はわかりません。"
				}

				// 返信文を送信
				// responseMessage の中に入っている文を返す
				_, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(responseMessage),
				).Do()
				if err != nil {
					log.Println(err)
					panic(err)
				}
			}
		}
	}
}
