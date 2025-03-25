package main

import "strings"


func badWordsFilter(post string) string {
	const bleep = "****"
	badWords := map[string]bool {
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	words := strings.Split(post, " ")
	for i, word := range words {
		temp := strings.ToLower(word)
		if badWords[temp] {
			words[i] = bleep
		}
	}
	post = strings.Join(words, " ")
	return post
}