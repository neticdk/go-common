package secrets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileProviderTestSuite struct {
	suite.Suite
	testDir string
}

func (suite *FileProviderTestSuite) SetupSuite() {
	suite.testDir = suite.T().TempDir()
}

func (suite *FileProviderTestSuite) TestGetSecret_Success() {
	filePath := filepath.Join(suite.testDir, "test_secret.txt")
	provider := NewFileProvider(Location(filePath))
	secretContent := []byte("test_secret_value")

	err := os.WriteFile(filePath, secretContent, 0o600)
	suite.Require().NoError(err)

	secret, err := provider.RetrieveSecret(suite.T().Context())

	suite.Require().NoError(err)
	assert.Equal(suite.T(), secretContent, secret.Value)
	assert.Equal(suite.T(), ProviderFile, secret.Provider)
	assert.Equal(suite.T(), Location(filePath), secret.Location)
}

func (suite *FileProviderTestSuite) TestGetSecret_FileNotFound() {
	filePath := filepath.Join(suite.testDir, "non_existent_file.txt")
	provider := NewFileProvider(Location(filePath))

	secret, err := provider.RetrieveSecret(suite.T().Context())

	suite.Require().Error(err)
	assert.Nil(suite.T(), secret)
	assert.Contains(suite.T(), err.Error(), "file \"")
	assert.Contains(suite.T(), err.Error(), "\" does not exist")
}

func (suite *FileProviderTestSuite) TestGetSecret_FileNotRegular() {
	filePath := filepath.Join(suite.testDir, "test_dir")
	provider := NewFileProvider(Location(filePath))

	err := os.Mkdir(filePath, 0o700)
	suite.Require().NoError(err)

	secret, err := provider.RetrieveSecret(suite.T().Context())

	suite.Require().Error(err)
	assert.Nil(suite.T(), secret)
	assert.Contains(suite.T(), err.Error(), "file \"")
	assert.Contains(suite.T(), err.Error(), "\" is not a regular file")
}

func (suite *FileProviderTestSuite) TestGetSecret_InvalidPath() {
	provider := NewFileProvider(Location("invalid/relative/path"))

	secret, err := provider.RetrieveSecret(suite.T().Context())

	assert.Nil(suite.T(), secret)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "invalid file path: invalid/relative/path")
}

func (suite *FileProviderTestSuite) TestGetSecret_ErrorCheckingFileStats() {
	filePath := filepath.Join(suite.testDir, "secret.txt")
	provider := NewFileProvider(Location(filePath))

	// Create a directory with restricted permissions.
	err := os.MkdirAll(filepath.Dir(filePath), 0o700)
	suite.Require().NoError(err, "failed to create test directory")
	defer func() {
		// Restore permissions to cleanup
		err := os.Chmod(filepath.Dir(filePath), 0o700)
		suite.Require().NoError(err, "failed to cleanup test directory")
	}()

	secretContent := []byte("test_secret_value")

	err = os.WriteFile(filePath, secretContent, 0o600)
	suite.Require().NoError(err)

	err = os.Chmod(filepath.Dir(filePath), 0o000) // remove permissions to trigger an error from stat
	suite.Require().NoError(err, "failed to set no permissions on test directory")

	_, err = provider.RetrieveSecret(suite.T().Context())

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "checking if file")
}

func (suite *FileProviderTestSuite) TestGetSecret_InsecurePermissions() {
	filePath := filepath.Join(suite.testDir, "insecure_secret.txt")
	provider := NewFileProvider(Location(filePath))
	secretContent := []byte("test_secret_value")

	err := os.WriteFile(filePath, secretContent, 0o644) // World-readable permissions
	suite.Require().NoError(err)

	secret, err := provider.RetrieveSecret(suite.T().Context())

	suite.Require().Error(err)
	assert.Nil(suite.T(), secret)
	assert.Contains(suite.T(), err.Error(), "file \"")
	assert.Contains(suite.T(), err.Error(), "\" has insecure permissions (world-readable)")
}

func TestFileProviderTestSuite(t *testing.T) {
	suite.Run(t, new(FileProviderTestSuite))
}
