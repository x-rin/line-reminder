package line_reminder

import (
	"os"
)

func (l *lineReminder) Check(id string) (string, error) {
	statusFlag, status, sErr := GetStatus(id)
	if sErr != nil {
		return "", sErr
	}

	if !statusFlag {
		target, pErr := l.client.GetProfile(id)
		if pErr != nil {
			return "", pErr
		}
		mErr := l.client.PostMessage("To " + target + "\n" + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			return "", mErr
		}
	}

	return status, nil
}
