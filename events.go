// Package mxevents provides a canonical taxonomy for email event classification.
// It serves as the source of truth for event types across the email delivery pipeline.
package mxevents

// CurrentTaxonomyVersion is the current version of the taxonomy. This is a semver packaged as an int where each two
// digits represent a version number.
// Examples:
// - 100 becomes 0.1.0
// - 110 becomes 0.1.10
// - 10203 becomes 1.2.3
const CurrentTaxonomyVersion = 100 // 0.1.0

// EventType represents a canonical email event type.
type EventType string

// Reason represents a bounce/failure reason subtype providing additional
// context for why a delivery failed or was rejected.
type Reason string

const (
	// ====================================
	// Origin Platform (CRM / Vendor's Comms Layer)
	// ====================================

	// EventOriginSuccess indicates the vendor's communications platform
	// successfully accepted the message into its pipeline for processing.
	// Example: HubSpot or Customer.io acknowledging a campaign email internally.
	EventOriginSuccess EventType = "origin-success"

	// EventOriginFailed indicates the vendor's communications platform rejected
	// the message before it reached the delivery vendor.
	// Example: invalid payload, missing recipient, or template rendering error.
	EventOriginFailed EventType = "origin-failed"

	// EventOriginDropped indicates the vendor's communications platform intentionally
	// discarded the message (e.g., internal suppression, policy enforcement).
	EventOriginDropped EventType = "origin-dropped"

	// ====================================
	// Gateway (ESP Layer)
	// ====================================

	// EventGatewayAccepted indicates the delivery gateway (ESP) acknowledged the message
	// and queued it for outbound delivery. This is the point where responsibility shifts
	// to the ESP's infrastructure.
	EventGatewayAccepted EventType = "gateway-accepted"

	// EventGatewaySuccess indicates the delivery gateway successfully prepared or handed off
	// the message for SMTP delivery attempts.
	EventGatewaySuccess EventType = "gateway-success"

	// EventGatewayFailed indicates the delivery gateway could not process or stage the send.
	// Example: misconfigured sender domain, authentication issue, or gateway internal error.
	EventGatewayFailed EventType = "gateway-failed"

	// EventGatewayDropped indicates the delivery gateway accepted the message but chose
	// not to attempt outbound delivery (e.g., gateway-side suppression, blocklist match).
	EventGatewayDropped EventType = "gateway-dropped"

	// ====================================
	// Mailbox (Target MX / Recipient Mail Server)
	// ====================================

	// EventMailboxAttempt indicates the ESP attempted an SMTP delivery to the recipient mailbox.
	// Useful for tracking retries vs. successful connections.
	EventMailboxAttempt EventType = "mailbox-attempt"

	// EventMailboxSuccess indicates the recipient mailbox accepted the message with a 2xx response.
	// This means "delivered to server," but not inbox placement.
	EventMailboxSuccess EventType = "mailbox-success"

	// EventMailboxFailed indicates a delivery attempt failed before a reliable SMTP-level
	// outcome could be determined.
	// Examples include MX lookup failures, connection or TLS negotiation failures, or
	// gateways that report only a generic "delivery failed" status.
	//
	// Semantics:
	// - Retry expectation: unknown
	// - Attribution: unknown
	//
	// Notes:
	// - This event may occur before or during the SMTP session.
	// - It does not assert temporary or permanent failure.
	EventMailboxFailed EventType = "mailbox-failed"

	// EventMailboxTempFail indicates a temporary delivery failure (SMTP 4xx or equivalent
	// provider signal) where retry may succeed.
	//
	// Semantics:
	// - Retry expectation: may succeed
	// - Attribution: unknown
	//
	// Notes:
	// - Temporary failures may still repeat or escalate, but this event alone does not
	//   assert permanent failure.
	EventMailboxTempFail EventType = "mailbox-tempfail"

	// EventMailboxSenderPermFail indicates a permanent delivery failure attributable to
	// sender-side causes.
	// Examples: authentication failure, policy rejection, reputation-based blocking,
	// DKIM/SPF/DMARC failures, sending limits exceeded.
	//
	// Semantics:
	// - Retry expectation: not expected to succeed
	// - Attribution: sender-side
	//
	// Notes:
	// - This event does not imply recipient mailbox invalidity.
	EventMailboxSenderPermFail EventType = "mailbox-sender-permfail"

	// EventMailboxRecipientPermFail indicates a permanent delivery failure attributable to
	// recipient-side causes.
	// Examples: mailbox does not exist, mailbox disabled, recipient domain does not exist.
	//
	// Semantics:
	// - Retry expectation: not expected to succeed
	// - Attribution: recipient-side
	//
	// Notes:
	// - Use only when evidence supports recipient-side permanence.
	EventMailboxRecipientPermFail EventType = "mailbox-recipient-permfail"

	// EventMailboxPermFail indicates a permanent delivery failure (SMTP 5xx or equivalent
	// provider signal) where retry is not expected to succeed, but the failure cannot be
	// confidently attributed to sender-side or recipient-side causes.
	//
	// Semantics:
	// - Retry expectation: not expected to succeed
	// - Attribution: unknown
	//
	// Notes:
	// - This event intentionally avoids assigning blame without sufficient evidence.
	// - It may be reclassified if additional signals become available.
	EventMailboxPermFail EventType = "mailbox-permfail"

	// EventMailboxQuarantined indicates the recipient mailbox accepted the message
	// but placed it in quarantine/spam (when gateway provides this signal).
	EventMailboxQuarantined EventType = "mailbox-quarantined"

	// ====================================
	// Engagement (End User Actions)
	// ====================================

	// EventEngagementOpen indicates an open/render was recorded (tracking pixel fired).
	// Caveat: auto-opens and privacy features (e.g., Apple MPP).
	EventEngagementOpen EventType = "engagement-open"

	// EventEngagementClick indicates the recipient clicked a tracked link in the email.
	// Stronger engagement signal than open.
	EventEngagementClick EventType = "engagement-click"

	// EventEngagementEngaged is a composite signal defined by the customer.
	// Could be injected via API if vendor doesn't provide this signal.
	EventEngagementEngaged EventType = "engagement-engaged"

	// ====================================
	// Status (List Management / Preferences)
	// ====================================
	//
	// These events represent recipient preference/state changes reported by the
	// platform (CRM/ESP). They are independent of the message lifecycle and are
	// primarily used for compliance and suppression logic.

	// EventStatusSubscribed indicates the recipient is globally subscribed
	// (eligible to receive messages subject to other policies and suppressions).
	EventStatusSubscribed EventType = "status-subscribed"

	// EventStatusUnsubscribed indicates the recipient opted out globally.
	// Treat as a hard suppression signal for all future sends from this sender.
	EventStatusUnsubscribed EventType = "status-unsubscribed"

	// EventStatusGroupSubscribed indicates the recipient opted in to a specific
	// topic/list/group (e.g., "Product Updates"). Not necessarily global.
	EventStatusGroupSubscribed EventType = "status-group-subscribed"

	// EventStatusGroupUnsubscribed indicates the recipient opted out of a specific
	// topic/list/group but may remain subscribed to others.
	EventStatusGroupUnsubscribed EventType = "status-group-unsubscribed"

	// EventStatusSpamReported indicates the recipient reported a message as spam
	// (e.g., via ISP feedback loop). Treat as an immediate global suppression
	// and record ISP/FBL details if available.
	EventStatusSpamReported EventType = "status-spam-reported"

	// EventStatusSpamCleared indicates the recipient's spam complaint status has been
	// manually cleared (e.g., by operator review, recipient request, or automated
	// rehabilitation). This allows the recipient to receive emails again.
	EventStatusSpamCleared EventType = "status-spam-cleared"
)
