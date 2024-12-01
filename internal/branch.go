package internal

import (
	"fmt"
	"strings"
	"time"
)

type Branch struct {
	Name string
}

func (b *Branch) CommitsSince(d time.Time) ([]Commit, error) {
	output, err := Git(
		"log",
		`--pretty=format:%H%x09%an%x09%ad`,
		"--date=iso-strict",
		fmt.Sprintf("--since=%s", d.Format(time.RFC3339)),
		fmt.Sprintf("%s", b.Name),
	)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	for _, line := range strings.Split(output, "\n") {
		var commit Commit
		err := commit.Scan(line)
		if err != nil {
			return nil, err
		}
		commits = append(commits, commit)
	}
	return commits, nil
}
