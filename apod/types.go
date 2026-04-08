package apod

// MediaType describes what kind of media today's APOD contains.
type MediaType string

const (
	MediaImage   MediaType = "image"
	MediaVideo   MediaType = "video"   // NASA-hosted mp4 or other iframe embed
	MediaYouTube MediaType = "youtube" // YouTube embed
)

// APOD holds all scraped fields for a single day's entry.
// Fields that don't apply to the media type are omitted from JSON output.
type APOD struct {
	Type        MediaType `json:"media_type"`
	Date        string    `json:"date"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Explanation string    `json:"explanation"`
	Source      string    `json:"source"`
}
