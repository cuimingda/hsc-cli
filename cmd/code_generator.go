package cmd

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const (
	groupCount      = 4
	groupLength     = 4
	lettersPerGroup = 2
)

var letterPool = []rune("cuimgnda")
var digitPool = []rune("23456789")

type CodeGenerator struct {
	rng *rand.Rand
}

func NewCodeGenerator(rng *rand.Rand) *CodeGenerator {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return &CodeGenerator{rng: rng}
}

func (g *CodeGenerator) Generate() string {
	letters := append([]rune(nil), letterPool...)
	g.rng.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})

	for i, letter := range letters {
		if g.rng.Intn(2) == 0 {
			letters[i] = unicode.ToUpper(letter)
		}
	}

	var builder strings.Builder
	builder.Grow(groupCount*groupLength + groupCount - 1)

	letterIndex := 0
	for groupIndex := range groupCount {
		letterPositions := g.letterPositions(groupIndex)
		for charIndex := range groupLength {
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

func (g *CodeGenerator) letterPositions(groupIndex int) [groupLength]bool {
	var positions [groupLength]bool

	if groupIndex == 0 {
		positions[0] = true
		positions[1+g.rng.Intn(groupLength-1)] = true
		return positions
	}

	indexes := [groupLength]int{0, 1, 2, 3}
	g.rng.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	for _, index := range indexes[:lettersPerGroup] {
		positions[index] = true
	}

	return positions
}
