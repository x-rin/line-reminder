package message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/golang/go/src/pkg/net/http"
)

func (m *message) Check(c *gin.Context) {
	//defer c.Request.Body.Close()
	//body, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil {
	//	fmt.Println("Read ERROR")
	//}
	//decoded, err := base64.StdEncoding.DecodeString(c.GetHeader("X-Line-Signature"))
	//if err != nil {
	//	fmt.Println("Encode ERROR")
	//}
	//hash := hmac.New(sha256.New, []byte(os.Getenv("CHANNELSECRET")))
	//hash.Write(body)
	//// Compare decoded signature and `hash.Sum(nil)` by using `hmac.Equal`
	//fmt.Println("decoded : %s", decoded)
	//fmt.Println("hash.Sum(nil) : %s", hash.Sum(nil))
	//if hmac.Equal(decoded, hash.Sum(nil)) {
	//	log.Println("TEST LOG")
	//	x, _ := ioutil.ReadAll(c.Request.Body)
	//	fmt.Printf("%s", string(x))
	//}
	//c.JSON(http.StatusOK, c)
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	received, _ := bot.ParseRequest(c.Request)

	for _, event := range received {
		fmt.Println(event.ReplyToken)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
