package artifact

// PullResult is a struct that represents the result of a pull operation
type PullResult struct {
	// Dir is the directory where the artifact was pulled to
	Dir string

	// If the puller discovers a version, it will be stored here
	Version string
}
