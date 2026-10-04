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

func TestObfuscatedIDRoundTrip(t *testing.T) {
	obfuscationKey = []byte("test-secret")

	seen := map[string]bool{}
	for _, id := range []uint64{0, 1, 2, 3, 100, 101, 12510, maxID} {
		code, err := EncodeID(id)
		if err != nil {
			t.Fatalf("EncodeID(%d) returned error: %v", id, err)
		}
		if seen[code] {
			t.Fatalf("EncodeID(%d) = %q collides", id, code)
		}
		seen[code] = true

		got, err := DecodeID(code)
		if err != nil || got != id {
			t.Fatalf("DecodeID(%q) = %d, %v; want %d", code, got, err, id)
		}
	}

	if code, _ := EncodeID(1); code == EncodeBase62(1) {
		t.Fatalf("EncodeID(1) = %q, not obfuscated", code)
	}
	if _, err := EncodeID(maxID + 1); err == nil {
		t.Fatal("EncodeID(maxID+1) expected an error")
	}
}
