package apod

import "regexp"

// PageURL is the canonical "today" permalink.
// This can be overridden via the APOD_URL environment variable.
var PageURL = "https://apod.nasa.gov/apod/astropix.html"

// BaseURL is prepended to relative image/video paths found in the HTML.
const BaseURL = "https://apod.nasa.gov/apod/"

// All regexes are compiled once at startup.
var (
	// Image: href="image/..." - just match the href, looks for any image tag nearby
	reFullImg = regexp.MustCompile(`(?i)href="(image/[^"]+\.jpg)"`)
	reThumb   = regexp.MustCompile(`(?i)src="(image/[^"]+\.jpg)"`)

	// Video detection: <video> tag or <source> tag with .mp4
	reVideoTag    = regexp.MustCompile(`(?i)<video`)
	reVideoSource = regexp.MustCompile(`(?i)<source\s+src="(image/[^"]+\.mp4)"`)

	// NASA-hosted video: <a href="...mp4"><video ...>
	reNASAVideo = regexp.MustCompile(`(?i)<a\s+href="([^"]+\.mp4)"\s*>\s*<video`)

	// YouTube embed: <iframe src="https://www.youtube.com/embed/ID...">
	reYouTube = regexp.MustCompile(`(?i)<iframe[^>]+src="(https://www\.youtube(?:-nocookie)?\.com/embed/[^"]+)"`)

	// Generic iframe fallback (Vimeo, etc.)
	reIframe = regexp.MustCompile(`(?i)<iframe[^>]+src="(https?://[^"]+)"`)

	// Date: "2026 April 8" anywhere in the document
	reDate = regexp.MustCompile(`(\d{4}\s+\w+\s+\d{1,2})`)

	// Title: content in <b> tags (will filter with isSkippedTitle)
	reTitle = regexp.MustCompile(`(?i)<b>\s*([^<]+?)\s*</b>`)

	// Credit: matches "Image Credit:" or "Video Credit:" or "Music Credit:" followed by content
	// Content may be plain text or wrapped in HTML tags (like <a> links)
	// Stop at next <b> tag or closing <center>
	reCredit = regexp.MustCompile(`(?is)(?:Image|Video|Music)\s+Credit\s*:\s*</b>(.+?)(?:<b>|</center>|$)`)

	// Explanation: text after "Explanation:" label, more flexible with whitespace
	reExplanation = regexp.MustCompile(`(?is)Explanation\s*:\s*</b>(.+?)(?:<p>|<hr|$)`)

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
