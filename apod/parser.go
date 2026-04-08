package apod

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// UserAgent is sent with every request to apod.nasa.gov.
// Override via the USER_AGENT environment variable before calling Fetch.
var UserAgent = "apod-stable"

// fetchHTML retrieves the raw HTML from url.
func fetchHTML(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body: %w", err)
	}

	htmlStr := string(body)
	// Remove various BOMs and handle encoding issues
	htmlStr = strings.TrimPrefix(htmlStr, "\ufeff")       // UTF-8 BOM
	htmlStr = strings.TrimPrefix(htmlStr, "\xef\xbb\xbf") // UTF-8 BOM alt

	// If we see null bytes, likely UTF-16, try to decode
	if strings.Contains(htmlStr, "\x00") {
		// Try UTF-16 decoding - for now just remove null bytes as a workaround
		htmlStr = strings.ReplaceAll(htmlStr, "\x00", "")
	}

	// Debug: show first 500 chars if verbose logging needed
	// log.Printf("DEBUG fetchHTML: charset=%s", resp.Header.Get("Content-Type"))
	return htmlStr, nil
}

// parse extracts all APOD fields from the raw HTML of an APOD page.
func parse(html string) (*APOD, error) {
	apod := &APOD{Source: PageURL}

	parseMedia(html, apod)
	parseDate(html, apod)
	parseTitle(html, apod)
	parseExplanation(html, apod)

	// Debug logging
	// log.Printf("DEBUG parse: type=%s date=%q title=%q credit=%q", apod.Type, apod.Date, apod.Title, apod.Credit)

	// if apod.Title == "" && apod.ImageURL == "" && apod.VideoURL == "" {
	// 	return nil, fmt.Errorf("parse failed: could not extract any content — page structure may have changed")
	// }
	return apod, nil
}

// parseMedia detects the media type and fills the relevant URL fields.
func parseMedia(html string, apod *APOD) {
	switch {
	// Check for <video> tag with <source> element (HTML5 video)
	case reVideoTag.MatchString(html):
		apod.Type = MediaVideo
		if m := reVideoSource.FindStringSubmatch(html); m != nil {
			url := m[1]
			if !strings.HasPrefix(url, "http") {
				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				url = BaseURL + strings.TrimPrefix(url, "/")
			}
			apod.URL = url
		}

	case reNASAVideo.MatchString(html):
		apod.Type = MediaVideo
		if m := reNASAVideo.FindStringSubmatch(html); m != nil {
			apod.URL = BaseURL + m[1]
		}

	case reYouTube.MatchString(html):
		apod.Type = MediaYouTube
		if m := reYouTube.FindStringSubmatch(html); m != nil {
			apod.URL = m[1]
		}

	case reIframe.MatchString(html):
		apod.Type = MediaVideo
		if m := reIframe.FindStringSubmatch(html); m != nil {
			apod.URL = m[1]
		}

	default:
		apod.Type = MediaImage
		if m := reFullImg.FindStringSubmatch(html); m != nil {
			url := m[1]
			// If URL doesn't start with http, prepend BaseURL
			if !strings.HasPrefix(url, "http") {
				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				url = BaseURL + strings.TrimPrefix(url, "/")
			}
			apod.URL = url
		}
	}
}

func parseDate(html string, apod *APOD) {
	if m := reDate.FindStringSubmatch(html); m != nil {
		apod.Date = m[1]
	}
}

func parseTitle(html string, apod *APOD) {
	for _, m := range reTitle.FindAllStringSubmatch(html, -1) {
		candidate := strings.TrimSpace(m[1])
		if isSkippedTitle(candidate) {
			continue
		}
		apod.Title = candidate
		return
	}
}

func parseExplanation(html string, apod *APOD) {
	if m := reExplanation.FindStringSubmatch(html); m != nil {
		apod.Explanation = cleanText(reStripTags.ReplaceAllString(m[1], " "))
	}
}

// isSkippedTitle reports whether a <b> tag content should be ignored as a title candidate.
func isSkippedTitle(candidate string) bool {
	if len(candidate) <= 2 || strings.Contains(candidate, "<") {
		return true
	}
	lower := strings.ToLower(candidate)
	for _, phrase := range titleSkipPhrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

// cleanText strips excess whitespace and newlines from scraped text.
func cleanText(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = reSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
