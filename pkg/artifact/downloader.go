package artifact

// Downloader is an interface for downloading files
type Downloader interface {
	Download(url, dest string) (int64, error)
}
