package anagram

import (
	"sort"
	"strings"
	"sync"
)

//FindAnagrams ...
func FindAnagrams(words []string) map[string][]string {
	result := newMapConcurrent()
	buckets := processInput(words)

	var wg = &sync.WaitGroup{}
	wg.Add(len(buckets))
	for _, bucket := range buckets {
		go solve(wg, result, bucket)
	}

	wg.Wait()
	return result.data
}

//Removes duplicates; converts words to lower case; sorts to buckets by length;
//Input is considered to be a slice of valid russian words in utf8 encoding
func processInput(words []string) map[int][]string {
	var out = make(map[int][]string)
	var m = make(map[string]struct{})

	for _, v := range words {
		v = strings.ToLower(v)
		if _, has := m[v]; !has {
			m[v] = struct{}{}
			n := len(v)
			out[n] = append(out[n], v)
		}
	}
	return out
}

func solve(wg *sync.WaitGroup, result *mapConcurrent, words []string) {
	defer wg.Done()
	var keyMap = make(map[string]string)
	var out = make(map[string][]string)

	for _, v := range words {
		hash := getAnagramHashFromString(v)
		if _, kmHas := keyMap[hash]; !kmHas {
			keyMap[hash] = v
		}
		out[hash] = append(out[hash], v)
	}

	for anagramHash, properKey := range keyMap {
		values := out[anagramHash]
		if len(values) > 1 {
			sort.Strings(values)
			result.AddValueSlice(properKey, values)
		}
	}
}

//Return unique hash for each anagram family.
//String characters should be lowercase utf8 encoded.
func getAnagramHashFromString(s string) string {
	var b = strings.Builder{}
	r := stringToSortedRuneSlice(s)
	for _, v := range r {
		b.WriteRune(v)
	}
	return b.String()
}

//Sorts strings runes in ascending order
func stringToSortedRuneSlice(s string) []rune {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}
