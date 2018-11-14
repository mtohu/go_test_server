package common

import (
	"fmt"
	"gotest/logs"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
		os.Exit(1)
	}
}

func WriteError(err string)  {
	los := logs.NewLogerr(true)
	los.Debug([]byte(err))
}
