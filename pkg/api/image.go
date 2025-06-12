package api

// Tag represents a single entry in the "tags" array
type Tag struct {
	Name        string `json:"name"`
	Phase       string `json:"phase"`
	PullSpec    string `json:"pullSpec"`
	DownloadURL string `json:"downloadURL"`
}

// Release represents the full JSON structure
type Release struct {
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}
