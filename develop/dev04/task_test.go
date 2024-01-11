package main

import "testing"

func Test_FindAnagrams(t *testing.T) {
	input := []string{
		"пятак", "листок", "пятка", "слиток", "тяпка", "столик", "кот", "ток", "отк", "токио",
	}
	expectedOutput := map[string][]string{
		"кот":    {"отк", "ток"},
		"листок": {"слиток", "столик"},
		"пятак":  {"пятка", "тяпка"},
	}

	currentOutput := FindAnagrams(&input)

	if !mapEquals(*currentOutput, expectedOutput) {
		t.Error("wrong answer")
	}
}

func mapEquals(currentAnagrams map[string][]string, expectedAnagrams map[string][]string) bool {
	if len(currentAnagrams) != len(expectedAnagrams) {
		return false
	}

	for key, value := range currentAnagrams {
		if !sliceEquals(expectedAnagrams[key], value) {
			return false
		}
	}

	return true
}

func sliceEquals(expectedAnagramsKeys []string, values []string) bool {
	if len(expectedAnagramsKeys) != len(values) {
		return false
	}

	for i, value := range expectedAnagramsKeys {
		if values[i] != value {
			return false
		}
	}

	return true
}
