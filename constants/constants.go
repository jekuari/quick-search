package constants

import (
	"os"
	"strconv"
)

type ContextKey struct {
	value string
}

var (
	ORDER_OF_STATIC_FILES = []string{
		"index.html",
		"index.htm",
	}

	REDIS_SEARCHES_CONTEXT_KEY    = ContextKey{value: "redis_searches"}
	REDIS_RATE_LIMITS_CONTEXT_KEY = ContextKey{value: "redis_rate_limits"}

	GOOGLE_SEARCH_URL           string
	GOOGLE_IM_FEELING_LUCKY_URL string
	PORT                        string
	HTML_FOLDER                 string
	PUBLIC_FOLDER               string
	REDIS_URL                   string
	REDIS_DB_RATE_LIMITS        int
	REDIS_DB_SEARCHES           int

	SEARCH_HOUR_RATE_LIMIT int
)

func LoadEnvVariables() {
	GOOGLE_SEARCH_URL = os.Getenv("GOOGLE_SEARCH_URL")
	if GOOGLE_SEARCH_URL == "" {
		panic("Error loading GOOGLE_SEARCH_URL")
	}
	GOOGLE_IM_FEELING_LUCKY_URL = os.Getenv("GOOGLE_IM_FEELING_LUCKY_URL")
	if GOOGLE_IM_FEELING_LUCKY_URL == "" {
		panic("Error loading GOOGLE_IM_FEELING_LUCKY_URL")
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		panic("Error loading PORT")
	}
	HTML_FOLDER = os.Getenv("HTML_FOLDER")
	if HTML_FOLDER == "" {
		panic("Error loading HTML_FOLDER")
	}
	PUBLIC_FOLDER = os.Getenv("PUBLIC_FOLDER")
	if PUBLIC_FOLDER == "" {
		panic("Error loading PUBLIC_FOLDER")
	}
	REDIS_URL = os.Getenv("REDIS_URL")
	if REDIS_URL == "" {
		panic("Error loading REDIS_URL")
	}

	// load hour rate limit
	hourRateLimit, err := strconv.Atoi(os.Getenv("SEARCH_HOUR_RATE_LIMIT"))
	if err != nil {
		panic("Error loading SEARCH_HOUR_RATE_LIMIT")
	}
	SEARCH_HOUR_RATE_LIMIT = hourRateLimit

	dbRateLimits, err := strconv.Atoi(os.Getenv("REDIS_DB_RATE_LIMITS"))
	if err != nil {
		panic("Error loading REDIS_DB_RATE_LIMITS")
	}
	REDIS_DB_RATE_LIMITS = dbRateLimits

	dbSearches, err := strconv.Atoi(os.Getenv("REDIS_DB_SEARCHES"))
	if err != nil {
		panic("Error loading REDIS_DB_SEARCHES")
	}
	REDIS_DB_SEARCHES = dbSearches
}
