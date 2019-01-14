package reminders

import "time"

type Reminder struct {
	ID       int
	Headers  map[string]interface{}
	Body     string
	Tags     []string
	DateTime time.Time
}
