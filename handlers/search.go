package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/jekuari/quick-search/cache"
	"github.com/jekuari/quick-search/constants"
	"github.com/jekuari/quick-search/logger"
)

func Search(ctx context.Context) http.HandlerFunc {
	redisClient := cache.GetRedisSearchesClient(ctx)
	if redisClient == nil {
		panic("Redis client not found in context")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search := query.Get("q")
		originalSearch := search

		if search == "" {
			Index(w, r)
			return
		}

		if rune(originalSearch[0]) != '!' {
			search = url.QueryEscape(search)
			redirectUrl := fmt.Sprintf("%v%v", constants.GOOGLE_SEARCH_URL, search)
			http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
			return
		}

		cachedUrl, err := redisClient.Get(ctx, originalSearch).Result()
		if err != nil {
			logger.Log("Cached url was not available", err)
		}
		if cachedUrl != "" {
			http.Redirect(w, r, cachedUrl, http.StatusSeeOther)
			return
		}

		search = url.QueryEscape(originalSearch[1:])

		imFeelingLuckyPage := fmt.Sprintf(constants.GOOGLE_IM_FEELING_LUCKY_URL, search)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Prevent following redirects
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Get(imFeelingLuckyPage)
		if err != nil {
			logger.Error("Could not get response from I'm feeling lucky", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// check for redirects
		if res.StatusCode == http.StatusMovedPermanently || res.StatusCode == http.StatusFound {
			logger.Log("Going to a direct redirect")
			// extract the url from the header, which looks like: https://www.google.com/url?q=https://gemini.google.com/
			location := res.Header.Get("Location")

			// extract the url from the query string
			u, err := url.Parse(location)
			if err != nil {
				logger.Error("Could not parse url", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			location = u.Query().Get("q")
			http.Redirect(w, r, location, http.StatusSeeOther)
			redisClient.Set(ctx, originalSearch, location, 7*24*time.Hour)
			return
		}

		content, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("Could not read response from I'm feeling lucky", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		anchorRegex := regexp.MustCompile(`<a href="(.+?)">`)

		match := anchorRegex.FindStringSubmatch(string(content))

		if len(match) < 2 {
			http.Error(w, "No match found", http.StatusInternalServerError)
			return
		}

		// convert search string to url encoded string

		redisClient.Set(ctx, originalSearch, match[1], 7*24*time.Hour)
		http.Redirect(w, r, match[1], http.StatusSeeOther)
	}
}
