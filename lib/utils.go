package lib

import (
	"quick_web_golang/config"
	"strconv"
)

func IsDev() bool {
	debug, _ := strconv.ParseBool(config.Get(config.IsDev))
	return debug
}
