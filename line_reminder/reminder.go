package line_reminder

import (
	"os"
)

func (l *LineReminder) PostReminder(id string) (string, error) {
	target, err := l.client.GetProfile(id)
	if err != nil {
		return "", err
	}

	rmdErr := l.client.PostMessage("To " + target + "\n" + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		return "", rmdErr
	}

	status := SetStatus(id, "false")

	return status, nil
}
