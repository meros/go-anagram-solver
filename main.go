// Try to find an anagram to a phrase matching a certain MD5

package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func charCountMap(value string) map[uint8]int {
	var charCount map[uint8]int = make(map[uint8]int)

	for i := 0; i < len(value); i++ {
		charCount[value[i]] += 1
	}

	return charCount
}

func issubset(subset map[uint8]int, set map[uint8]int) bool {
	for key, value := range subset {
		if set[key] < value {
			return false
		}
	}

	return true
}

func subtract(subset map[uint8]int, set map[uint8]int) map[uint8]int {
	var result map[uint8]int = make(map[uint8]int)

	for key, value := range set {
		result[key] = value - subset[key]
	}

	return result
}

var anagram string = "poultry outwits ants"
var anagramCM = charCountMap(anagram)
var resultMD5 string = "4624d200580677270a54ccff86b9610e"

// Check if subresult is a solution, and if not then try all potential words
func tryNextWord(words []string, subresult string) (bool, string) {
	if len(subresult) == len(anagram) {
		// Test md5
		subresultMD5 := md5.New()
		io.WriteString(subresultMD5, subresult)
		return (fmt.Sprintf("%x", subresultMD5.Sum(nil)) == resultMD5), subresult
	}

	// We are not ready to test for solution, limit word collection and try all words
	subresultCM := charCountMap(subresult)
	leftCM := subtract(subresultCM, anagramCM)

	var filteredWords []string = make([]string, 0)
	for _, word := range words {
		if issubset(charCountMap(word), leftCM) &&
			// Special case, if there is only 1 space left when subtracting
			// subresult character count map from anagram character count map
			// then also limit on exact word length to fill rest of potential result
			(leftCM[' '] > 1 || len(word)+len(subresult)+1 == len(anagram)) {
			filteredWords = append(filteredWords, word)
		}
	}

	for _, word := range filteredWords {
		newSubresult := word
		if subresult != "" {
			newSubresult = strings.Join([]string{subresult, word}, " ")
		}

		// Did this result in the solution being found?
		found, result := tryNextWord(filteredWords, newSubresult)
		if found {
			return true, result
		}
	}

	// Solution is not in this branch, return false and let callee try something else
	return false, ""
}

// We want to sort the words by length (longest first, to speed up search)
type LongestFirst []string

func (a LongestFirst) Len() int           { return len(a) }
func (a LongestFirst) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LongestFirst) Less(i, j int) bool { return len(a[i]) > len(a[j]) }

func main() {
	fmt.Println("Looking for solution to 'down the rabbit hole' problem")

	words, err := readLines("wordlist")
	if err != nil {
		fmt.Println("could not read wordlist, make sure you run this program in the same folder as the wordlist file")
		return
	}

	sort.Sort(LongestFirst(words))
	found, result := tryNextWord(words, "")
	if !found {
		fmt.Println("Failed to find result")
	} else {
		fmt.Println("Result:", result)
	}
}
