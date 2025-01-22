package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestRandHexString(t *testing.T) {
	const allowedChars = "0123456789abcdef"

	testCases := []int{0, 1, 5, 10, 50, 100}

	for _, testLength := range testCases {
		t.Run(fmt.Sprintf("Length %d", testLength), func(t *testing.T) {
			randomStr := RandString(testLength)

			if len(randomStr) != testLength {
				t.Fatalf("Expected length %d, got %d", testLength, len(randomStr))
			}

			for _, ch := range randomStr {
				if !strings.ContainsRune(allowedChars, ch) {
					t.Fatalf("Invalid character '%c' in string", ch)
				}
			}

			randomStr2 := RandString(testLength)
			if testLength > 0 && randomStr == randomStr2 {
				t.Fatalf("Two random strings are identical: %s and %s. Is it truly random?", randomStr, randomStr2)
			}
		})
	}
}

func BenchmarkRandHexString(b *testing.B) {

	testCases := []int{0, 1, 5, 10, 41, 50, 100}

	for _, testLength := range testCases {
		for i := 0; i < b.N; i++ {
			RandString(testLength)
		}
	}
}
