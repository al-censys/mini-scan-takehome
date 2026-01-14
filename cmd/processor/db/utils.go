package db

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func fileExists(filename string) (exists bool, err error) {
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		exists = false
		err = nil
		return
	}

	if err != nil {
		return
	}

	if !info.IsDir() {
		exists = true
	}

	return
}

func IAmNotImplemented() (err error) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to retrieve call stack")
	}

	name := runtime.FuncForPC(pc).Name()

	idx := strings.LastIndex(name, "/") + 1

	err = fmt.Errorf(name[idx:] + " not implemented")

	return
}
