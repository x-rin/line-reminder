package line_reminder

import "net/http"

type LineReminder interface {
	PostReport(id string) (string, error)
	Check(id string) (string, error)
	PostReminder(id string) (string, error)
	GetWebHook(req *http.Request) (string, error)
}

type lineReminder struct {
	client LineClient
}

func NewLineReminder(client LineClient) LineReminder {
	return &lineReminder{
		client: client,
	}
}
