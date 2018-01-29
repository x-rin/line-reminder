package line_reminder

import (
	"os"
)

func (l *LineReminder) PostReport(id string) (string, error) {
	source, err := l.client.GetProfile(id)
	if err != nil {
		return "", nil
	}

	rptErr := l.client.PostMessage(os.Getenv("REPORT_MESSAGE") + "\nby " + source)
	if rptErr != nil {
		return "", rptErr
	}

	status := SetStatus(id, "true")

	rpyErr := l.client.PostMessage(os.Getenv("REPLY_SUCCESS"))
	if rpyErr != nil {
		return "", rpyErr
	}

	return status, nil
}
