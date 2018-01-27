package reminder

import (
	"os"
)

func Check(id string) (string, error) {
	config := NewLineConfig()
	target, pErr := config.GetProfile(id)
	if pErr != nil {
		return "", pErr
	}

	statusFlag, status, sErr := GetStatus(id)
	if sErr != nil {
		return "", sErr
	}

	if !statusFlag {
		mErr := config.PostMessage(target + ": " + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			return "", mErr
		}
	}

	return status, nil
}
