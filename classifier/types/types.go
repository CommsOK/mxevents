package types

import "github.com/commsok/mxevents"

type Facts struct {
	Sender             SenderFacts
	Recipient          RecipientFacts
	SMTPResponse       string
	SMTPCode           string
	SMTPDeliveryStatus string
}

type SenderFacts struct {
	Vendor    string
	EventName string
}

type RecipientFacts struct {
	Vendor string
}

type ClassificationResult struct {
	EventType  mxevents.EventType
	Reason     mxevents.Reason
	Confidence float32
	Facts      *Facts
}
