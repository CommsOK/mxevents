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
	// SourceVendor is the event source vendor/platform (ESP webhook, CRM, etc.).
	// Examples: mxevents.SourceVendorSendGrid, mxevents.SourceVendorHubSpot.
	SourceVendor SourceVendor
	EventName    string
}

type RecipientFacts struct {
	// MailboxVendorURI is a mailbox provider behavior bucket identifier.
	// Values look like domains but are NOT the recipient's domain.
	MailboxVendorURI MailboxVendorURI
	RecipientDomain  string
}

type ClassificationResult struct {
	TaxonomyVersion int
	EventType       EventType
	Reason          Reason
	Confidence      float32
	Facts           *EventFacts
}
