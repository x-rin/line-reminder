package handler

import (
	"fmt"
	"net/http"

	"github.com/kutsuzawa/line-reminder/service"
	"github.com/kutsuzawa/line-reminder/util"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type LineHandler struct {
	logger        *zap.Logger
	reportMessage string
	replyMessage  string
	groupID       string
	service       service.LineService
}

// NewLineHandler - コントローラーを生成
func NewLineHandler(groupID string, service service.LineService, logger *zap.Logger, reportMessage, replyMessage string) *LineHandler {
	return &LineHandler{
		groupID:       groupID,
		service:       service,
		logger:        logger,
		reportMessage: reportMessage,
		replyMessage:  replyMessage,
	}
}

func (lc *LineHandler) Report(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("UserID").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	source, err := lc.service.GetNameByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("%s\nby %s", lc.reportMessage, source)
	if err := lc.service.Send(lc.groupID, msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = util.SetStatus(id, "true")
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

func (lc *LineHandler) Reply(w http.ResponseWriter, r *http.Request) {
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

func (lc *LineHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
