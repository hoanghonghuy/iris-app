package shared

import "testing"

func TestNormalizePagination_DefaultAndBounds(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "default values when limit and offset are invalid",
			limit:      0,
			offset:     -1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "clamp limit to max",
			limit:      101,
			offset:     5,
			wantLimit:  100,
			wantOffset: 5,
		},
		{
			name:       "keep valid values",
			limit:      15,
			offset:     10,
			wantLimit:  15,
			wantOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := NormalizePagination(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit {
				t.Fatalf("expected limit %d, got %d", tt.wantLimit, gotLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Fatalf("expected offset %d, got %d", tt.wantOffset, gotOffset)
			}
		})
	}
}
