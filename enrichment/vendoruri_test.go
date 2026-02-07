package enrichment

import (
	"context"
	"testing"

	"github.com/commsok/mxevents"
)

func TestCommonMailboxVendorURIEnricher_Enrich(t *testing.T) {
	tests := []struct {
		name           string
		facts          *mxevents.EventFacts
		expectedVendor mxevents.MailboxVendorURI
	}{
		// Code path: MailboxVendorURI already set - should skip
		{
			name: "does not overwrite existing MailboxVendorURI",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					MailboxVendorURI: mxevents.MailboxVendorURIOutlook,
					RecipientDomain:  "gmail.com",
				},
			},
			expectedVendor: mxevents.MailboxVendorURIOutlook,
		},
		// Code path: RecipientDomain is empty - should skip
		{
			name: "skips enrichment when RecipientDomain is empty",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					RecipientDomain: "",
				},
			},
			expectedVendor: "",
		},
		// Code path: Domain found in map - should set MailboxVendorURI
		{
			name: "sets MailboxVendorURI for known domain",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					RecipientDomain: "gmail.com",
				},
			},
			expectedVendor: mxevents.MailboxVendorURIGmail,
		},
		// Code path: Domain not found in map - should not set MailboxVendorURI
		{
			name: "does not set MailboxVendorURI for unknown domain",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					RecipientDomain: "example.com",
				},
			},
			expectedVendor: "",
		},
		// Code path: Case-insensitive matching
		{
			name: "domain matching is case-insensitive",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					RecipientDomain: "Gmail.com",
				},
			},
			expectedVendor: mxevents.MailboxVendorURIGmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enricher := &CommonMailboxVendorURIEnricher{}
			ctx := context.Background()

			err := enricher.Enrich(&ctx, tt.facts)
			if err != nil {
				t.Fatalf("Enrich() error = %v", err)
			}

			if tt.facts.Recipient.MailboxVendorURI != tt.expectedVendor {
				t.Errorf("MailboxVendorURI = %q, want %q", tt.facts.Recipient.MailboxVendorURI, tt.expectedVendor)
			}
		})
	}
}
