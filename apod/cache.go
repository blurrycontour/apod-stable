package apod

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// cache holds the marshalled JSON payload for the current UTC day.
type cache struct {
	mu        sync.RWMutex
	payload   []byte
	fetchedAt time.Time
}

var dayCache cache

// Fetch returns the JSON-encoded APOD for today, hitting the live page at most
// once per UTC calendar day and serving subsequent requests from memory.
func Fetch() ([]byte, error) {
	now := time.Now().UTC()

	dayCache.mu.RLock()
	if dayCache.payload != nil && sameDay(dayCache.fetchedAt, now) {
		data := dayCache.payload
		dayCache.mu.RUnlock()
		return data, nil
	}
	dayCache.mu.RUnlock()

	// Cache miss — scrape and repopulate under write lock.
	dayCache.mu.Lock()
	defer dayCache.mu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited.
	if dayCache.payload != nil && sameDay(dayCache.fetchedAt, now) {
		return dayCache.payload, nil
	}

	html, err := fetchHTML(PageURL)
	if err != nil {
		return nil, err
	}

	apod, err := parse(html)
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(apod, "", "  ")
	if err != nil {
		return nil, err
	}

	dayCache.payload = out
	dayCache.fetchedAt = now
	log.Printf("APOD refreshed: [%s] %q (%s)", apod.Type, apod.Title, apod.Date)
	return out, nil
}

// sameDay reports whether two UTC instants fall on the same calendar date.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
