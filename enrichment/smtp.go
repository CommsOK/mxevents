package enrichment

import (
	"context"

	"github.com/commsok/mxevents"
	"libsisimai.org/sisimai/v5/smtp/reply"
	"libsisimai.org/sisimai/v5/smtp/status"
)

// SMTPEnricher extracts SMTP facts from raw response data using Sisimai's parsing.
type SMTPEnricher struct {
}

// Enrich parses the SMTP response and fills in missing SMTPCode and SMTPDeliveryStatus fields.
func (e *SMTPEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
	if facts.SMTPResponse == "" {
		return nil
	}

	// Use existing fields as hints for extraction
	codeHint := facts.SMTPCode
	statusHint := facts.SMTPDeliveryStatus

	// Fill SMTPCode if missing
	if facts.SMTPCode == "" {
		if code := reply.Find(facts.SMTPResponse, statusHint); code != "" && reply.Test(code) {
			facts.SMTPCode = code
		}
	}

	// Fill SMTPDeliveryStatus if missing
	if facts.SMTPDeliveryStatus == "" {
		if ds := status.Find(facts.SMTPResponse, codeHint); ds != "" && status.Test(ds) {
			facts.SMTPDeliveryStatus = ds
		}
	}

	return nil
}
