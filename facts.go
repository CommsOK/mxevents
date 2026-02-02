package mxevents

type EventFacts struct {
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
	VendorURI       string
	RecipientDomain string
}

type ClassificationResult struct {
	TaxonomyVersion int
	EventType       EventType
	Reason          Reason
	Confidence      float32
	Facts           *EventFacts
}
