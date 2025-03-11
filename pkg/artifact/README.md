# artifact

This package provides a simple way to manage downloadable artifacts in a
project. It supports downloading and extracting artifacts of different types.

Supported types:

- http(s) archive file (zip, tar, tar.gz)
- git repository
- Github release
- helm chart
- json-bundler

First create an Artifact instance describing the artifact to download.

Then use the puller.New() function to create a new puller instance.

Example:

```go
a := &artifact.Artifact{
	Name:       "my-artifact",
	Repository: "https://my-repos.com/my-repo",
	Branch:     "main",
	Tag:        "v1.0.0",
	CommitHash: "e7d1f4c...",
	Version:    "1.0.0",
	URL:        "https://my-repos.com/my-repo/archive/v1.0.0.zip",
	AssetName:  "my-artifact-1.0.0.zip",
	SubDir:     "my-subdir",
}


p := puller.NewPuller(puller.WithLogger(logger))
res, err := p.Pull(ctx, puller.MethodHelmChart, a, puller.WithHelmOptions(...))
```

The concrete puller method determines the options.
