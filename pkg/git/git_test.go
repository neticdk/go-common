package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_Init(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-common-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	g := NewRepository(nil)
	cfg := config.NewConfig()
	cfg.Author.Name = "Test"
	cfg.Author.Email = "test@example.com"

	t.Run("Init", func(t *testing.T) {
		_, err = g.Init(tmpDir, "main")
		assert.NoError(t, err)
		assert.DirExists(t, filepath.Join(tmpDir, ".git"))
	})

	t.Run("CreateRemote", func(t *testing.T) {
		_, err = g.CreateRemote("origin", "https://example.com/repo.git")
		assert.NoError(t, err)
	})

	t.Run("Add", func(t *testing.T) {
		err := os.WriteFile(filepath.Join(tmpDir, "test1.txt"), []byte("test1"), 0o640)
		assert.NoError(t, err)
		err = os.WriteFile(filepath.Join(tmpDir, "test2.txt"), []byte("test2"), 0o640)
		assert.NoError(t, err)
		err = g.Add("test1.txt", "test2.txt")
		assert.NoError(t, err)
	})

	t.Run("Commit", func(t *testing.T) {
		g.Repo().SetConfig(cfg)
		hash, err := g.Commit("Initial commit", nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestGitRepository_InitAndCommit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-common-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	g := NewRepository(nil)
	cfg := config.NewConfig()
	cfg.Author.Name = "Test"
	cfg.Author.Email = "test@example.com"

	t.Run("InitAndCommit", func(t *testing.T) {
		err = os.WriteFile(filepath.Join(tmpDir, "test1.txt"), []byte("test1"), 0o640)
		assert.NoError(t, err)
		err = g.InitAndCommit(tmpDir, "https://example.com/repo.git", cfg, nil)
		assert.NoError(t, err)
		assert.DirExists(t, filepath.Join(tmpDir, ".git"))

		// Check if the remote is created
		remotes, err := g.Repo().Remotes()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(remotes))
		assert.Equal(t, "origin", remotes[0].Config().Name)

		// Check if the commit is made
		head, err := g.Repo().Head()
		assert.NoError(t, err)
		assert.NotNil(t, head)
	})
}
