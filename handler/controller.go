package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kutsuzawa/line-reminder/service"
	"github.com/kutsuzawa/line-reminder/util"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type LineController struct {
	logger         *zap.Logger
	remindMessage  string
	reportMessage  string
	replyMessage   string
	checkedMessage string
	groupID        string
	service        service.LineService
	id             string
}

// NewLineController - コントローラーを生成
func NewLineController(groupID string, service service.LineService, logger *zap.Logger, remindMessage, reportMessage, replyMessage, checkedMessage string) *LineController {
	return &LineController{
		groupID:        groupID,
		service:        service,
		logger:         logger,
		remindMessage:  remindMessage,
		reportMessage:  reportMessage,
		replyMessage:   replyMessage,
		checkedMessage: checkedMessage,
	}
}

func (lc *LineController) getID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lc.id = r.URL.Query().Get("id")
		next.ServeHTTP(w, r)
	}
}

// Check - 対象のstatusをcheckして、falseの場合メッセージを送信する
func (lc *LineController) Check(id, message string) (string, error) {
	status, err := util.GetStatus(id)
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

func (lc *LineController) check(w http.ResponseWriter, r *http.Request) {
	//TODO: id が空のときの処理を書いてあげたほう丁寧
	status, err := util.GetStatus(lc.id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !status {
		target, err := lc.service.GetNameByID(lc.id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// e.g To cappyzawa
		// Good Morning
		msg := fmt.Sprintf("To %s\n%s", target, lc.checkedMessage)
		if err := lc.service.Send(lc.groupID, msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// Remind - 対象にメッセージを送信し、statusをfalseに更新する
func (lc *LineController) Remind(id, message string) (string, error) {
	target, err := lc.service.GetNameByID(id)
	if err != nil {
		return "", err
	}
	status, err := util.SetStatus(id, "false")
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("To %s\n%s", target, message)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		return "", err
	}
	return strconv.FormatBool(status), nil
}

func (lc *LineController) remind(w http.ResponseWriter, r *http.Request) {
	target, err := lc.service.GetNameByID(lc.id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = util.SetStatus(lc.id, "false")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("To %s\n%s", target, lc.remindMessage)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

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
	status, err := util.SetStatus(id, "true")
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

func (lc *LineController) report(w http.ResponseWriter, r *http.Request) {
	source, err := lc.service.GetNameByID(lc.id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("%s\nby %s", lc.reportMessage, source)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = util.SetStatus(lc.id, "true")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO: なんかダサい
	replyMsg := "えらい！"
	if err := lc.service.Send(lc.groupID, replyMsg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

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
		statusBool, err := util.SetStatus(event.Source.UserID, "true")
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

func (lc *LineController) reply(w http.ResponseWriter, r *http.Request) {
	event, err := lc.service.Hear(r)
	message := lc.replyMessage
	word := lc.reportMessage
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if msg.Text == word {
		_, err := util.SetStatus(event.Source.UserID, "true")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := lc.service.Reply(event.ReplyToken, message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if msg.Text == "info" {
		groupID := event.Source.GroupID
		userID := event.Source.UserID
		if err := lc.service.Reply(event.ReplyToken, fmt.Sprintf("GroupID: %s\n\nUserID: %s", groupID, userID)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)

}

type logger interface {
	// TODO: ださい
	Sync()
}

func (lc *LineController) Run(port string) {
	endpointPrefix := "api/v1"
	http.HandleFunc(endpointPrefix+"/check", lc.getID(lc.check))
	http.HandleFunc(endpointPrefix+"/remind", lc.getID(lc.remind))
	http.HandleFunc(endpointPrefix+"/report", lc.getID(lc.report))
	http.HandleFunc(endpointPrefix+"/reply", lc.getID(lc.reply))

	http.ListenAndServe(":"+port, nil)
}
