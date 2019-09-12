package main

import (
	"testing"
	"time"
)

func TestUptime(t *testing.T) {
	strs := map[time.Duration]string{
		39 * time.Minute:                                 " 22:10:06 up 39 min,  0 users,  load average: 3.24, 3.26, 3.06",
		1*time.Hour + 2*time.Minute:                      " 22:33:09 up  1:02,  0 users,  load average: 2.98, 3.31, 3.31",
		620*24*time.Hour + 22*time.Hour + 37*time.Minute: " 22:20:33 up 620 days, 22:37,  1 user,  load average: 0.03, 0.10, 0.10",
		2*24*time.Hour + 21*time.Hour + 41*time.Minute:   " 19:12:05 up 2 days, 21:41,  0 users,  load average: 3.48, 3.36, 3.33",
		1*24*time.Hour + 21*time.Hour + 41*time.Minute:   " 19:12:05 up 1 day, 21:41,  0 users,  load average: 3.48, 3.36, 3.33",
	}

	for expected, str := range strs {
		t.Run(expected.String(), func(t *testing.T) {
			uptime := parseUptime(str)
			t.Log(uptime)
			if uptime != expected {
				t.Errorf("Expected: %s  Actual: %s", expected, uptime)
			}
		})
	}
}
