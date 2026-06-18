package math

import "testing"

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"正数相乘", 3, 4, 12},
		{"负数相乘", -2, 5, -10},
		{"负负得正", -3, -3, 9},
		{"乘以零", 7, 0, 0},
		{"零乘以零", 0, 0, 0},
		{"乘以一", 1, 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Multiply(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
