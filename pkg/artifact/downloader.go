package artifact

// Downloader is an interface for downloading files
type Downloader interface {
	// Download downloads a file from the given URL to the given destination
	Download(url, dest string) (int64, error)
}
