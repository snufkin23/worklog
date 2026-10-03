package platform

import (
	"testing"
	"time"
)

func TestDayBounds(t *testing.T) {
	loc, err := LoadAppLocation()
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	tests := []struct {
		name      string
		now       time.Time
		wantStart time.Time
	}{
		{
			name:      "just after midnight in Nepal belongs to the new day",
			now:       time.Date(2026, 10, 2, 18, 16, 0, 0, time.UTC), // 00:01 on Oct 3 in Nepal
			wantStart: time.Date(2026, 10, 2, 18, 15, 0, 0, time.UTC),
		},
		{
			name:      "one minute before Nepal midnight is still the old day",
			now:       time.Date(2026, 10, 2, 18, 14, 0, 0, time.UTC), // 23:59 on Oct 2 in Nepal
			wantStart: time.Date(2026, 10, 1, 18, 15, 0, 0, time.UTC),
		},
		{
			name:      "midday",
			now:       time.Date(2026, 10, 2, 6, 0, 0, 0, time.UTC), // 11:45 in Nepal
			wantStart: time.Date(2026, 10, 1, 18, 15, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := DayBounds(tt.now, loc)
			if !start.Equal(tt.wantStart) {
				t.Errorf("start = %v, want %v", start.UTC(), tt.wantStart)
			}
			if want := tt.wantStart.Add(24 * time.Hour); !end.Equal(want) {
				t.Errorf("end = %v, want %v", end.UTC(), want)
			}
		})
	}
}
