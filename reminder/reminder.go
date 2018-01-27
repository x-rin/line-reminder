package reminder

import (
	"os"
)

func PostReminder(id string) (string, error) {
	config := NewLineConfig()
	target, err := config.GetProfile(id)
	if err != nil {
		return "", err
	}

	rmdErr := config.PostMessage(target + ": " + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		return "", rmdErr
	}

	status := SetStatus(id, "false")

	return status, nil
}
