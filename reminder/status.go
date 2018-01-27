package reminder

import (
	"os"
	"strconv"
	"strings"
)

func getStatus(id string) (bool, string, error) {
	statusKey := strings.ToUpper(id) + "_STATUS"
	status := os.Getenv(statusKey)
	statusFlag, err := strconv.ParseBool(status)
	return statusFlag, status, err
}

func setStatus(id string, status string) string {
	statusKey := strings.ToUpper(id) + "_STATUS"
	os.Setenv(statusKey, status)
	changedStatus := os.Getenv(statusKey)
	return changedStatus
}
