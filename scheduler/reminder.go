package scheduler

import (
	"fmt"
	"time"

	"github.com/kutsuzawa/line-reminder/factory"
	"github.com/kutsuzawa/line-reminder/util"
	"go.uber.org/zap"
)

type Reminder struct {
	Message string
	GroupID string
	Hours   []string

	Logger         *zap.Logger
	ServiceFactory factory.ServiceFactory
}

func (r *Reminder) calculateRemainTime(timeStr string) (time.Duration, string, error) {
	now, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, "", err
	}
	for _, remindHour := range r.Hours {
		remindTime, err := time.Parse("15:04", remindHour)
		if err != nil {
			return 0, "", err
		}
		if now.Before(remindTime) {
			remain := remindTime.Sub(now)
			return remain, remindHour, nil
		}
	}

	// when current time is later than r.hours, next remind will be occured in next day
	remindTime, err := time.Parse("15:04", r.Hours[0])
	if err != nil {
		return 0, "", err
	}
	return remindTime.Add(24 * time.Hour).Sub(now), r.Hours[0], nil
}

// Schedule schedules to remind based on specified time.
func (r *Reminder) Schedule(targets []string) error {
	r.Logger.Info("remind scheduler is stated")
	// initialize
	nowJST := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	remain, nextHour, err := r.calculateRemainTime(nowJST.Format("15:04"))
	if err != nil {
		r.Logger.Error("failed to calculate remain time",
			zap.Error(err),
		)
		return err
	}
	r.Logger.Info("before reminding",
		zap.Time("now", nowJST),
		zap.Float64("remain", remain.Hours()),
		zap.String("next remind", nextHour),
	)
	for {
		select {
		case <-time.After(remain):
			for _, id := range targets {
				if err := r.Remind(id); err != nil {
					return err
				}
			}
			remain, nextHour, err = r.calculateRemainTime(nextHour)
			if err != nil {
				r.Logger.Error("failed to calculate remain time",
					zap.Error(err),
				)
				return err
			}
			r.Logger.Info("before reminding",
				zap.Time("now", nowJST),
				zap.Float64("remain", remain.Hours()),
				zap.String("next remind", nextHour),
			)
		}
	}

}

func (r *Reminder) Remind(id string) error {
	r.Logger.Info("start to remind",
		zap.String("id", id),
	)
	lineService, err := r.ServiceFactory.LineService()
	if err != nil {
		r.Logger.Error("failed to create line service",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}
	target, err := lineService.GetNameByID(id)
	if err != nil {
		r.Logger.Error("failed to get name by id",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}
	_, err = util.SetStatus(id, "false")
	if err != nil {
		r.Logger.Error("failed to set status",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}
	msg := fmt.Sprintf("To %s\n%s", target, r.Message)
	if err := lineService.Send(r.GroupID, msg); err != nil {
		r.Logger.Error("failed to send reminder",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}
	r.Logger.Info("reminder is done",
		zap.String("id", id),
	)
	return nil
}
