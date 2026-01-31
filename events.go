// Package mxevents provides a canonical taxonomy for email event classification.
// It serves as the source of truth for event types across the email delivery pipeline.
package mxevents

// Event represents a canonical email event type.
type Event string

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

// ====================================
// Bounce / Failure Reasons (Sisimai-derived)
// ====================================
//
// These Reason values are produced by the Sisimai library based on SMTP reply
// codes, enhanced status codes (when present), and diagnostic message parsing.
//
// Important constraints:
//
//   - Reasons describe *why* a delivery attempt failed, as inferred from the
//     remote server's response. They do NOT reliably indicate who is "at fault"
//     (sender vs recipient) on their own.
//   - Temporary vs permanent is determined primarily via Sisimai's hardbounce
//     signal (or SMTP 4xx vs 5xx when available), not by the reason name itself.
//   - Sender/recipient attribution may be inferred later using additional logic,
//     but is not guaranteed by the reason alone.
//
// These docs describe ONLY what Sisimai (and therefore CommsOK) can realistically
// detect from bounce data, not all theoretical causes.
//
// Source: https://libsisimai.org/en/reason/
const (

	// ReasonUserUnknown indicates the remote server reported that the recipient
	// mailbox does not exist or is not recognized.
	//
	// Detectable signals:
	// - SMTP replies like 550 5.1.1
	// - Diagnostic messages containing "user unknown", "no such user", etc.
	//
	// Reliable inferences:
	// - Permanent failure (hardbounce = true)
	// - Recipient address is invalid
	ReasonUserUnknown Reason = "UserUnknown"

	// ReasonHostUnknown indicates the recipient domain or mail host could not be resolved.
	//
	// Detectable signals:
	// - DNS lookup failures
	// - SMTP diagnostics like "host not found", "domain does not exist"
	//
	// Reliable inferences:
	// - Often permanent, sometimes temporary (depends on DNS failure type)
	// - Recipient domain problem
	ReasonHostUnknown Reason = "HostUnknown"

	// ReasonHasMoved indicates the recipient address is no longer valid and has
	// been replaced or moved elsewhere.
	//
	// Detectable signals:
	// - SMTP replies referencing address relocation or forwarding-only status
	//
	// Reliable inferences:
	// - Permanent failure
	// - Recipient-side condition
	ReasonHasMoved Reason = "HasMoved"

	// ReasonMailboxFull indicates the recipient mailbox is over quota.
	//
	// Detectable signals:
	// - SMTP replies like 452 4.2.2 or equivalent wording
	//
	// Reliable inferences:
	// - Usually temporary (softbounce)
	// - Recipient-side capacity issue
	ReasonMailboxFull Reason = "MailboxFull"

	// ReasonVacation indicates the recipient is in an auto-reply / vacation state
	// and delivery may be deferred.
	//
	// Detectable signals:
	// - Auto-responder wording in SMTP diagnostics
	//
	// Reliable inferences:
	// - Temporary failure
	// - No delivery rejection semantics by itself
	ReasonVacation Reason = "Vacation"

	// ReasonSpamDetected indicates the remote server classified the message as spam
	// or bulk unsolicited content.
	//
	// Detectable signals:
	// - SMTP replies referencing spam, bulk mail, content scoring
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Content/reputation related rejection
	// - Does NOT reliably prove sender-wide vs recipient-specific scope
	ReasonSpamDetected Reason = "SpamDetected"

	// ReasonBadReputation indicates rejection due to sender reputation signals.
	//
	// Detectable signals:
	// - Diagnostics mentioning reputation, history, or policy trust
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Reputation-based rejection
	ReasonBadReputation Reason = "BadReputation"

	// ReasonBlocked indicates the message was blocked by a rule or deny list.
	//
	// Detectable signals:
	// - SMTP replies containing "blocked", "denied", "refused"
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Policy-based rejection, scope unclear
	ReasonBlocked Reason = "Blocked"

	// ReasonPolicyViolation indicates rejection due to policy enforcement other
	// than spam classification.
	//
	// Detectable signals:
	// - Diagnostics referencing policy rules or restrictions
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Policy-based rejection
	ReasonPolicyViolation Reason = "PolicyViolation"

	// ReasonAuthFailure indicates authentication-related rejection.
	//
	// Detectable signals:
	// - SMTP replies mentioning SPF, DKIM, DMARC, AUTH, or credentials
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Authentication-related rejection
	ReasonAuthFailure Reason = "AuthFailure"

	// ReasonRequirePTR indicates rejection due to missing or invalid reverse DNS.
	//
	// Detectable signals:
	// - Diagnostics referencing PTR or reverse DNS
	//
	// Reliable inferences:
	// - Permanent failure
	// - Infrastructure misconfiguration
	ReasonRequirePTR Reason = "RequirePTR"

	// ReasonFailedSTARTTLS indicates TLS negotiation failure.
	//
	// Detectable signals:
	// - STARTTLS failure messages
	// - TLS handshake or certificate errors
	//
	// Reliable inferences:
	// - Temporary or permanent failure
	// - Transport security failure
	ReasonFailedSTARTTLS Reason = "FailedSTARTTLS"

	// ReasonEmailTooLarge indicates the message exceeded size limits.
	//
	// Detectable signals:
	// - SMTP replies like 552 5.3.4
	//
	// Reliable inferences:
	// - Usually permanent
	// - Message size violation
	ReasonEmailTooLarge Reason = "EmailTooLarge"

	// ReasonVirusDetected indicates malware or virus detection.
	//
	// Detectable signals:
	// - Diagnostics mentioning virus, malware, or infected content
	//
	// Reliable inferences:
	// - Permanent failure
	// - Content security rejection
	ReasonVirusDetected Reason = "VirusDetected"

	// ReasonContentError indicates malformed message content.
	//
	// Detectable signals:
	// - Invalid MIME structure
	// - Header parsing errors
	//
	// Reliable inferences:
	// - Usually permanent
	// - Message construction error
	ReasonContentError Reason = "ContentError"

	// ReasonNotCompliantRFC indicates protocol or RFC violations.
	//
	// Detectable signals:
	// - Invalid SMTP commands, malformed headers, bad addresses
	//
	// Reliable inferences:
	// - Usually permanent
	// - Protocol-level non-compliance
	ReasonNotCompliantRFC Reason = "NotCompliantRFC"

	// ReasonSyntaxError indicates syntax errors in SMTP dialogue or addresses.
	//
	// Detectable signals:
	// - SMTP syntax error responses
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Protocol syntax issue
	ReasonSyntaxError Reason = "SyntaxError"

	// ReasonNoRelaying indicates the remote server refused to relay the message.
	//
	// Detectable signals:
	// - "Relaying denied" diagnostics
	//
	// Reliable inferences:
	// - Permanent failure
	// - Relay policy rejection
	ReasonNoRelaying Reason = "NoRelaying"

	// ReasonRateLimited indicates throttling or rate enforcement.
	//
	// Detectable signals:
	// - SMTP replies like 421 or diagnostics mentioning rate limits
	//
	// Reliable inferences:
	// - Temporary failure
	// - Volume or concurrency throttling
	ReasonRateLimited Reason = "RateLimited"

	// ReasonSystemFull indicates remote system capacity exhaustion.
	//
	// Detectable signals:
	// - Queue full, disk full diagnostics
	//
	// Reliable inferences:
	// - Usually temporary
	// - Remote system capacity issue
	ReasonSystemFull Reason = "SystemFull"

	// ReasonSystemError indicates a generic remote system error.
	//
	// Detectable signals:
	// - Non-specific "system error" responses
	//
	// Reliable inferences:
	// - Usually temporary
	// - Infrastructure-level failure
	ReasonSystemError Reason = "SystemError"

	// ReasonNetworkError indicates connectivity or routing failures.
	//
	// Detectable signals:
	// - Connection timeouts
	// - TCP/IP errors
	//
	// Reliable inferences:
	// - Temporary failure
	// - Network-layer issue
	ReasonNetworkError Reason = "NetworkError"

	// ReasonExpired indicates delivery attempts were retried and eventually abandoned.
	//
	// Detectable signals:
	// - Explicit "expired" or retry timeout indicators from the MTA
	//
	// Reliable inferences:
	// - Final failure after temporary retries
	ReasonExpired Reason = "Expired"

	// ReasonFiltered indicates the message was filtered without explicit spam classification.
	//
	// Detectable signals:
	// - Diagnostics mentioning filtering or moderation
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Filtering decision, exact cause unclear
	ReasonFiltered Reason = "Filtered"

	// ReasonRejected indicates a generic rejection without detailed classification.
	//
	// Detectable signals:
	// - "Rejected" diagnostics with no additional context
	//
	// Reliable inferences:
	// - Failure occurred
	// - Cause unknown
	ReasonRejected Reason = "Rejected"

	// ReasonNotAccept indicates refusal to accept the message.
	//
	// Detectable signals:
	// - "Not accepted" diagnostics
	//
	// Reliable inferences:
	// - Failure occurred
	// - Cause unclear
	ReasonNotAccept Reason = "NotAccept"

	// ReasonMailerError indicates a generic mailer failure.
	//
	// Detectable signals:
	// - Mail system error messages without detail
	//
	// Reliable inferences:
	// - Failure occurred
	// - Classification uncertain
	ReasonMailerError Reason = "MailerError"

	// ReasonSecurityError indicates rejection due to security controls.
	//
	// Detectable signals:
	// - Diagnostics mentioning security, abuse, or protection mechanisms
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Security-related rejection
	ReasonSecurityError Reason = "SecurityError"

	// ReasonSuspend indicates the sender or recipient account is suspended.
	//
	// Detectable signals:
	// - Account disabled or suspended diagnostics
	//
	// Reliable inferences:
	// - Usually permanent
	// - Account-level restriction
	ReasonSuspend Reason = "Suspend"

	// ReasonFeedback indicates enforcement based on feedback loop data.
	//
	// Detectable signals:
	// - ESP-side suppression tied to complaints
	//
	// Reliable inferences:
	// - Not an SMTP bounce
	// - Pre-delivery suppression
	ReasonFeedback Reason = "Feedback"

	// ReasonSuppressed indicates the message was suppressed before delivery attempt.
	//
	// Detectable signals:
	// - ESP or platform reports suppression without SMTP attempt
	//
	// Reliable inferences:
	// - Not an SMTP bounce
	// - Intentional non-delivery
	ReasonSuppressed Reason = "Suppressed"

	// ReasonUndefined indicates the reason could not be determined.
	//
	// Detectable signals:
	// - Missing or unparsable diagnostics
	//
	// Reliable inferences:
	// - Failure occurred
	// - No further classification possible
	ReasonUndefined Reason = "Undefined"
)
