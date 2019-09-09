package main

import (
	"strconv"
	"testing"
)

func TestUptime(t *testing.T) {
	strs := []string{
		" 22:10:06 up 39 min,  0 users,  load average: 3.24, 3.26, 3.06",
		" 22:33:09 up  1:02,  0 users,  load average: 2.98, 3.31, 3.31",
		" 22:20:33 up 620 days, 22:37,  1 user,  load average: 0.03, 0.10, 0.10",
	}

	for i, str := range strs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			uptime := parseUptime(str)
			t.Log(uptime)
		})
	}
}
