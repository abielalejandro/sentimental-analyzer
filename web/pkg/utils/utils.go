package utils

import (
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func newCloudEvent(
	id string,
	eType string,
	source string,
	data interface{}) cloudevents.Event {

	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetSource(source)
	event.SetType(eType)

	return event
}

func NewTextCloudEvent(
	id string,
	eType string,
	source string,
	data interface{}) cloudevents.Event {

	event := newCloudEvent(id, eType, source, data)
	event.SetDataContentType("text/plain")
	event.SetData(cloudevents.TextPlain, data)

	return event
}

func NewJsonCloudEvent(
	id string,
	eType string,
	source string,
	data interface{}) cloudevents.Event {

	event := newCloudEvent(id, eType, source, data)
	event.SetDataContentType("application/json")
	event.SetData(cloudevents.ApplicationJSON, data)

	return event
}
