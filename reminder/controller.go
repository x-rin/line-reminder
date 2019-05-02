package reminder

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineController struct {
	groupID string
	service LineService
}

// NewLineController - コントローラーを生成
func NewLineController(groupID string, service LineService) *LineController {
	return &LineController{
		groupID: groupID,
		service: service,
	}
}

// Check - 対象のstatusをcheckして、falseの場合メッセージを送信する
func (lc *LineController) Check(id, message string) (string, error) {
	status, err := GetStatus(id)
	if err != nil {
		return "", err
	}
	if !status {
		target, err := lc.service.GetNameByID(id)
		if err != nil {
			return "", err
		}
		// e.g To cappyzawa
		// Good Morning
		msg := fmt.Sprintf("To %s\n%s", target, message)
		if err := lc.service.Send(lc.groupID, msg); err != nil {
			return "", err
		}
	}
	return strconv.FormatBool(status), nil
}

// Remind - 対象にメッセージを送信し、statusをfalseに更新する
func (lc *LineController) Remind(id, message string) (string, error) {
	target, err := lc.service.GetNameByID(id)
	if err != nil {
		return "", err
	}
	status, err := SetStatus(id, "false")
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("To %s\n%s", target, message)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		return "", err
	}
	return strconv.FormatBool(status), nil
}

// Report - リマインダーを実行したことを報告し、statusをtrueに更新する
func (lc *LineController) Report(id, message string) (string, error) {
	source, err := lc.service.GetNameByID(id)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("%s\nby %s", message, source)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		return "", err
	}
	status, err := SetStatus(id, "true")
	if err != nil {
		return "", err
	}
	//TODO: なんかダサい
	replyMsg := "えらい！"
	if err := lc.service.Send(lc.groupID, replyMsg); err != nil {
		return "", err
	}
	return strconv.FormatBool(status), nil
}

// ReplyByWord - 対象の投稿を受け取り、statusをtrueに更新した後に返信する
func (lc *LineController) ReplyByWord(req *http.Request, message, word string) (string, error) {
	event, err := lc.service.Hear(req)
	if err != nil {
		return "", err
	}
	msg, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return "", nil
	}
	var status = "false"
	if msg.Text == word {
		statusBool, err := SetStatus(event.Source.UserID, "true")
		if err != nil {
			return "", err
		}
		status = strconv.FormatBool(statusBool)
		if err := lc.service.Reply(event.ReplyToken, message); err != nil {
			return "", err
		}
	} else if msg.Text == "info" {
		groupID := event.Source.GroupID
		userID := event.Source.UserID
		if err := lc.service.Reply(event.ReplyToken, fmt.Sprintf("GroupID: %s\n\nUserID: %s", groupID, userID)); err != nil {
			return "", err
		}
	}
	return status, nil
}
