package line_reminder

import (
	"os"
	"strconv"
	"strings"
)

func GetStatus(id string) (bool, string, error) {
	statusKey := strings.ToUpper(id) + "_STATUS"
	status := os.Getenv(statusKey)
	statusFlag, err := strconv.ParseBool(status)
	if err != nil {
		return false, "", err
	}
	return statusFlag, status, nil
}

func SetStatus(id string, status string) string {
	statusKey := strings.ToUpper(id) + "_STATUS"
	os.Setenv(statusKey, status)
	changedStatus := os.Getenv(statusKey)
	return changedStatus
}
