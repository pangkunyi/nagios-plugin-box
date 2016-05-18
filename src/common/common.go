package server

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

func NagiosParam(param string) (warning, critical float64, err error) {
	fields := strings.Split(param, "|")
	if len(fields) != 2 {
		err = fmt.Errorf("invalid param:%s", param)
		return
	}
	if warning, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return
	}
	critical, err = strconv.ParseFloat(fields[1], 64)
	return
}

func HighState(states ...int) int {
	highState := 0
	for _, state := range states {
		if state > highState {
			highState = state
		}
	}
	return highState
}

func StateString(state int) string {
	switch state {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	case 3:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}
