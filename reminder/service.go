package reminder

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
	"net/http"
)

// LineService - LineAPIを使用するメソッドを定義
type LineService interface {
	GetTargetName(id string) (string, error)
	Send(message string) error
	Hear(request *http.Request) (linebot.Event, error)
	Reply(token string, message string) error
}

type lineService struct {
	api LineAPI
}

// NewLineService - LineServiceを生成
func NewLineService(api LineAPI) LineService {
	return &lineService{
		api: api,
	}
}

func (ls *lineService) GetTargetName(id string) (string, error) {
	target, err := ls.api.GetProfile(id).Do()
	if err != nil {
		return "", err
	}
	return target.DisplayName, nil
}

func (ls *lineService) Send(message string) error {
	if _, err := ls.api.PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage(message)).Do(); err != nil {
		return err
	}
	return nil
}

func (ls *lineService) Hear(request *http.Request) (linebot.Event, error){
	received, err := ls.api.ParseRequest(request)
	if err != nil {
		return linebot.Event{}, err
	}
	var retEvnet linebot.Event
	for _, event := range received {
		//log.Println("groupId: " + event.Source.GroupID)
		//log.Println("userId: " + event.Source.UserID)
		retEvnet = event
	}
	return retEvnet, nil
}

func (ls *lineService) Reply(token string, message string) error {
	if _, err := ls.api.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return err
	}
	return nil
}

