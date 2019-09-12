package main

import (
	"strconv"
	"strings"
	"time"
)

func parseUptime(output string) time.Duration {
	parts := strings.Fields(output)
	uptimeStrs := parts[2:]

	var (
		days, hours, minutes int
		uptime               time.Duration
		err                  error
	)

	if strings.HasPrefix(uptimeStrs[1], "day") {
		days, err = strconv.Atoi(uptimeStrs[0])
		if err != nil {
			return uptime
		}
		uptimeStrs = uptimeStrs[2:]
	}

	if uptimeStrs[1] == "min," {
		minutes, err = strconv.Atoi(uptimeStrs[0])
		if err != nil {
			return uptime
		}
	} else {
		hourMinutes := uptimeStrs[0]
		hourMinutes = hourMinutes[:len(hourMinutes)-1]
		parts := strings.Split(hourMinutes, ":")
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return uptime
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return uptime
		}
	}

	uptime = time.Duration(days)*24*time.Hour + time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return uptime
}
