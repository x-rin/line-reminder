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
		mErr := con.PostMessage("To " + target + "\n" + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			return "", mErr
		}
	}

	return status, nil
}
