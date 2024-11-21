package sbom

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImageNamesFromManifest(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		want     []string
		wantErr  bool
	}{
		{
			name: "simple pod manifest",
			manifest: `---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
  - name: sidecar
    image: sidecar:latest`,
			want: []string{"nginx:1.14.2", "sidecar:latest"},
		},
		{
			name: "multiple pod manifests",
			manifest: `---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
  - name: sidecar
    image: sidecar:latest
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx2
spec:
  containers:
  - name: nginx2
    image: nginx:1.14.3
  - name: sidecar
    image: sidecar:latest
`,
			want: []string{"nginx:1.14.2", "sidecar:latest", "nginx:1.14.3", "sidecar:latest"},
		},
		{
			name: "deployment manifest",
			manifest: `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
      - name: proxy
        image: proxy:1.0.0`,
			want: []string{"nginx:1.14.2", "proxy:1.0.0"},
		},
		{
			name: "statefulset manifest",
			manifest: `---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  template:
    spec:
      containers:
      - name: webapp
        image: webapp:v1
      - name: db
        image: mysql:5.7`,
			want: []string{"webapp:v1", "mysql:5.7"},
		},
		{
			name: "empty manifest",
			manifest: `---
apiVersion: v1
kind: Pod
metadata:
  name: empty
spec:
  containers: []`,
			want: []string{},
		},
		{
			name: "invalid manifest",
			manifest: `
invalid:
  yaml: [`,
			wantErr: true,
		},
		{
			name: "manifest without containers",
			manifest: `---
apiVersion: v1
kind: Service
metadata:
  name: my-service`,
			want: []string{},
		},
		{
			name: "manifest with initContainers and containers",
			manifest: `---
apiVersion: v1
kind: Pod
metadata:
  name: with-init
spec:
  initContainers:
  - name: init
    image: init:v1
  containers:
  - name: app
    image: app:v1`,
			want: []string{"app:v1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			reader := strings.NewReader(tt.manifest)

			got, err := extractImageNamesFromManifest(ctx, reader)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
