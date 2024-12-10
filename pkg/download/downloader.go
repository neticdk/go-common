package download

type Downloader interface {
	Download(url, dest string) (int64, error)
}
