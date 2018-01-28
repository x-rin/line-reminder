package reminder

import (
	"os"
)

func (con *LineConfig) PostReminder(id string) (string, error) {
	target, err := con.GetProfile(id)
	if err != nil {
		return "", err
	}

	rmdErr := con.PostMessage("To " + target + "\n" + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		return "", rmdErr
	}

	status := SetStatus(id, "false")

	return status, nil
}
