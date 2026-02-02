package enrichment

import (
	"context"
	"testing"

	"github.com/commsok/mxevents"
)

func TestCommonVendorURIEnricher_Enrich(t *testing.T) {
	tests := []struct {
		name           string
		facts          *mxevents.EventFacts
		expectedVendor string
	}{
		// Code path: VendorURI already set - should skip
		{
			name: "does not overwrite existing VendorURI",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					VendorURI:       mxevents.VendorURIOutlook,
					RecipientDomain: "gmail.com",
				},
			},
			expectedVendor: mxevents.VendorURIOutlook,
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
		// Code path: Domain found in map - should set VendorURI
		{
			name: "sets VendorURI for known domain",
			facts: &mxevents.EventFacts{
				Recipient: mxevents.RecipientFacts{
					RecipientDomain: "gmail.com",
				},
			},
			expectedVendor: mxevents.VendorURIGmail,
		},
		// Code path: Domain not found in map - should not set VendorURI
		{
			name: "does not set VendorURI for unknown domain",
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
			expectedVendor: mxevents.VendorURIGmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enricher := &CommonVendorURIEnricher{}
			ctx := context.Background()

			err := enricher.Enrich(&ctx, tt.facts)
			if err != nil {
				t.Fatalf("Enrich() error = %v", err)
			}

			if tt.facts.Recipient.VendorURI != tt.expectedVendor {
				t.Errorf("VendorURI = %q, want %q", tt.facts.Recipient.VendorURI, tt.expectedVendor)
			}
		})
	}
}
