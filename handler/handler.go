package handler

import (
	"fmt"
	"net/http"

	"github.com/kutsuzawa/line-reminder/factory"
	"github.com/kutsuzawa/line-reminder/util"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type LineHandler struct {
	logger        *zap.Logger
	reportMessage string
	replyMessage  string
	groupID       string

	serviceFactory factory.ServiceFactory
}

// NewLineHandler - コントローラーを生成
func NewLineHandler(groupID string, serviceFactory factory.ServiceFactory, logger *zap.Logger, reportMessage, replyMessage string) *LineHandler {
	return &LineHandler{
		groupID:        groupID,
		serviceFactory: serviceFactory,
		logger:         logger,
		reportMessage:  reportMessage,
		replyMessage:   replyMessage,
	}
}

func (lc *LineHandler) Report(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("UserID").(string)
	if !ok {
		lc.logger.Error("\"UserID\" is missing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lineService, err := lc.serviceFactory.LineService()
	if err != nil {
		lc.logger.Error("failed to create lineService")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	source, err := lineService.GetNameByID(id)
	if err != nil {
		lc.logger.Error("failed to get source")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("%s\nby %s", lc.reportMessage, source)
	if err := lineService.Send(lc.groupID, msg); err != nil {
		lc.logger.Error("failed to send message to line")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = util.SetStatus(id, "true")
	if err != nil {
		lc.logger.Error("failed to set status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO: なんかダサい
	replyMsg := "えらい！"
	if err := lineService.Send(lc.groupID, replyMsg); err != nil {
		lc.logger.Error("failed to send reply message to line")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (lc *LineHandler) Reply(w http.ResponseWriter, r *http.Request) {
	lineService, err := lc.serviceFactory.LineService()
	if err != nil {
		lc.logger.Error("failed to create lineService")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	event, err := lineService.Hear(r)
	if err != nil {
		lc.logger.Error("failed to get event from line")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message := lc.replyMessage
	word := lc.reportMessage
	msg, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		lc.logger.Error("message type is not TextMessage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if msg.Text == word {
		_, err := util.SetStatus(event.Source.UserID, "true")
		if err != nil {
			lc.logger.Error("failed to set status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := lineService.Reply(event.ReplyToken, message); err != nil {
			lc.logger.Error("failed to reply message to reply")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if msg.Text == "info" {
		groupID := event.Source.GroupID
		userID := event.Source.UserID
		if err := lineService.Reply(event.ReplyToken, fmt.Sprintf("GroupID: %s\n\nUserID: %s", groupID, userID)); err != nil {
			lc.logger.Error("failed to send info")
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
