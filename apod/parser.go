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
var UserAgent = "apod-scraper/1.0 (+https://github.com/you/apod-server)"

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
	return string(body), nil
}

// parse extracts all APOD fields from the raw HTML of an APOD page.
func parse(html string) (*APOD, error) {
	apod := &APOD{SourceURL: PageURL}

	parseMedia(html, apod)
	parseDate(html, apod)
	parseTitle(html, apod)
	parseCredit(html, apod)
	parseExplanation(html, apod)

	if apod.Title == "" && apod.ImageURL == "" && apod.VideoURL == "" {
		return nil, fmt.Errorf("parse failed: could not extract any content — page structure may have changed")
	}
	return apod, nil
}

// parseMedia detects the media type and fills the relevant URL fields.
func parseMedia(html string, apod *APOD) {
	switch {
	case reNASAVideo.MatchString(html):
		apod.Type = MediaVideo
		if m := reNASAVideo.FindStringSubmatch(html); m != nil {
			apod.VideoURL = BaseURL + m[1]
		}

	case reYouTube.MatchString(html):
		apod.Type = MediaYouTube
		if m := reYouTube.FindStringSubmatch(html); m != nil {
			apod.VideoURL = m[1]
		}

	case reIframe.MatchString(html):
		apod.Type = MediaVideo
		if m := reIframe.FindStringSubmatch(html); m != nil {
			apod.VideoURL = m[1]
		}

	default:
		apod.Type = MediaImage
		if m := reFullImg.FindStringSubmatch(html); m != nil {
			apod.ImageURL = BaseURL + m[1]
		}
		if m := reThumb.FindStringSubmatch(html); m != nil {
			apod.ThumbnailURL = BaseURL + m[1]
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

func parseCredit(html string, apod *APOD) {
	if m := reCredit.FindStringSubmatch(html); m != nil {
		apod.Credit = cleanText(reStripTags.ReplaceAllString(m[1], " "))
	}
}

func parseExplanation(html string, apod *APOD) {
	if m := reExplanation.FindStringSubmatch(html); m != nil {
		apod.Explanation = cleanText(reStripTags.ReplaceAllString(m[1], " "))
	}
}

// isSkippedTitle reports whether a <b> tag content should be ignored as a title candidate.
func isSkippedTitle(candidate string) bool {
	if len(candidate) <= 10 || strings.Contains(candidate, "<") {
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
