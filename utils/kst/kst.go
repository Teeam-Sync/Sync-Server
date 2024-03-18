package utils_kst

import (
	"time"

	logger "github.com/Teeam-Sync/Sync-Server/logging"
)

var Kst *time.Location

func MustLoadKST() *time.Location {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return loc
}
