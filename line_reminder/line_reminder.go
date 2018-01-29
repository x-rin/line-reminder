package line_reminder

type lineReminder struct {
	client *lineClient
}

func NewLineReminder() *lineReminder{
	return &lineReminder{
		client: NewLineClient(),
	}
}
