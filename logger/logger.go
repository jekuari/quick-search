package logger

import (
	"fmt"
	"strings"
	"time"
)

func Log(val ...any) {
	date := time.Now().UTC().Format(time.RFC3339)

	vals := []string{}
	for _, v := range val {
		vals = append(vals, fmt.Sprintf("%v", v))
	}

	concatenated := strings.Join(vals, "")

	fmt.Printf("%v | %v\n", date, concatenated)
}

func Error(title string, err error) {
	Log("Error: ", title, " | description: ", err)
}
