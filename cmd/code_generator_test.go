package cmd

import (
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

func TestCodeGeneratorGenerateProducesValidCodeWithDefaultGroupSize(t *testing.T) {
	generator, err := NewCodeGenerator(rand.New(rand.NewSource(1)), 4)
	if err != nil {
		t.Fatalf("create generator: %v", err)
	}

	for range 256 {
		assertValidCode(t, generator.Generate(), 4)
	}
}

func TestCodeGeneratorGenerateProducesValidCodeWithGroupSizeFive(t *testing.T) {
	generator, err := NewCodeGenerator(rand.New(rand.NewSource(2)), 5)
	if err != nil {
		t.Fatalf("create generator: %v", err)
	}

	for range 256 {
		assertValidCode(t, generator.Generate(), 5)
	}
}

func assertValidCode(t *testing.T, code string, groupSize int) {
	t.Helper()

	expectedLength := groupCount*groupSize + groupCount - 1
	if len(code) != expectedLength {
		t.Fatalf("expected code length %d, got %d: %q", expectedLength, len(code), code)
	}

	parts := strings.Split(code, "-")
	if len(parts) != 4 {
		t.Fatalf("expected 4 groups, got %d: %q", len(parts), code)
	}

	allowedLetters := map[rune]struct{}{
		'c': {},
		'u': {},
		'i': {},
		'm': {},
		'g': {},
		'n': {},
		'd': {},
		'a': {},
	}

	seenLetters := make(map[rune]int, len(allowedLetters))

	for _, part := range parts {
		if len(part) != groupSize {
			t.Fatalf("expected group length %d, got %d in %q", groupSize, len(part), code)
		}

		lettersInGroup := 0
		digitsInGroup := 0

		for _, r := range part {
			switch {
			case unicode.IsLetter(r):
				lower := unicode.ToLower(r)
				if _, ok := allowedLetters[lower]; !ok {
					t.Fatalf("unexpected letter %q in %q", r, code)
				}

				seenLetters[lower]++
				if seenLetters[lower] > 1 {
					t.Fatalf("letter %q repeated in %q", lower, code)
				}

				lettersInGroup++
			case unicode.IsDigit(r):
				if r < '2' || r > '9' {
					t.Fatalf("unexpected digit %q in %q", r, code)
				}
				digitsInGroup++
			default:
				t.Fatalf("unexpected character %q in %q", r, code)
			}
		}

		if lettersInGroup != 2 || digitsInGroup != groupSize-2 {
			t.Fatalf("expected each group to contain 2 letters and %d digits: %q", groupSize-2, code)
		}
	}

	if !unicode.IsLetter(rune(parts[0][0])) {
		t.Fatalf("expected the first character to be a letter: %q", code)
	}

	if len(seenLetters) != len(allowedLetters) {
		t.Fatalf("expected all letters to appear exactly once: %q", code)
	}
}
