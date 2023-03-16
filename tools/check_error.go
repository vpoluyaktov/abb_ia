package tools

import "github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"

func CheckError(e error) {
    if e != nil {
		//logger.Error(e.Error())
        panic(e)
    }
}