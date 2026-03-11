package cmd

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestCodeGeneratorGenerateProducesValidCodeWithDefaultGroupSize(t *testing.T) {
	generator, err := NewCodeGenerator(rand.New(rand.NewSource(1)), 4, defaultLetters)
	if err != nil {
		t.Fatalf("create generator: %v", err)
	}

	for range 256 {
		assertValidCode(t, generator.Generate(), 4, defaultLetters)
	}
}

func TestCodeGeneratorGenerateProducesValidCodeWithGroupSizeFive(t *testing.T) {
	generator, err := NewCodeGenerator(rand.New(rand.NewSource(2)), 5, defaultLetters)
	if err != nil {
		t.Fatalf("create generator: %v", err)
	}

	for range 256 {
		assertValidCode(t, generator.Generate(), 5, defaultLetters)
	}
}

func TestCodeGeneratorGenerateProducesValidCodeWithCustomLetters(t *testing.T) {
	generator, err := NewCodeGenerator(rand.New(rand.NewSource(3)), 4, "AbCdEfGhIj")
	if err != nil {
		t.Fatalf("create generator: %v", err)
	}

	for range 256 {
		assertValidCode(t, generator.Generate(), 4, "AbCdEfGhIj")
	}
}

func TestNormalizeLettersRemovesDuplicatesCaseInsensitively(t *testing.T) {
	got, err := normalizeLetters("AaBbCcDdEeFfGgHhI")
	if err != nil {
		t.Fatalf("normalize letters: %v", err)
	}

	want := []rune("abcdefghi")
	if !slices.Equal(got, want) {
		t.Fatalf("expected %q, got %q", string(want), string(got))
	}
}

func TestNewCodeGeneratorRejectsLettersWithNonLetterCharacters(t *testing.T) {
	_, err := NewCodeGenerator(rand.New(rand.NewSource(4)), 4, "abcd1234")
	if err == nil {
		t.Fatal("expected error for letters containing non-letters")
	}

	if !strings.Contains(err.Error(), "only letters are allowed") {
		t.Fatalf("expected non-letter validation error, got %v", err)
	}
}

func TestNewCodeGeneratorRejectsTooFewUniqueLettersAfterDeduplication(t *testing.T) {
	_, err := NewCodeGenerator(rand.New(rand.NewSource(5)), 4, "AaBbCcDd")
	if err == nil {
		t.Fatal("expected error for too few unique letters")
	}

	if !strings.Contains(err.Error(), "need at least 8 unique letters") {
		t.Fatalf("expected unique letter count validation error, got %v", err)
	}
}

func assertValidCode(t *testing.T, code string, groupSize int, allowedLetters string) {
	t.Helper()

	normalizedLetters, err := normalizeLetters(allowedLetters)
	if err != nil {
		t.Fatalf("normalize allowed letters: %v", err)
	}

	expectedLength := groupCount*groupSize + groupCount - 1
	if utf8.RuneCountInString(code) != expectedLength {
		t.Fatalf("expected code length %d, got %d: %q", expectedLength, utf8.RuneCountInString(code), code)
	}

	parts := strings.Split(code, "-")
	if len(parts) != 4 {
		t.Fatalf("expected 4 groups, got %d: %q", len(parts), code)
	}

	allowedLetterSet := make(map[rune]struct{}, len(normalizedLetters))
	for _, letter := range normalizedLetters {
		allowedLetterSet[letter] = struct{}{}
	}

	seenLetters := make(map[rune]int, requiredLetters)

	for _, part := range parts {
		if utf8.RuneCountInString(part) != groupSize {
			t.Fatalf("expected group length %d, got %d in %q", groupSize, utf8.RuneCountInString(part), code)
		}

		lettersInGroup := 0
		digitsInGroup := 0

		for _, r := range part {
			switch {
			case unicode.IsLetter(r):
				lower := unicode.ToLower(r)
				if _, ok := allowedLetterSet[lower]; !ok {
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

	firstCharacter, _ := utf8.DecodeRuneInString(parts[0])
	if !unicode.IsLetter(firstCharacter) {
		t.Fatalf("expected the first character to be a letter: %q", code)
	}

	if len(seenLetters) != requiredLetters {
		t.Fatalf("expected %d unique letters to appear exactly once: %q", requiredLetters, code)
	}
}
