package rollup

import (
	"fmt"
	"time"
)

type Interval struct {
	Name   string
	Unit   string
	Window int
}

func ParseInterval(value string) (Interval, error) {
	switch value {
	case "minutes":
		return Interval{Name: value, Unit: "minute", Window: 60}, nil
	case "hours":
		return Interval{Name: value, Unit: "hour", Window: 24}, nil
	case "days":
		return Interval{Name: value, Unit: "day", Window: 30}, nil
	case "months":
		return Interval{Name: value, Unit: "month", Window: 12}, nil
	default:
		return Interval{}, fmt.Errorf("invalid interval %q", value)
	}
}

func (i Interval) Buckets(now time.Time) []time.Time {
	end := truncate(now.UTC(), i.Unit)
	buckets := make([]time.Time, i.Window)
	for index := range buckets {
		offset := index - i.Window + 1
		if i.Unit == "month" {
			buckets[index] = end.AddDate(0, offset, 0)
		} else {
			duration := map[string]time.Duration{
				"minute": time.Minute,
				"hour":   time.Hour,
				"day":    24 * time.Hour,
			}[i.Unit]
			buckets[index] = end.Add(time.Duration(offset) * duration)
		}
	}
	return buckets
}

func (i Interval) Start(now time.Time) time.Time {
	return i.Buckets(now)[0]
}

func truncate(value time.Time, unit string) time.Time {
	switch unit {
	case "minute":
		return value.Truncate(time.Minute)
	case "hour":
		return value.Truncate(time.Hour)
	case "day":
		return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
	case "month":
		return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
	default:
		return value
	}
}
