package reminder

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"net/http"
	"os"
)

// LineController - コントローラーができることの定義
type LineController interface {
	Check(id string) (string, error)
	Remind(id string) (string, error)
	Report(id string) (string, error)
	Reply(req *http.Request) (string, error)
}

type lineController struct {
	service LineService
}

// NewLineController - コントローラーを生成
func NewLineController(service LineService) LineController {
	return &lineController{
		service: service,
	}
}

// Check - 対象のstatusをcheckして、falseの場合メッセージを送信する
func (lc *lineController) Check(id string) (string, error) {
	statusFlag, status, err := GetStatus(id)
	if err != nil {
		return "", err
	}
	if !statusFlag {
		target, err := lc.service.GetTargetName(id)
		if err != nil {
			return "", err
		}
		//e.g To cappyzawa
		// Good Morning
		message := "To " + target + "\n" + os.Getenv("CHECKED_MESSAGE")
		if err := lc.service.Send(message); err != nil {
			return "", err
		}
	}
	return status, nil
}

// Remind - 対象にメッセージを送信し、statusをfalseに更新する
func (lc *lineController) Remind(id string) (string, error) {
	target, err := lc.service.GetTargetName(id)
	if err != nil {
		return "", err
	}
	message := "To " + target + "\n" + os.Getenv("REMINDER_MESSAGE")
	if err := lc.service.Send(message); err != nil {
		return "", err
	}
	status := SetStatus(id, "false")
	return status, nil
}

// Report - リマインダーを実行したことを報告し、statusをtrueに更新する
func (lc *lineController) Report(id string) (string, error) {
	source, err := lc.service.GetTargetName(id)
	if err != nil {
		return "", nil
	}
	message := os.Getenv("REPORT_MESSAGE") + "\nby " + source
	if err := lc.service.Send(message); err != nil {
		return "", err
	}
	status := SetStatus(id, "true")
	if err := lc.service.Send(os.Getenv("REPLY_SUCCESS")); err != nil {
		return "", err
	}
	return status, nil
}

// Reply - 対象の投稿を受け取り、statusをtrueに更新した後に返信する
func (lc *lineController) Reply(req *http.Request) (string, error) {
	event, err := lc.service.Hear(req)
	if err != nil {
		return "", err
	}
	msg, err := extractMessage(event)
	if err != nil {
		return "", err
	}
	var status = "false"
	if msg == os.Getenv("REPORT_MESSAGE") {
		status = SetStatus(event.Source.UserID, "true")
		if err := lc.service.Reply(event.ReplyToken, os.Getenv("REPLY_SUCCESS")); err != nil {
			return "", err
		}
	}
	return status, nil
}

func extractMessage(event linebot.Event) (string, error) {
	textMsg := new(linebot.TextMessage)
	byteMsg, _ := event.Message.MarshalJSON()
	if err := json.Unmarshal(byteMsg, textMsg); err != nil {
		return "", err
	}
	return textMsg.Text, nil
}
