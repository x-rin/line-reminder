package scheduler

import (
	"fmt"
	"time"

	"github.com/kutsuzawa/line-reminder/service"
	"github.com/kutsuzawa/line-reminder/util"
)

// Checker checks status periodically.
type Checker struct {
	Message  string
	GroupID  string
	Line     service.LineService
	Duration time.Duration
}

// Schedule schedulee targets's status periodically
func (c *Checker) Schedule(targets []string) error {
	ticker := time.NewTicker(c.Duration).C

	for {
		select {
		case <-ticker:
			if err := c.Check(targets); err != nil {
				return err
			}
		}
	}
}

// Check checks targets's status
// If status is false, this post message to pager.
func (c *Checker) Check(targets []string) error {
	for _, t := range targets {
		status, err := util.GetStatus(t)
		if err != nil {
			return err
		}
		if !status {
			name, err := c.Line.GetNameByID(t)
			if err != nil {
				return err
			}

			// e.g To cappyzawa
			// Good Morning
			msg := fmt.Sprintf("To %s\n%s", name, c.Message)
			if err := c.Line.Send(c.GroupID, msg); err != nil {
				return err
			}
		}
	}
	return nil
}
