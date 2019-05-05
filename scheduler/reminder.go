package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/kutsuzawa/line-reminder/factory"
	"github.com/kutsuzawa/line-reminder/util"
)

type Reminder struct {
	Message string
	GroupID string
	Hours   []string

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

func (r *Reminder) Schedule(targets []string) error {
	// initialize
	nowJST := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	remain, nextHour, err := r.calculateRemainTime(nowJST.Format("15:04"))
	log.Printf("now: %v, remain: %v, next remind: %s\n", nowJST, remain, nextHour)
	if err != nil {
		return err
	}

	for {
		select {
		case <-time.After(remain):
			for _, id := range targets {
				if err := r.remind(id); err != nil {
					return err
				}
			}
			remain, nextHour, err = r.calculateRemainTime(nextHour)
			log.Printf("now: %v, remain: %v, next remind: %s\n", nowJST, remain, nextHour)
			if err != nil {
				return err
			}
		}
	}

}

func (r *Reminder) remind(id string) error {
	lineService, err := r.ServiceFactory.LineService()
	if err != nil {
		return err
	}
	target, err := lineService.GetNameByID(id)
	if err != nil {
		return err
	}
	_, err = util.SetStatus(id, "false")
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("To %s\n%s", target, r.Message)
	if err := lineService.Send(r.GroupID, msg); err != nil {
		return err
	}
	return nil
}
