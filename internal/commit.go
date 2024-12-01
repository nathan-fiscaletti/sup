package internal

import (
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	Hash   string
	Author string
	Date   time.Time
}

func (c *Commit) Scan(value string) error {
	parts := strings.Split(value, "\t")
	if len(parts) != 3 {
		return fmt.Errorf("invalid commit format: %s", value)
	}
	c.Hash = parts[0]
	c.Author = parts[1]
	var err error
	c.Date, err = time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return err
	}
	return nil
}

func (c *Commit) Subject() (string, error) {
	return c.Format("%s")
}

func (c *Commit) Body() (string, error) {
	return c.Format("%b")
}

func (c *Commit) RawMessage() (string, error) {
	return c.Format("%B")
}

func (c *Commit) ShortHash() (string, error) {
	return c.Format("%h")
}

func (c *Commit) Format(f string) (string, error) {
	res, err := Git("log", "-1", "--pretty=format:"+f, c.Hash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}
