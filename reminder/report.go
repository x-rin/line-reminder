package reminder

import (
	"os"
)

func (con *LineConfig) PostReport(id string) (string, error) {
	source, err := con.GetProfile(id)
	if err != nil {
		return "", nil
	}

	rptErr := con.PostMessage(os.Getenv("REPORT_MESSAGE") + "\nby " + source)
	if rptErr != nil {
		return "", rptErr
	}

	status := SetStatus(id, "true")

	rpyErr := con.PostMessage(os.Getenv("REPLY_SUCCESS"))
	if rpyErr != nil {
		return "", rpyErr
	}

	return status, nil
}
