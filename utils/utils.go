package utils

import (
	"log"
)

func FatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func PrintErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func StrFilter(s []string, fn func(string) bool) []string {
	var p []string // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func StrNotEmpty(s string) bool {
	return s != ""
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
