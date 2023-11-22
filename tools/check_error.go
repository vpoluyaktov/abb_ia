package tools

import "abb_ia/internal/logger"

func CheckError(e error) {
	if e != nil {
		logger.Error(e.Error())
		panic(e)
	}
}
