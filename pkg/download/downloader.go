package download

type Downloader interface {
	// Download downloads the file from the given URL and saves it to the given
	// destination.
	Download(url, dest string) (int64, error)
}
