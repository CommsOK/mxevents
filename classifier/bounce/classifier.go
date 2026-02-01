package bounce

import (
	"context"
	"strings"

	"github.com/commsok/mxevents"
	"libsisimai.org/sisimai/v5/reason"
	"libsisimai.org/sisimai/v5/siba"
)

// Classifier classifies bounce events using sisimai's reason detection.
type Classifier struct {
}

// Classify analyzes SMTP response data and classifies the bounce event.
func (c *Classifier) Classify(ctx *context.Context, facts *mxevents.EventFacts, taxonomyVersion int) (*mxevents.ClassificationResult, error) {
	if facts.SMTPResponse == "" && facts.SMTPCode == "" && facts.SMTPDeliveryStatus == "" {
		return nil, nil
	}

	// Build a sisimai Fact struct from our EventFacts
	fact := buildSisimaiFactFromEventFacts(facts)

	// Use sisimai's reason.Find to detect the bounce reason
	sisimaiReason := reason.Find(fact)
	if sisimaiReason == "" {
		sisimaiReason = "undefined"
	}

	// Map sisimai reason to our canonical reason
	bounceReason := mapSisimaiReason(sisimaiReason)

	// Determine if this is a hard bounce based on SMTP code
	isHardBounce := isHardBounceFromCode(facts.SMTPCode, facts.SMTPDeliveryStatus)

	// Classify into sender, recipient, or generic bounce
	eventType := classifyBounceType(sisimaiReason, isHardBounce)

	// Determine confidence based on available data
	confidence := calculateConfidence(facts, sisimaiReason)

	return &mxevents.ClassificationResult{
		TaxonomyVersion: taxonomyVersion,
		EventType:       eventType,
		Reason:          bounceReason,
		Confidence:      confidence,
		Facts:           facts,
	}, nil
}

// buildSisimaiFactFromEventFacts constructs a sisimai Fact struct from EventFacts.
func buildSisimaiFactFromEventFacts(facts *mxevents.EventFacts) *siba.Fact {
	return &siba.Fact{
		DiagnosticCode: facts.SMTPResponse,
		DiagnosticType: "SMTP",
		ReplyCode:      facts.SMTPCode,
		DeliveryStatus: facts.SMTPDeliveryStatus,
	}
}

// isHardBounceFromCode determines if a bounce is permanent based on SMTP codes.
func isHardBounceFromCode(smtpCode, deliveryStatus string) bool {
	// Check SMTP reply code first (5xx = permanent, 4xx = temporary)
	if len(smtpCode) > 0 {
		if smtpCode[0] == '5' {
			return true
		}
		if smtpCode[0] == '4' {
			return false
		}
	}

	// Check delivery status code (5.x.x = permanent, 4.x.x = temporary)
	if len(deliveryStatus) > 0 {
		if deliveryStatus[0] == '5' {
			return true
		}
		if deliveryStatus[0] == '4' {
			return false
		}
	}

	// Default to hard bounce if we can't determine
	return true
}

// recipientReasons are bounce reasons attributable to the recipient.
var recipientReasons = map[string]bool{
	"userunknown": true, // Recipient mailbox doesn't exist
	"hostunknown": true, // Recipient domain doesn't exist
	"hasmoved":    true, // Recipient has moved
	"mailboxfull": true, // Recipient mailbox is full
	"suspend":     true, // Recipient account suspended
	"vacation":    true, // Recipient auto-reply
}

// senderReasons are bounce reasons attributable to the sender.
var senderReasons = map[string]bool{
	"authfailure":     true, // SPF/DKIM/DMARC failure
	"badreputation":   true, // IP reputation issue
	"blocked":         true, // IP/hostname blocked
	"requireptr":      true, // Missing PTR record
	"norelaying":      true, // Relay denied
	"rejected":        true, // Rejected due to envelope from
	"spamdetected":    true, // Spam filter rejection
	"policyviolation": true, // Policy violation
	"virusdetected":   true, // Virus detected
	"contenterror":    true, // Content error
	"notcompliantrfc": true, // RFC non-compliance
	"mesgtoobig":      true, // Message too big
	"exceedlimit":     true, // Size limit exceeded
	"syntaxerror":     true, // Syntax error in message
	"securityerror":   true, // Security issue
}

// classifyBounceType determines the event type based on the bounce reason and permanence.
func classifyBounceType(sisimaiReason string, isHardBounce bool) mxevents.EventType {
	reasonLower := strings.ToLower(sisimaiReason)

	// Temporary failures are always mailbox-tempfail
	if !isHardBounce {
		return mxevents.EventMailboxTempFail
	}

	// Check if recipient-related
	if recipientReasons[reasonLower] {
		return mxevents.EventMailboxRecipientPermFail
	}

	// Check if sender-related
	if senderReasons[reasonLower] {
		return mxevents.EventMailboxSenderPermFail
	}

	// Generic permanent failure for everything else
	return mxevents.EventMailboxPermFail
}

// mapSisimaiReason maps a sisimai reason to our canonical Reason type.
func mapSisimaiReason(sisimaiReason string) mxevents.Reason {
	// Sisimai reasons are lowercase, our canonical reasons are PascalCase
	reasonMap := map[string]mxevents.Reason{
		"userunknown":     mxevents.BounceReasonUserUnknown,
		"hostunknown":     mxevents.BounceReasonHostUnknown,
		"hasmoved":        mxevents.BounceReasonHasMoved,
		"mailboxfull":     mxevents.BounceReasonMailboxFull,
		"vacation":        mxevents.BounceReasonVacation,
		"spamdetected":    mxevents.BounceReasonSpamDetected,
		"badreputation":   mxevents.BounceReasonBadReputation,
		"blocked":         mxevents.BounceReasonBlocked,
		"policyviolation": mxevents.BounceReasonPolicyViolation,
		"authfailure":     mxevents.BounceReasonAuthFailure,
		"requireptr":      mxevents.BounceReasonRequirePTR,
		"failedstarttls":  mxevents.BounceReasonFailedSTARTTLS,
		"mesgtoobig":      mxevents.BounceReasonEmailTooLarge,
		"exceedlimit":     mxevents.BounceReasonEmailTooLarge,
		"virusdetected":   mxevents.BounceReasonVirusDetected,
		"contenterror":    mxevents.BounceReasonContentError,
		"notcompliantrfc": mxevents.BounceReasonNotCompliantRFC,
		"syntaxerror":     mxevents.BounceReasonSyntaxError,
		"norelaying":      mxevents.BounceReasonNoRelaying,
		"speeding":        mxevents.BounceReasonRateLimited,
		"toomanyconn":     mxevents.BounceReasonRateLimited,
		"systemfull":      mxevents.BounceReasonSystemFull,
		"systemerror":     mxevents.BounceReasonSystemError,
		"networkerror":    mxevents.BounceReasonNetworkError,
		"expired":         mxevents.BounceReasonExpired,
		"filtered":        mxevents.BounceReasonFiltered,
		"rejected":        mxevents.BounceReasonRejected,
		"notaccept":       mxevents.BounceReasonNotAccept,
		"mailererror":     mxevents.BounceReasonMailerError,
		"securityerror":   mxevents.BounceReasonSecurityError,
		"suspend":         mxevents.BounceReasonSuspend,
		"feedback":        mxevents.BounceReasonFeedback,
		"suppressed":      mxevents.BounceReasonSuppressed,
		"undefined":       mxevents.BounceReasonUndefined,
		"onhold":          mxevents.BounceReasonUndefined,
	}

	if reason, ok := reasonMap[strings.ToLower(sisimaiReason)]; ok {
		return reason
	}

	return mxevents.BounceReasonUndefined
}

// calculateConfidence determines the confidence level of the classification.
func calculateConfidence(facts *mxevents.EventFacts, sisimaiReason string) float32 {
	var confidence float32 = 0.5

	// Higher confidence if we have SMTP code
	if facts.SMTPCode != "" {
		confidence += 0.15
	}

	// Higher confidence if we have delivery status
	if facts.SMTPDeliveryStatus != "" {
		confidence += 0.15
	}

	// Higher confidence if we have diagnostic message
	if facts.SMTPResponse != "" {
		confidence += 0.1
	}

	// Lower confidence for undefined reasons
	if strings.ToLower(sisimaiReason) == "undefined" || strings.ToLower(sisimaiReason) == "onhold" {
		confidence -= 0.2
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}
