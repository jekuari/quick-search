package handlers

import (
	"context"
	"net/http"

	"github.com/jekuari/quick-search/middleware"
)

func Setup(ctx context.Context) {
	// contextualize the handlers
	searchHandler := Search(ctx)
	publicHandler := Public(ctx)

	// handle the routes
	http.Handle("/", addMiddleware(ctx, searchHandler))
	http.Handle("/public/*", addMiddleware(ctx, publicHandler))
}

// Adds all the middlewares to the handlers
func addMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return middleware.AntiSpam(
		ctx,
		middleware.NoCache(
			next,
		),
	)
}
