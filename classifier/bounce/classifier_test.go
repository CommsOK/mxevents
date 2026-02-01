package bounce

import (
	"context"
	"testing"

	"github.com/commsok/mxevents"
)

func TestClassifier_Classify(t *testing.T) {
	classifier := &Classifier{}
	ctx := context.Background()

	tests := []struct {
		name           string
		facts          *mxevents.EventFacts
		expectedType   mxevents.EventType
		expectedReason mxevents.Reason
		expectNil      bool
		minConfidence  float32
	}{
		{
			name: "user unknown - recipient permfail",
			facts: &mxevents.EventFacts{
				SMTPCode:           "550",
				SMTPDeliveryStatus: "5.1.1",
				SMTPResponse:       "550 5.1.1 The email account that you tried to reach does not exist",
			},
			expectedType:   mxevents.EventMailboxRecipientPermFail,
			expectedReason: mxevents.BounceReasonUserUnknown,
			minConfidence:  0.8,
		},
		{
			name: "mailbox full - recipient permfail",
			facts: &mxevents.EventFacts{
				SMTPCode:           "552",
				SMTPDeliveryStatus: "5.2.2",
				SMTPResponse:       "552 5.2.2 Mailbox full",
			},
			expectedType:   mxevents.EventMailboxRecipientPermFail,
			expectedReason: mxevents.BounceReasonMailboxFull,
			minConfidence:  0.8,
		},
		{
			name: "mailbox full temporary - tempfail",
			facts: &mxevents.EventFacts{
				SMTPCode:           "452",
				SMTPDeliveryStatus: "4.2.2",
				SMTPResponse:       "452 4.2.2 Mailbox full, try again later",
			},
			expectedType:   mxevents.EventMailboxTempFail,
			expectedReason: mxevents.BounceReasonMailboxFull,
			minConfidence:  0.8,
		},
		{
			name: "authentication failure - tempfail (address may be valid)",
			facts: &mxevents.EventFacts{
				SMTPCode:           "550",
				SMTPDeliveryStatus: "5.7.26",
				SMTPResponse:       "550 5.7.26 This message does not pass authentication checks (SPF and DKIM)",
			},
			// IsToxic() returns false for auth failures - the address itself is valid,
			// it's a sender configuration issue that could be fixed
			expectedType:   mxevents.EventMailboxTempFail,
			expectedReason: mxevents.BounceReasonAuthFailure,
			minConfidence:  0.8,
		},
		{
			name: "security error (5.7.1) - tempfail (address may be valid)",
			facts: &mxevents.EventFacts{
				SMTPCode:           "550",
				SMTPDeliveryStatus: "5.7.1",
				SMTPResponse:       "550 5.7.1 Connection refused, your IP is listed in a blocklist",
			},
			// IsToxic() returns false for security errors - the address itself is valid
			expectedType:   mxevents.EventMailboxTempFail,
			expectedReason: mxevents.BounceReasonSecurityError,
			minConfidence:  0.8,
		},
		{
			name: "spam detected - tempfail (address may be valid)",
			facts: &mxevents.EventFacts{
				SMTPCode:           "550",
				SMTPDeliveryStatus: "5.7.1",
				SMTPResponse:       "550 5.7.1 Message rejected as spam",
			},
			// IsToxic() returns false for spam - the address itself is valid
			expectedType:   mxevents.EventMailboxTempFail,
			expectedReason: mxevents.BounceReasonSpamDetected,
			minConfidence:  0.8,
		},
		{
			name: "rate limited - tempfail",
			facts: &mxevents.EventFacts{
				SMTPCode:           "421",
				SMTPDeliveryStatus: "4.7.0",
				SMTPResponse:       "421 4.7.0 Too many connections, try again later",
			},
			expectedType:   mxevents.EventMailboxTempFail,
			expectedReason: mxevents.BounceReasonRateLimited,
			minConfidence:  0.8,
		},
		{
			name: "host unknown - recipient permfail",
			facts: &mxevents.EventFacts{
				SMTPCode:           "550",
				SMTPDeliveryStatus: "5.1.2",
				SMTPResponse:       "550 5.1.2 Host unknown",
			},
			expectedType:   mxevents.EventMailboxRecipientPermFail,
			expectedReason: mxevents.BounceReasonHostUnknown,
			minConfidence:  0.8,
		},
		{
			name: "empty facts - returns nil",
			facts: &mxevents.EventFacts{
				SMTPCode:           "",
				SMTPDeliveryStatus: "",
				SMTPResponse:       "",
			},
			expectNil: true,
		},
		{
			name: "only smtp code - tempfail (undefined reason not toxic)",
			facts: &mxevents.EventFacts{
				SMTPCode: "550",
			},
			// IsToxic() returns false for undefined reasons
			expectedType:  mxevents.EventMailboxTempFail,
			minConfidence: 0.4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := classifier.Classify(&ctx, tt.facts, mxevents.CurrentTaxonomyVersion)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectNil {
				if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.EventType != tt.expectedType {
				t.Errorf("expected event type %s, got %s", tt.expectedType, result.EventType)
			}

			if tt.expectedReason != "" && result.Reason != tt.expectedReason {
				t.Errorf("expected reason %s, got %s", tt.expectedReason, result.Reason)
			}

			if result.Confidence < tt.minConfidence {
				t.Errorf("expected confidence >= %f, got %f", tt.minConfidence, result.Confidence)
			}

			if result.TaxonomyVersion != mxevents.CurrentTaxonomyVersion {
				t.Errorf("expected taxonomy version %d, got %d", mxevents.CurrentTaxonomyVersion, result.TaxonomyVersion)
			}
		})
	}
}

func TestClassifyBounceType(t *testing.T) {
	tests := []struct {
		reason       string
		isHardBounce bool
		expected     mxevents.EventType
	}{
		// Recipient-related hard bounces
		{"userunknown", true, mxevents.EventMailboxRecipientPermFail},
		{"hostunknown", true, mxevents.EventMailboxRecipientPermFail},
		{"mailboxfull", true, mxevents.EventMailboxRecipientPermFail},
		{"suspend", true, mxevents.EventMailboxRecipientPermFail},

		// Sender-related hard bounces
		{"authfailure", true, mxevents.EventMailboxSenderPermFail},
		{"blocked", true, mxevents.EventMailboxSenderPermFail},
		{"spamdetected", true, mxevents.EventMailboxSenderPermFail},
		{"badreputation", true, mxevents.EventMailboxSenderPermFail},

		// Generic hard bounces
		{"networkerror", true, mxevents.EventMailboxPermFail},
		{"systemerror", true, mxevents.EventMailboxPermFail},
		{"undefined", true, mxevents.EventMailboxPermFail},

		// All soft bounces are tempfail regardless of reason
		{"userunknown", false, mxevents.EventMailboxTempFail},
		{"authfailure", false, mxevents.EventMailboxTempFail},
		{"networkerror", false, mxevents.EventMailboxTempFail},
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			result := classifyBounceType(tt.reason, tt.isHardBounce)
			if result != tt.expected {
				t.Errorf("classifyBounceType(%q, %v) = %s, want %s",
					tt.reason, tt.isHardBounce, result, tt.expected)
			}
		})
	}
}

func TestMapSisimaiReason(t *testing.T) {
	tests := []struct {
		sisimaiReason  string
		expectedReason mxevents.Reason
	}{
		{"userunknown", mxevents.BounceReasonUserUnknown},
		{"hostunknown", mxevents.BounceReasonHostUnknown},
		{"mailboxfull", mxevents.BounceReasonMailboxFull},
		{"authfailure", mxevents.BounceReasonAuthFailure},
		{"blocked", mxevents.BounceReasonBlocked},
		{"spamdetected", mxevents.BounceReasonSpamDetected},
		{"undefined", mxevents.BounceReasonUndefined},
		{"unknownreason", mxevents.BounceReasonUndefined}, // unknown maps to undefined
		{"USERUNKNOWN", mxevents.BounceReasonUserUnknown}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.sisimaiReason, func(t *testing.T) {
			result := mapSisimaiReason(tt.sisimaiReason)
			if result != tt.expectedReason {
				t.Errorf("mapSisimaiReason(%q) = %s, want %s",
					tt.sisimaiReason, result, tt.expectedReason)
			}
		})
	}
}
