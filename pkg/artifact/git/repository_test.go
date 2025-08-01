package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPullRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		art := &artifact.Artifact{
			Name:    "my-repo",
			URL:     "https://github.com/example/my-repo.git",
			BaseDir: tmpDir,
		}

		mockRepo := git.NewMockRepository(t)
		mockRepo.On("PlainCloneContext", mock.Anything, mock.Anything, mock.Anything).Return(&gogit.Repository{}, nil)

		opts := &RepositoryOptions{
			Repository: mockRepo,
		}

		result, err := PullRepository(context.Background(), art, opts)
		require.NoError(t, err)
		assert.NotNil(t, result)

		artifactDir := filepath.Join(tmpDir, "my-repo")
		assert.DirExists(t, artifactDir)
	})

	t.Run("success_with_single_file_artifact", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		art := &artifact.Artifact{
			Name:         "my-repo",
			URL:          "https://github.com/example/my-repo.git",
			BaseDir:      tmpDir,
			RelativePath: "single-file.yaml",
		}

		mockRepo := git.NewMockRepository(t)
		mockRepo.On("PlainCloneContext", mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				cloneDestPath := args.Get(1).(string)

				err = os.WriteFile(filepath.Join(cloneDestPath, "single-file.yaml"), []byte("file content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(cloneDestPath, "other-file.yaml"), []byte("file content"), 0644)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(cloneDestPath, "example"), 0755)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(cloneDestPath, ".git"), 0755)
				require.NoError(t, err)
			}).
			Return(&gogit.Repository{}, nil)

		opts := &RepositoryOptions{
			Repository: mockRepo,
		}

		result, err := PullRepository(context.Background(), art, opts)

		require.NoError(t, err)
		require.NotNil(t, result)

		finalDir := art.DestDir()
		finalFile := filepath.Join(finalDir, "single-file.yaml")
		noExistDir := filepath.Join(finalDir, "example")
		noExistFile := filepath.Join(finalDir, "other-file.yaml")

		assert.DirExists(t, finalDir, "The final destination directory should be created")
		assert.FileExists(t, finalFile, "The artifact file should exist in the destination directory")
		assert.NoDirExists(t, noExistDir, "The final destination should not contain example directory")
		assert.NoFileExists(t, noExistFile, "The final destination should not contain other-file.yaml")
	})

	t.Run("error_no_repository", func(t *testing.T) {
		art := &artifact.Artifact{
			Name: "my-repo",
			URL:  "https://github.com/example/my-repo.git",
		}
		opts := &RepositoryOptions{}

		_, err := PullRepository(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository is required")
	})

	t.Run("error_clone", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		mockRepo := git.NewMockRepository(t)
		mockRepo.On("PlainCloneContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("clone error"))

		art := &artifact.Artifact{
			Name:    "my-repo",
			URL:     "https://github.com/example/my-repo.git",
			BaseDir: tmpDir,
		}
		opts := &RepositoryOptions{
			Repository: mockRepo,
		}

		_, err = PullRepository(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cloning repository")
	})

	t.Run("error_checkout_commit", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		storer := memory.NewStorage()
		worktree := memfs.New()
		repo, err := gogit.Init(storer, worktree)
		require.NoError(t, err)

		mockRepo := git.NewMockRepository(t)
		mockRepo.On("PlainCloneContext", mock.Anything, mock.Anything, mock.Anything).Return(repo, nil)

		art := &artifact.Artifact{
			Name:       "my-repo",
			URL:        "https://github.com/example/my-repo.git",
			BaseDir:    tmpDir,
			CommitHash: "invalid-hash",
		}
		opts := &RepositoryOptions{
			Repository: mockRepo,
		}

		_, err = PullRepository(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "checking out commit")
	})
}
