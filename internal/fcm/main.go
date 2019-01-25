package fcm

import (
	"mynsb-api/internal/4U"
	"mynsb-api/internal/filesint"
	"github.com/NaySoftware/go-fcm"
	"mynsb-api/internal/events"
	"mynsb-api/internal/db"
)

var ServerKey string
const (
	FourUTopic = "/topics/4U"
	EventTopic = "/topics/event"
)

// Initialise our data
func init() {
	txt, _ := filesint.DataDump("sensitive", "/fcm/serverkey.txt")

	ServerKey = string(txt)
}











// TOPIC SENDERS

// alertNewFourU takes an issue and sends it to all the devices subscribed to the FourU topic
func alertNewFourU(issue fouru.Issue) {

	// Data to send through to the topic
	data := map[string]string {
		"issue_name": issue.Name,
		"issue_desc": issue.Desc,
		"issue_link": issue.Link,
	}

	// Set up the client
	client := fcm.NewFcmClient(ServerKey)

	// Send the data through to the topic
	client.NewFcmMsgTo(FourUTopic, data)
	client.Send()
}

// alertEventRunning takes an event, sets up a client and sends it to all devices subscribed to the Event topic
func alertEventRunning(event events.Event) {

	// data to send through to the topic
	data := map[string]string {
		"event_name": event.Name,
		"event_desc": event.ShortDesc,
		"event_time": event.Start.String(),
	}

	// Build a client
	client := fcm.NewFcmClient(ServerKey)

	// Send the data through to the topic
	client.NewFcmMsgTo(EventTopic, data)
	client.Send()
}












// CORE FUNCTIONS

// SendNew4U gets executed when there is a new 4U Issue
func SendNew4U() {

	// Connect to database
	db.Conn("student")
	defer db.DB.Close()

	res, _ := db.DB.Query("SELECT * FROM four_u WHERE article_id = (SELECT max(article_id) FROM four_u)")
	issue := fouru.Issue{}
	issue.ReplaceWith(res)

	// Send the new Four U thingy
	alertNewFourU(issue)
}