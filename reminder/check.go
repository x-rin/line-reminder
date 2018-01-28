package reminder

import (
	"os"
)

func (con *LineConfig) Check(id string) (string, error) {
	target, pErr := con.GetProfile(id)
	if pErr != nil {
		return "", pErr
	}

	statusFlag, status, sErr := GetStatus(id)
	if sErr != nil {
		return "", sErr
	}

	if !statusFlag {
		mErr := con.PostMessage(target + ": " + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			return "", mErr
		}
	}

	reStatus := SetStatus(id, status)
	return reStatus, nil
}
