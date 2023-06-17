package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Run server at 8080")

	http.HandleFunc("/api/substring", LongSubstringHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func LongSubstringHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "text/plain" {
		http.Error(w, "Only text/plain content are allowed", http.StatusUnsupportedMediaType)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(reqBody) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	respData := getSubstring(string(reqBody))

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respData))
}

func getSubstring(s string) string {
	if len(s) == 0 {
		return ""
	}

	subL, subR := 0, 0 // substring edge indexes
	l, maxLen := 0, 0
	charSet := make(map[byte]bool)
	for r := range s {
		// delete up to the repeated character.
		for charSet[s[r]] {
			delete(charSet, s[l])
			l++
		}
		charSet[s[r]] = true
		// find max substring.
		if maxLen < r-l+1 {
			subL, subR = l, r
			maxLen = r - l + 1
		}
	}

	var sb strings.Builder
	for i := subL; i <= subR; i++ {
		sb.WriteByte(s[i])
	}
	return sb.String()
}
