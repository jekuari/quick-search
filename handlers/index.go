package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jekuari/quick-search/constants"
)

func Index(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastChar := cwd[len(cwd)-1:]
	if lastChar != "/" {
		cwd = fmt.Sprintf("%v/", cwd)
	}

	filePath := fmt.Sprintf("%v%v%v", cwd, constants.HTML_FOLDER, "/index.html")
	http.ServeFile(w, r, filePath)
}
