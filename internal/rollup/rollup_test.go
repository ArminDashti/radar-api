package rollup

import (
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name       string
		wantWindow int
		wantUnit   string
	}{
		{"minutes", 60, "minute"},
		{"hours", 24, "hour"},
		{"days", 30, "day"},
		{"months", 12, "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseInterval(tt.name)
			if err != nil {
				t.Fatalf("ParseInterval() error = %v", err)
			}
			if spec.Window != tt.wantWindow || spec.Unit != tt.wantUnit {
				t.Fatalf("got %+v, want window %d unit %s", spec, tt.wantWindow, tt.wantUnit)
			}
		})
	}
}

func TestBucketsAreOldestFirst(t *testing.T) {
	now := time.Date(2026, 8, 13, 1, 37, 0, 0, time.UTC)
	spec, _ := ParseInterval("hours")
	buckets := spec.Buckets(now)

	if len(buckets) != 24 {
		t.Fatalf("bucket count = %d, want 24", len(buckets))
	}
	if !buckets[0].Before(buckets[len(buckets)-1]) {
		t.Fatal("buckets must be ordered oldest to newest")
	}
	if got := buckets[len(buckets)-1]; !got.Equal(time.Date(2026, 8, 13, 1, 0, 0, 0, time.UTC)) {
		t.Fatalf("last bucket = %s", got)
	}
}
