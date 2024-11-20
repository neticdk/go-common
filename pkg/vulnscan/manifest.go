package vulnscan

import "io"

type manifestScanner struct {
	manifestPath string
	manifest     io.Reader
}
