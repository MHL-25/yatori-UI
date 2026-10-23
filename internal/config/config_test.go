package config

import (
	"testing"
	"time"
)

func TestInSelfCheckWindow(t *testing.T) {
	mk := func(sw int, start, end string) User {
		u := User{}
		u.CoursesCustom.SelfCheckSw = sw
		u.CoursesCustom.SelfCheckStartDate = start
		u.CoursesCustom.SelfCheckEndDate = end
		return u
	}
	at := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.Local)
	}
	cases := []struct {
		name string
		user User
		now  time.Time
		want bool
	}{
		{"未启用", mk(0, "2026-10-01", "2026-10-31"), at(2026, 11, 5), false},
		{"区间内", mk(1, "2026-10-01", "2026-10-31"), at(2026, 10, 23), true},
		{"开始当天", mk(1, "2026-10-01", "2026-10-31"), at(2026, 10, 1), true},
		{"结束当天", mk(1, "2026-10-01", "2026-10-31"), at(2026, 10, 31), true},
		{"早于开始", mk(1, "2026-10-01", "2026-10-31"), at(2026, 9, 30), false},
		{"晚于结束", mk(1, "2026-10-01", "2026-10-31"), at(2026, 11, 1), false},
		{"仅开始日期", mk(1, "2026-10-01", ""), at(2026, 12, 1), true},
		{"仅结束日期", mk(1, "", "2026-10-31"), at(2026, 1, 1), true},
		{"全空不自动", mk(1, "", ""), at(2026, 1, 1), false},
		{"格式错误不自动", mk(1, "bad", "2026-10-31"), at(2026, 1, 1), false},
	}
	for _, c := range cases {
		got, _ := InSelfCheckWindow(c.user, c.now)
		if got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}
