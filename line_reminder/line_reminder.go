package line_reminder

type LineReminder struct {
	client *lineClient
}

func NewLineReminder() *LineReminder{
	return &LineReminder{
		client: NewLineClient(),
	}
}
