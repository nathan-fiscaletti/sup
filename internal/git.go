package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Git(args ...string) (string, error) {
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput() // Use CombinedOutput to capture both stdout and stderr
	if err != nil {
		return string(out), err
	}
	// normalize output
	normalizedOutput := strings.ReplaceAll(string(out), "\r\n", "\n")
	return normalizedOutput, nil
}

func IsGitRepository() (bool, error) {
	// Check if .git directory exists
	_, err := os.Stat(".git")
	if err == nil {
		// .git directory exists
		return true, nil
	}
	if os.IsNotExist(err) {
		// .git directory does not exist
		return false, nil
	}
	return false, err
}

type RepositoryDetails struct {
	// URL is the URL of the repository
	URL string
	// UrlMd5Sum is the MD5 sum of the repositories URL
	UrlMd5Sum string
}

func (r *RepositoryDetails) GetChangeCache() (*Cache[Updates], error) {
	userDir, err := GetUserDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get user directory: %v", err)
	}

	cachePath := filepath.Join(userDir, ".sup", r.UrlMd5Sum+".json")

	return &Cache[Updates]{
		Location: cachePath,
	}, nil
}

// GetRepositoryDetails returns the repository ID for the current repository
func GetRepositoryDetails() (*RepositoryDetails, error) {
	repoId, err := Git("remote", "get-url", "origin")
	if err != nil {
		return nil, err
	}

	repoId = strings.TrimSpace(repoId)
	hasher := md5.New()
	hasher.Write([]byte(repoId))
	return &RepositoryDetails{
		URL:       repoId,
		UrlMd5Sum: hex.EncodeToString(hasher.Sum(nil)),
	}, nil
}
