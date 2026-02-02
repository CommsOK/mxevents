package mxevents

type EventFacts struct {
	Sender             SenderFacts
	Recipient          RecipientFacts
	SMTPResponse       string
	SMTPCode           string
	SMTPDeliveryStatus string
	// ReasonCode is the vendor-specific reason code/enum for the event outcome.
	// Example: HubSpot's "PREVIOUSLY_BOUNCED", SendGrid's "bounced"
	ReasonCode string
	// ReasonMessage is the human-readable explanation provided by the vendor.
	ReasonMessage string
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
