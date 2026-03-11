package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const (
	defaultGroupSize = 4
	defaultLetters   = "cuimngda"
	lettersPerGroup  = 2
	groupCount       = 4
	requiredLetters  = groupCount * lettersPerGroup
)

var digitPool = []rune("23456789")

type CodeGenerator struct {
	rng       *rand.Rand
	groupSize int
	letters   []rune
}

func NewCodeGenerator(rng *rand.Rand, groupSize int, letters string) (*CodeGenerator, error) {
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

	return &CodeGenerator{
		rng:       rng,
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

			builder.WriteRune(digitPool[g.rng.Intn(len(digitPool))])
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
