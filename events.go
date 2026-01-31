// Package mxevents provides a canonical taxonomy for email event classification.
// It serves as the source of truth for event types across the email delivery pipeline.
package mxevents

// Event represents a canonical email event type.
type Event string

const (
	// ====================================
	// Origin Platform (CRM / Vendor's Comms Layer)
	// ====================================

	// EventOriginSuccess indicates the vendor's communications platform
	// successfully accepted the message into its pipeline for processing.
	// Example: HubSpot or Customer.io acknowledging a campaign email internally.
	EventOriginSuccess Event = "origin-success"

	// EventOriginFailed indicates the vendor's communications platform rejected
	// the message before it reached the delivery vendor.
	// Example: invalid payload, missing recipient, or template rendering error.
	EventOriginFailed Event = "origin-failed"

	// EventOriginDropped indicates the vendor's communications platform intentionally
	// discarded the message (e.g., internal suppression, policy enforcement).
	EventOriginDropped Event = "origin-dropped"

	// ====================================
	// Gateway (ESP Layer)
	// ====================================

	// EventGatewayAccepted indicates the delivery gateway (ESP) acknowledged the message
	// and queued it for outbound delivery. This is the point where responsibility shifts
	// to the ESP's infrastructure.
	EventGatewayAccepted Event = "gateway-accepted"

	// EventGatewaySuccess indicates the delivery gateway successfully prepared or handed off
	// the message for SMTP delivery attempts.
	EventGatewaySuccess Event = "gateway-success"

	// EventGatewayFailed indicates the delivery gateway could not process or stage the send.
	// Example: misconfigured sender domain, authentication issue, or gateway internal error.
	EventGatewayFailed Event = "gateway-failed"

	// EventGatewayDropped indicates the delivery gateway accepted the message but chose
	// not to attempt outbound delivery (e.g., gateway-side suppression, blocklist match).
	EventGatewayDropped Event = "gateway-dropped"

	// ====================================
	// Mailbox (Target MX / Recipient Mail Server)
	// ====================================

	// EventMailboxAttempt indicates the ESP attempted an SMTP delivery to the recipient mailbox.
	// Useful for tracking retries vs. successful connections.
	EventMailboxAttempt Event = "mailbox-attempt"

	// EventMailboxSuccess indicates the recipient mailbox accepted the message with a 2xx response.
	// This means "delivered to server," but not inbox placement.
	EventMailboxSuccess Event = "mailbox-success"

	// EventMailboxFailed indicates a delivery attempt failed without a clear temp/permanent signal.
	// Example: gateway only reports "delivery failed."
	EventMailboxFailed Event = "mailbox-failed"

	// EventMailboxTempFail indicates a temporary SMTP failure (4xx). Gateway may retry.
	EventMailboxTempFail Event = "mailbox-tempfail"

	// EventMailboxSenderPermFail indicates a permanent SMTP failure (5xx) due to sender-side issues.
	// Examples: authentication failure, sender domain blocked, sender IP blacklisted,
	// DKIM/SPF/DMARC policy failures, sending limit exceeded.
	EventMailboxSenderPermFail Event = "mailbox-sender-permfail"

	// EventMailboxRecipientPermFail indicates a permanent SMTP failure (5xx) due to recipient-side issues.
	// Examples: mailbox does not exist, mailbox disabled, recipient domain does not exist,
	// recipient rejected message. This is the default when failure type is uncertain.
	EventMailboxRecipientPermFail Event = "mailbox-recipient-permfail"

	// EventMailboxQuarantined indicates the recipient mailbox accepted the message
	// but placed it in quarantine/spam (when gateway provides this signal).
	EventMailboxQuarantined Event = "mailbox-quarantined"

	// ====================================
	// Engagement (End User Actions)
	// ====================================

	// EventEngagementOpen indicates an open/render was recorded (tracking pixel fired).
	// Caveat: auto-opens and privacy features (e.g., Apple MPP).
	EventEngagementOpen Event = "engagement-open"

	// EventEngagementClick indicates the recipient clicked a tracked link in the email.
	// Stronger engagement signal than open.
	EventEngagementClick Event = "engagement-click"

	// EventEngagementEngaged is a composite signal defined by the customer.
	// Could be injected via API if vendor doesn't provide this signal.
	EventEngagementEngaged Event = "engagement-engaged"

	// ====================================
	// Status (List Management / Preferences)
	// ====================================
	//
	// These events represent recipient preference/state changes reported by the
	// platform (CRM/ESP). They are independent of the message lifecycle and are
	// primarily used for compliance and suppression logic.

	// EventStatusSubscribed indicates the recipient is globally subscribed
	// (eligible to receive messages subject to other policies and suppressions).
	EventStatusSubscribed Event = "status-subscribed"

	// EventStatusUnsubscribed indicates the recipient opted out globally.
	// Treat as a hard suppression signal for all future sends from this sender.
	EventStatusUnsubscribed Event = "status-unsubscribed"

	// EventStatusGroupSubscribed indicates the recipient opted in to a specific
	// topic/list/group (e.g., "Product Updates"). Not necessarily global.
	EventStatusGroupSubscribed Event = "status-group-subscribed"

	// EventStatusGroupUnsubscribed indicates the recipient opted out of a specific
	// topic/list/group but may remain subscribed to others.
	EventStatusGroupUnsubscribed Event = "status-group-unsubscribed"

	// EventStatusSpamReported indicates the recipient reported a message as spam
	// (e.g., via ISP feedback loop). Treat as an immediate global suppression
	// and record ISP/FBL details if available.
	EventStatusSpamReported Event = "status-spam-reported"

	// EventStatusSpamCleared indicates the recipient's spam complaint status has been
	// manually cleared (e.g., by operator review, recipient request, or automated
	// rehabilitation). This allows the recipient to receive emails again.
	EventStatusSpamCleared Event = "status-spam-cleared"
)
