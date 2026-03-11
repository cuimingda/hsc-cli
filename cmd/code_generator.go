package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	defaultGroupSize = 4
	defaultDigits    = "23456789"
	defaultLetters   = "cuimngda"
	lettersPerGroup  = 2
	groupCount       = 4
	maxDigits        = 10
	minDigits        = 1
	requiredLetters  = groupCount * lettersPerGroup
)

type CodeGenerator struct {
	rng       *rand.Rand
	digits    []rune
	groupSize int
	letters   []rune
}

func NewCodeGenerator(rng *rand.Rand, groupSize int, letters string, digits string) (*CodeGenerator, error) {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if err := validateGroupSize(groupSize); err != nil {
		return nil, err
	}

	normalizedLetters, err := normalizeLetters(letters)
	if err != nil {
		return nil, err
	}

	normalizedDigits, err := normalizeDigits(digits)
	if err != nil {
		return nil, err
	}

	return &CodeGenerator{
		rng:       rng,
		digits:    normalizedDigits,
		groupSize: groupSize,
		letters:   normalizedLetters,
	}, nil
}

func (g *CodeGenerator) Generate() string {
	letters := append([]rune(nil), g.letters...)
	g.rng.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	for i, letter := range letters {
		if g.rng.Intn(2) == 0 {
			letters[i] = unicode.ToUpper(letter)
		}
	}

	var builder strings.Builder
	builder.Grow(groupCount*g.groupSize + groupCount - 1)

	letterIndex := 0
	for groupIndex := range groupCount {
		letterPositions := g.letterPositions(groupIndex)
		for charIndex := range g.groupSize {
			if letterPositions[charIndex] {
				builder.WriteRune(letters[letterIndex])
				letterIndex++
				continue
			}

			builder.WriteRune(g.digits[g.rng.Intn(len(g.digits))])
		}

		if groupIndex < groupCount-1 {
			builder.WriteByte('-')
		}
	}

	return builder.String()
}

func validateGroupSize(groupSize int) error {
	switch groupSize {
	case 4, 5:
		return nil
	default:
		return fmt.Errorf("invalid value %d for --group-size: allowed values are 4 or 5", groupSize)
	}
}

func normalizeLetters(letters string) ([]rune, error) {
	normalized := make([]rune, 0, len(letters))
	seen := make(map[rune]struct{}, len(letters))

	for _, letter := range letters {
		if !unicode.IsLetter(letter) {
			return nil, fmt.Errorf("invalid value %q for --letters: only letters are allowed", letters)
		}

		lower := unicode.ToLower(letter)
		if _, ok := seen[lower]; ok {
			continue
		}

		seen[lower] = struct{}{}
		normalized = append(normalized, lower)
	}

	if len(normalized) < requiredLetters {
		return nil, fmt.Errorf(
			"invalid value %q for --letters: need at least %d unique letters after case-insensitive deduplication",
			letters,
			requiredLetters,
		)
	}

	return normalized, nil
}

func normalizeDigits(digits string) ([]rune, error) {
	digitCount := utf8.RuneCountInString(digits)
	if digitCount < minDigits || digitCount > maxDigits {
		return nil, fmt.Errorf(
			"invalid value %q for --digits: must contain %d to %d digits",
			digits,
			minDigits,
			maxDigits,
		)
	}

	normalized := make([]rune, 0, digitCount)
	seen := make(map[rune]struct{}, digitCount)

	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid value %q for --digits: only digits 0-9 are allowed", digits)
		}

		if _, ok := seen[digit]; ok {
			return nil, fmt.Errorf("invalid value %q for --digits: digits must not repeat", digits)
		}

		seen[digit] = struct{}{}
		normalized = append(normalized, digit)
	}

	return normalized, nil
}

func (g *CodeGenerator) letterPositions(groupIndex int) []bool {
	positions := make([]bool, g.groupSize)

	if groupIndex == 0 {
		positions[0] = true
		positions[1+g.rng.Intn(g.groupSize-1)] = true
		return positions
	}

	indexes := make([]int, g.groupSize)
	for i := range indexes {
		indexes[i] = i
	}

	g.rng.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	for _, index := range indexes[:lettersPerGroup] {
		positions[index] = true
	}

	return positions
}
