package utils

import (
	"math"
	"testing"
)

func TestBase62RoundTrip(t *testing.T) {
	tests := []struct {
		number uint64
		code   string
	}{
		{number: 0, code: "a"},
		{number: 1, code: "b"},
		{number: 61, code: "9"},
		{number: 62, code: "ba"},
		{number: 12510, code: "dpW"},
		{number: math.MaxUint64, code: EncodeBase62(math.MaxUint64)},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			if got := EncodeBase62(test.number); got != test.code {
				t.Fatalf("EncodeBase62(%d) = %q, want %q", test.number, got, test.code)
			}

			got, err := DecodeBase62(test.code)
			if err != nil {
				t.Fatalf("DecodeBase62(%q) returned error: %v", test.code, err)
			}
			if got != test.number {
				t.Fatalf("DecodeBase62(%q) = %d, want %d", test.code, got, test.number)
			}
		})
	}
}

func TestDecodeBase62RejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"", "a-", "00000000000000000"} {
		if _, err := DecodeBase62(value); err == nil {
			t.Errorf("DecodeBase62(%q) expected an error", value)
		}
	}
}
