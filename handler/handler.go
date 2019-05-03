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
	logger         *zap.Logger
	remindMessage  string
	reportMessage  string
	replyMessage   string
	checkedMessage string
	groupID        string
	service        service.LineService
	id             string
}

// NewLineHandler - コントローラーを生成
func NewLineHandler(groupID string, service service.LineService, logger *zap.Logger, remindMessage, reportMessage, replyMessage, checkedMessage string) *LineHandler {
	return &LineHandler{
		groupID:        groupID,
		service:        service,
		logger:         logger,
		remindMessage:  remindMessage,
		reportMessage:  reportMessage,
		replyMessage:   replyMessage,
		checkedMessage: checkedMessage,
	}
}

func (lc *LineHandler) getID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lc.id = r.URL.Query().Get("id")
		next.ServeHTTP(w, r)
	}
}

func (lc *LineHandler) check(w http.ResponseWriter, r *http.Request) {
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

func (lc *LineHandler) remind(w http.ResponseWriter, r *http.Request) {
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

func (lc *LineHandler) report(w http.ResponseWriter, r *http.Request) {
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

func (lc *LineHandler) reply(w http.ResponseWriter, r *http.Request) {
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

func (lc *LineHandler) Run(port string) error {
	endpointPrefix := "/api/v1"
	http.HandleFunc(endpointPrefix+"/check", lc.getID(lc.check))
	http.HandleFunc(endpointPrefix+"/remind", lc.getID(lc.remind))
	http.HandleFunc(endpointPrefix+"/report", lc.getID(lc.report))
	http.HandleFunc(endpointPrefix+"/webhook", lc.getID(lc.reply))

	return http.ListenAndServe(":"+port, nil)
}
