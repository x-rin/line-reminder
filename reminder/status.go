package reminder

import (
	"os"
	"strconv"
	"strings"
)

// GetStatus - 対象のStatusを取得する
func GetStatus(id string) (bool, error) {
	statusKey := strings.ToUpper(id) + "_STATUS"
	statusStr := os.Getenv(statusKey)
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		return false, err
	}
	return status, nil
}

// SetStatus - 対象のStatusをセットする
func SetStatus(id string, status string) (bool, error) {
	statusKey := strings.ToUpper(id) + "_STATUS"
	os.Setenv(statusKey, status)
	statusStr := os.Getenv(statusKey)
	statusBool, err := strconv.ParseBool(statusStr)
	if err != nil {
		return false, err
	}
	return statusBool, nil
}
