package sentence_gen

import (
	"math/rand"
	"personal/type-training/words"
	"strings"
	"unicode"
)

func getRandomItemFrom(array []string) string {
	return array[rand.Intn(len(array))]
}

func getWordOfCategory(category rune) string {
	switch category {
	case 'A':
		return getRandomItemFrom(words.Adjectives)
	case 'N':
		return getRandomItemFrom(words.Nouns)
	case 'V':
		return getRandomItemFrom(words.Verbs)
	case 'D':
		return getRandomItemFrom(words.Adverbs)
	case 'T':
		return getRandomItemFrom(words.Articles)
	case 'C':
		return getRandomItemFrom(words.Conjunctions)
	case 'P':
		return getRandomItemFrom(words.Prepositions)
	default:
		return ""
	}
}

func generateSentence(template string) string {
	var sentence []string

	for _, char := range template {
		sentence = append(sentence, getWordOfCategory(char))
	}

	sentence[0] = capitalizeFirstLetter(sentence[0])

	return strings.Join(sentence, " ")
}

func capitalizeFirstLetter(word string) string {
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func pickSentenceStructure() string {
	var structures = []string{"ANV", "NVDAN", "TANV", "ANVDAN", "TANVDAN", "NVCNVDAN", "TANVANV", "TANVCDANV", "TANVCTANVDAN", "NVCANVDPN", "TANVDPN", "ANVPNTDAN"}
	return getRandomItemFrom(structures)
}

func GenerateSentences(n int) string {
	var text []string

	for i := 0; i < n; i++ {
		text = append(text, generateSentence(pickSentenceStructure()))
	}

	result := strings.Join(text, ". ")

	return result + ". "
}
