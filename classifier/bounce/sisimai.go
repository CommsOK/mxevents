package bounce

import (
	"context"
	"strings"

	"github.com/commsok/mxevents"
	"libsisimai.org/sisimai/v5/reason"
	"libsisimai.org/sisimai/v5/siba"
)

// SisimaiClassifier classifies bounce events using sisimai's reason detection.
type SisimaiClassifier struct {
}

// Classify analyzes SMTP response data and classifies the bounce event.
func (c *SisimaiClassifier) Classify(ctx *context.Context, facts *mxevents.EventFacts, taxonomyVersion int) (*mxevents.ClassificationResult, error) {
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

	// Set reason on fact so IsToxic() can evaluate it
	fact.Reason = sisimaiReason

	// Map sisimai reason to our canonical reason
	bounceReason := mapSisimaiReason(sisimaiReason)

	// Use sisimai's IsToxic() to determine if this is a permanent failure.
	// IsToxic() returns true when the address should be removed from mailing lists:
	// - Hard bounces: userunknown, hostunknown, hasmoved, notaccept, suspend, suppressed
	// - Soft bounces with 5xx: mailboxfull, filtered, norelaying
	// - Feedback loops: abuse, fraud, opt-out
	// This means 5xx infrastructure errors (networkerror, systemerror) are treated as
	// temporary since the address itself may still be valid.
	isHardBounce := fact.IsToxic()

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

// sisimaiReasonMap is the canonical mapping from Sisimai reason keys to mxevents Reasons.
//
// IMPORTANT: This mapping assumes the mxevents Reason taxonomy includes these additional Reasons:
// - smtp.ContentError
// - network.Error
// Optionally, you may also add smtp.RequirePTR; otherwise we map requireptr -> smtp.AuthFailure.
//
// If you do not want to add these Reasons, see the comments below for fallback mappings.
var sisimaiReasonMap = map[string]mxevents.Reason{
	// ====================================
	// SMTP (smtp.*)
	// ====================================

	// Recipient/mailbox existence/state
	"userunknown": mxevents.ReasonSMTPUserUnknown,
	"hasmoved":    mxevents.ReasonSMTPHasMoved,
	"mailboxfull": mxevents.ReasonSMTPMailboxFull,

	// Content / policy / filtering
	"spamdetected":    mxevents.ReasonSMTPSpamDetected,
	"policyviolation": mxevents.ReasonSMTPPolicyViolation,
	"blocked":         mxevents.ReasonSMTPBlocked,

	// Auth posture (SMTP evaluation)
	"authfailure": mxevents.ReasonSMTPAuthFailure,

	// Message size
	"mesgtoobig":  mxevents.ReasonSMTPMessageTooLarge,
	"exceedlimit": mxevents.ReasonSMTPMessageTooLarge,

	// Malware
	"virusdetected": mxevents.ReasonSMTPVirusDetected,

	// Protocol / content correctness
	"syntaxerror": mxevents.ReasonSMTPSyntaxError,

	// These are *not* just syntax; they are malformed MIME / RFC compliance failures.
	// Requires mxevents.ReasonSMTPContentError = "smtp.ContentError".
	"contenterror":    mxevents.ReasonSMTPContentError,
	"notcompliantrfc": mxevents.ReasonSMTPContentError,

	// Relay policy
	"norelaying": mxevents.ReasonSMTPNoRelaying,

	// Throttling / deferrals
	"speeding":    mxevents.ReasonSMTPRateLimited,
	"toomanyconn": mxevents.ReasonSMTPRateLimited,
	"ratelimited": mxevents.ReasonSMTPRateLimited,

	// Final failure after retries
	"expired": mxevents.ReasonSMTPExpired,

	// ====================================
	// Network / transport (network.*)
	// ====================================

	// DNS / MX resolution
	"hostunknown": mxevents.ReasonNetworkDnsFailure,

	// Broad bucket in Sisimai: may include timeouts, resets, routing, etc.
	// Requires mxevents.ReasonNetworkError = "network.Error".
	"networkerror": mxevents.ReasonNetworkError,

	// TLS handshake / STARTTLS failures
	"failedstarttls": mxevents.ReasonNetworkTlsFailure,

	// ====================================
	// Policy (policy.*)
	// ====================================

	"badreputation": mxevents.ReasonPolicyBadReputation,

	// These are often policy-like signals in diagnostics; they can be coarse.
	"filtered":  mxevents.ReasonPolicyBlocked,
	"rejected":  mxevents.ReasonPolicyBlocked,
	"notaccept": mxevents.ReasonPolicyRestricted,

	// "securityerror" and "suspend" can be account/security enforcement.
	// Keeping coarse is OK; revisit if you see strong patterns in real payloads.
	"securityerror": mxevents.ReasonPolicyBlocked,
	"suspend":       mxevents.ReasonPolicyBlocked,

	// requireptr is specifically reverse DNS/PTR posture. It is not a generic policy block.
	// If you add mxevents.ReasonSMTPRequirePTR = "smtp.RequirePTR", map to that instead.
	"requireptr": mxevents.ReasonSMTPAuthFailure,

	// ====================================
	// System / operational (system.*)
	// ====================================

	// Remote/system capacity. We treat Sisimai's "systemfull" as backpressure.
	"systemfull": mxevents.ReasonSystemBackpressure,

	// Generic system errors
	"systemerror": mxevents.ReasonSystemError,
	"mailererror": mxevents.ReasonSystemError,

	// ====================================
	// Complaint / feedback (complaint.*)
	// ====================================

	// Sisimai's "feedback" is complaint-like but can be broader than strict ISP FBL.
	// If you add mxevents.ReasonComplaintSignal = "complaint.Signal", map feedback to that.
	// Otherwise, use FeedbackLoop as the closest actionable complaint bucket.
	"feedback": mxevents.ReasonComplaintFeedbackLoop,

	// ====================================
	// Suppression (suppression.*)
	// ====================================

	// Pre-delivery suppression (not SMTP)
	"suppressed": mxevents.ReasonSuppressionSuppressed,

	// ====================================
	// Unknown / fallback (unknown.*)
	// ====================================

	"undefined": mxevents.ReasonUnknown,

	// "onhold" and "vacation" exist in Sisimai but are not reliably actionable
	// in a deliverability hygiene product without additional context.
	"onhold":   mxevents.ReasonUnknown,
	"vacation": mxevents.ReasonUnknown,
}

// mapSisimaiReason maps a sisimai reason to our canonical Reason type.
func mapSisimaiReason(sisimaiReason string) mxevents.Reason {
	if reason, ok := sisimaiReasonMap[strings.ToLower(sisimaiReason)]; ok {
		return reason
	}

	return mxevents.ReasonUnknown
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
