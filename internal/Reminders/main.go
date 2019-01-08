package reminders

import "time"

type Reminder struct{
	ReminderId int
	Headers map[string]interface{}
	Body string
	Tags []string
	ReminderDateTime time.Time
}
