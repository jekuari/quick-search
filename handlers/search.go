package handlers

import (
	"context"
	/* 	"crypto/tls" */
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
			logger.Log(("Redirecting to index.html"))
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
			logger.Log("Error getting cached url: ", err)
		}
		if cachedUrl != "" {
			logger.Log("Redirecting to cached url: ", cachedUrl)
			http.Redirect(w, r, cachedUrl, http.StatusSeeOther)
			return
		}

		search = url.QueryEscape(originalSearch[1:])

		logger.Log("search: ", search)

		imFeelingLuckyPage := fmt.Sprintf(constants.GOOGLE_IM_FEELING_LUCKY_URL, search)

		/* tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr} */

		res, err := http.Get(imFeelingLuckyPage)
		if err != nil {
			logger.Log("Error getting response: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Log("Response status: ", res.Status)

		content, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Log("Error reading response: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		logger.Log("Content: ", string(content))

		anchorRegex := regexp.MustCompile(`<a href="(.+?)">`)

		match := anchorRegex.FindStringSubmatch(string(content))

		if len(match) < 2 {
			http.Error(w, "No match found", http.StatusInternalServerError)
			return
		}

		logger.Log("Match: ", match[1])

		// convert search string to url encoded string

		redisClient.Set(ctx, originalSearch, match[1], 30*24*time.Hour)
		http.Redirect(w, r, match[1], http.StatusSeeOther)

		logger.Log("Redirecting to: ", match[1])
	}
}
