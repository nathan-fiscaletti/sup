package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Cache[T any] struct {
	Location string
}

func (c *Cache[T]) Get() (*T, error) {
	var value T

	// Check if the cache file exists
	_, err := os.Stat(c.Location)
	if err != nil {
		return nil, err
	}

	// load the cache data
	data, err := os.ReadFile(c.Location)
	if err != nil {
		return nil, err
	}

	// unmarshal the data
	err = json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (c *Cache[T]) Set(value *T) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(c.Location), 0755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = os.WriteFile(c.Location, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetUserDirectory() (string, error) {
	// Attempt to get the current user
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}

	// Get the user's home directory
	homeDir := usr.HomeDir

	// Alternatively, check the environment variables for the home directory
	// for better compatibility in certain environments
	if homeDir == "" {
		homeDir = os.Getenv("HOME") // Unix-like systems
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE") // Windows
		}
	}

	return homeDir, nil
}
