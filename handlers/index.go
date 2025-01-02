package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jekuari/quick-search/.git/logger"
	"github.com/jekuari/quick-search/constants"
)

func Index(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := fmt.Sprintf("%v/%v%v", cwd, constants.HTML_FOLDER, "/index.html")
	logger.Log("Serving file: ", filePath)
	http.ServeFile(w, r, filePath)
}
