package reminder

import (
	"os"
	"strconv"
	"strings"
)

// GetStatus - 対象のStatusを取得する
func GetStatus(id string) (bool, string, error) {
	statusKey := strings.ToUpper(id) + "_STATUS"
	status := os.Getenv(statusKey)
	statusFlag, err := strconv.ParseBool(status)
	if err != nil {
		return false, "", err
	}
	return statusFlag, status, nil
}

// SetStatus - 対象のStatusをセットする
func SetStatus(id string, status string) string {
	statusKey := strings.ToUpper(id) + "_STATUS"
	os.Setenv(statusKey, status)
	changedStatus := os.Getenv(statusKey)
	return changedStatus
}
