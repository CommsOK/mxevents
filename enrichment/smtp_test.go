package enrichment

import (
	"context"
	"testing"

	"github.com/commsok/mxevents"
)

func TestSMTPEnricher_Enrich(t *testing.T) {
	tests := []struct {
		name       string
		facts      *mxevents.EventFacts
		wantCode   string
		wantStatus string
	}{
		{
			name: "extracts code and status from bounce response",
			facts: &mxevents.EventFacts{
				SMTPResponse: "550 5.1.1 The email account that you tried to reach does not exist",
			},
			wantCode:   "550",
			wantStatus: "5.1.1",
		},
		{
			name: "extracts mailbox full response",
			facts: &mxevents.EventFacts{
				SMTPResponse: "452 4.2.2 Mailbox full",
			},
			wantCode:   "452",
			wantStatus: "4.2.2",
		},
		{
			name: "does not overwrite existing SMTPCode",
			facts: &mxevents.EventFacts{
				SMTPCode:     "551",
				SMTPResponse: "550 5.1.1 User unknown",
			},
			wantCode:   "551",
			wantStatus: "5.1.1",
		},
		{
			name: "does not overwrite existing SMTPDeliveryStatus",
			facts: &mxevents.EventFacts{
				SMTPDeliveryStatus: "5.1.2",
				SMTPResponse:       "550 5.1.1 User unknown",
			},
			wantCode:   "550",
			wantStatus: "5.1.2",
		},
		{
			name: "handles temporary failure",
			facts: &mxevents.EventFacts{
				SMTPResponse: "421 4.7.0 Try again later",
			},
			wantCode:   "421",
			wantStatus: "4.7.0",
		},
		{
			name: "empty response does nothing",
			facts: &mxevents.EventFacts{
				SMTPResponse: "",
			},
			wantCode:   "",
			wantStatus: "",
		},
		{
			name: "extracts from verbose response",
			facts: &mxevents.EventFacts{
				SMTPResponse: "host mx.example.com[1.2.3.4] said: 550 5.7.1 Service unavailable; client blocked",
			},
			wantCode:   "550",
			wantStatus: "5.7.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enricher := &SMTPEnricher{}
			ctx := context.Background()

			err := enricher.Enrich(&ctx, tt.facts)
			if err != nil {
				t.Fatalf("Enrich() error = %v", err)
			}

			if tt.facts.SMTPCode != tt.wantCode {
				t.Errorf("SMTPCode = %q, want %q", tt.facts.SMTPCode, tt.wantCode)
			}

			if tt.facts.SMTPDeliveryStatus != tt.wantStatus {
				t.Errorf("SMTPDeliveryStatus = %q, want %q", tt.facts.SMTPDeliveryStatus, tt.wantStatus)
			}
		})
	}
}
