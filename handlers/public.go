package handlers

import (
	"context"
	"net/http"

	"github.com/jekuari/quick-search/constants"
)

func Public(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, constants.HTML_FOLDER+r.URL.Path[1:])
	}
}
