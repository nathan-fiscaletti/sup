package internal

import (
	"encoding/json"
	"strings"
	"time"
)

type Update struct {
	BranchName string    `json:"branch"`
	Date       time.Time `json:"date"`
}

func (c *Update) Branch() *Branch {
	return &Branch{Name: c.BranchName}
}

type Updates struct {
	Updates []Update `json:"updates"`
}

// Get will return the difference between two Changes objects based on
// the date of the last commit. If a later date is found in the other Changes
// object, it will be added to the returned slice.
func (c *Updates) Compare(other *Updates) []Update {
	var difference []Update
	for _, change := range c.Updates {
		found := false
		for _, otherChange := range other.Updates {
			if change.BranchName == otherChange.BranchName {
				found = true
				if change.Date.After(otherChange.Date) {
					difference = append(difference, change)
				}
			}
		}
		if !found {
			difference = append(difference, change)
		}
	}
	return difference
}

func GetRemoteUpdates() (*Updates, error) {
	output, err := Git("for-each-ref", "--sort=-committerdate", `--format={"branch":"%(refname:short)","date":"%(committerdate:iso-strict)"}`, "refs/remotes/")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	changes := []Update{}
	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if len(stripped) < 1 {
			continue
		}
		var change Update
		err := json.Unmarshal([]byte(stripped), &change)
		if err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}

	return &Updates{changes}, nil
}
