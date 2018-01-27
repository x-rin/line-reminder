package reminder

import (
	"os"
)

func PostReport(id string) (string, error) {
	config := NewLineConfig()
	source, err := config.GetProfile(id)
	if err != nil {
		return "", nil
	}

	rptErr := config.PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if rptErr != nil {
		return "", rptErr
	}

	status := SetStatus(id, "true")

	rpyErr := config.PostMessage(os.Getenv("REPLY_SUCCESS"))
	if rpyErr != nil {
		return "", rpyErr
	}

	return status, nil
}
