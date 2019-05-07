package scheduler

import (
	"fmt"
	"time"

	"github.com/kutsuzawa/line-reminder/factory"
	"github.com/kutsuzawa/line-reminder/util"
	"go.uber.org/zap"
)

// Checker checks status periodically.
type Checker struct {
	Message  string
	GroupID  string
	Duration time.Duration

	Logger         *zap.Logger
	ServiceFactory factory.ServiceFactory
}

// Schedule schedules checking targets's status periodically.
func (c *Checker) Schedule(targets []string) error {
	c.Logger.Info("check scheduler is started")
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
		c.Logger.Info("start to check",
			zap.String("id", t),
		)
		status, err := util.GetStatus(t)
		if err != nil {
			c.Logger.Error("failed to get status",
				zap.String("id", t),
				zap.Error(err),
			)
			return err
		}
		if !status {
			lineService, err := c.ServiceFactory.LineService()
			if err != nil {
				c.Logger.Error("failed to create line service",
					zap.String("id", t),
					zap.Error(err),
				)
				return err
			}
			name, err := lineService.GetNameByID(t)
			if err != nil {
				c.Logger.Error("failed to get name by id",
					zap.String("id", t),
					zap.Error(err),
				)
				return err
			}

			// e.g To cappyzawa
			// Good Morning
			msg := fmt.Sprintf("To %s\n%s", name, c.Message)
			if err := lineService.Send(c.GroupID, msg); err != nil {
				c.Logger.Error("failed to send message for check",
					zap.String("id", t),
					zap.Error(err),
				)
				return err
			}
			c.Logger.Info("cheking is done",
				zap.String("id", t),
			)
		}
	}
	return nil
}
