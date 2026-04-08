package apod

import "regexp"

const (
	// PageURL is the canonical "today" permalink.
	PageURL = "https://apod.nasa.gov/apod/astropix.html"

	// BaseURL is prepended to relative image/video paths found in the HTML.
	BaseURL = "https://apod.nasa.gov/apod/"
)

// All regexes are compiled once at startup.
var (
	// Image: <a href="image/..."><img src="image/...">
	reFullImg = regexp.MustCompile(`(?i)<a\s+href="(image/[^"]+)"\s*>\s*<img`)
	reThumb   = regexp.MustCompile(`(?i)<img\s[^>]*src="(image/[^"]+)"`)

	// NASA-hosted video: <a href="image/....mp4"><video ...>
	reNASAVideo = regexp.MustCompile(`(?i)<a\s+href="(image/[^"]+\.mp4)"\s*>\s*<video`)

	// YouTube embed: <iframe src="https://www.youtube.com/embed/ID...">
	reYouTube = regexp.MustCompile(`(?i)<iframe[^>]+src="(https://www\.youtube(?:-nocookie)?\.com/embed/[^"]+)"`)

	// Generic iframe fallback (Vimeo, etc.)
	reIframe = regexp.MustCompile(`(?i)<iframe[^>]+src="(https?://[^"]+)"`)

	// Date: "2026 April 7"
	reDate = regexp.MustCompile(`(\d{4} \w+ \d{1,2})`)

	// Title: first <b>…</b> longer than 10 chars that isn't a known label
	reTitle = regexp.MustCompile(`(?i)<b>\s*([^<]{10,}?)\s*</b>`)

	// Credit block — covers "Image Credit", "Video Credit", "Music Credit", "Credit"
	reCredit = regexp.MustCompile(`(?is)(?:Image|Video|Music)?\s*Credit[^:]*:.*?</b>(.*?)</center>`)

	// Explanation block
	reExplanation = regexp.MustCompile(`(?is)<b>Explanation:</b>\s*(.*?)<p>`)

	// Utility: strip HTML tags, collapse whitespace
	reStripTags = regexp.MustCompile(`<[^>]+>`)
	reSpaces    = regexp.MustCompile(`\s{2,}`)
)

// titleSkipPhrases are substrings that disqualify a <b> tag from being the title.
var titleSkipPhrases = []string{
	"explanation",
	"image credit",
	"video credit",
	"music credit",
	"tomorrow",
	"jigsaw",
	"authors",
	"explore",
}
