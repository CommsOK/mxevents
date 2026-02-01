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

	// EventMailboxFailed indicates a delivery attempt failed without a clear temp/permanent signal.
	// Example: gateway only reports "delivery failed."
	EventMailboxFailed EventType = "mailbox-failed"

	// EventMailboxTempFail indicates a temporary SMTP failure (4xx). Gateway may retry.
	EventMailboxTempFail EventType = "mailbox-tempfail"

	// EventMailboxSenderPermFail indicates a permanent SMTP failure (5xx) due to sender-side issues.
	// Examples: authentication failure, sender domain blocked, sender IP blacklisted,
	// DKIM/SPF/DMARC policy failures, sending limit exceeded.
	EventMailboxSenderPermFail EventType = "mailbox-sender-permfail"

	// EventMailboxRecipientPermFail indicates a permanent SMTP failure (5xx) due to recipient-side issues.
	// Examples: mailbox does not exist, mailbox disabled, recipient domain does not exist,
	// recipient rejected message. This is the default when failure type is uncertain.
	EventMailboxRecipientPermFail EventType = "mailbox-recipient-permfail"

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

	// BounceReasonUserUnknown indicates the remote server reported that the recipient
	// mailbox does not exist or is not recognized.
	//
	// Detectable signals:
	// - SMTP replies like 550 5.1.1
	// - Diagnostic messages containing "user unknown", "no such user", etc.
	//
	// Reliable inferences:
	// - Permanent failure (hardbounce = true)
	// - Recipient address is invalid
	BounceReasonUserUnknown Reason = "UserUnknown"

	// BounceReasonHostUnknown indicates the recipient domain or mail host could not be resolved.
	//
	// Detectable signals:
	// - DNS lookup failures
	// - SMTP diagnostics like "host not found", "domain does not exist"
	//
	// Reliable inferences:
	// - Often permanent, sometimes temporary (depends on DNS failure type)
	// - Recipient domain problem
	BounceReasonHostUnknown Reason = "HostUnknown"

	// BounceReasonHasMoved indicates the recipient address is no longer valid and has
	// been replaced or moved elsewhere.
	//
	// Detectable signals:
	// - SMTP replies referencing address relocation or forwarding-only status
	//
	// Reliable inferences:
	// - Permanent failure
	// - Recipient-side condition
	BounceReasonHasMoved Reason = "HasMoved"

	// BounceReasonMailboxFull indicates the recipient mailbox is over quota.
	//
	// Detectable signals:
	// - SMTP replies like 452 4.2.2 or equivalent wording
	//
	// Reliable inferences:
	// - Usually temporary (softbounce)
	// - Recipient-side capacity issue
	BounceReasonMailboxFull Reason = "MailboxFull"

	// BounceReasonVacation indicates the recipient is in an auto-reply / vacation state
	// and delivery may be deferred.
	//
	// Detectable signals:
	// - Auto-responder wording in SMTP diagnostics
	//
	// Reliable inferences:
	// - Temporary failure
	// - No delivery rejection semantics by itself
	BounceReasonVacation Reason = "Vacation"

	// BounceReasonSpamDetected indicates the remote server classified the message as spam
	// or bulk unsolicited content.
	//
	// Detectable signals:
	// - SMTP replies referencing spam, bulk mail, content scoring
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Content/reputation related rejection
	// - Does NOT reliably prove sender-wide vs recipient-specific scope
	BounceReasonSpamDetected Reason = "SpamDetected"

	// BounceReasonBadReputation indicates rejection due to sender reputation signals.
	//
	// Detectable signals:
	// - Diagnostics mentioning reputation, history, or policy trust
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Reputation-based rejection
	BounceReasonBadReputation Reason = "BadReputation"

	// BounceReasonBlocked indicates the message was blocked by a rule or deny list.
	//
	// Detectable signals:
	// - SMTP replies containing "blocked", "denied", "refused"
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Policy-based rejection, scope unclear
	BounceReasonBlocked Reason = "Blocked"

	// BounceReasonPolicyViolation indicates rejection due to policy enforcement other
	// than spam classification.
	//
	// Detectable signals:
	// - Diagnostics referencing policy rules or restrictions
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Policy-based rejection
	BounceReasonPolicyViolation Reason = "PolicyViolation"

	// BounceReasonAuthFailure indicates authentication-related rejection.
	//
	// Detectable signals:
	// - SMTP replies mentioning SPF, DKIM, DMARC, AUTH, or credentials
	//
	// Reliable inferences:
	// - Failure may be temporary or permanent
	// - Authentication-related rejection
	BounceReasonAuthFailure Reason = "AuthFailure"

	// BounceReasonRequirePTR indicates rejection due to missing or invalid reverse DNS.
	//
	// Detectable signals:
	// - Diagnostics referencing PTR or reverse DNS
	//
	// Reliable inferences:
	// - Permanent failure
	// - Infrastructure misconfiguration
	BounceReasonRequirePTR Reason = "RequirePTR"

	// BounceReasonFailedSTARTTLS indicates TLS negotiation failure.
	//
	// Detectable signals:
	// - STARTTLS failure messages
	// - TLS handshake or certificate errors
	//
	// Reliable inferences:
	// - Temporary or permanent failure
	// - Transport security failure
	BounceReasonFailedSTARTTLS Reason = "FailedSTARTTLS"

	// BounceReasonEmailTooLarge indicates the message exceeded size limits.
	//
	// Detectable signals:
	// - SMTP replies like 552 5.3.4
	//
	// Reliable inferences:
	// - Usually permanent
	// - Message size violation
	BounceReasonEmailTooLarge Reason = "EmailTooLarge"

	// BounceReasonVirusDetected indicates malware or virus detection.
	//
	// Detectable signals:
	// - Diagnostics mentioning virus, malware, or infected content
	//
	// Reliable inferences:
	// - Permanent failure
	// - Content security rejection
	BounceReasonVirusDetected Reason = "VirusDetected"

	// BounceReasonContentError indicates malformed message content.
	//
	// Detectable signals:
	// - Invalid MIME structure
	// - Header parsing errors
	//
	// Reliable inferences:
	// - Usually permanent
	// - Message construction error
	BounceReasonContentError Reason = "ContentError"

	// BounceReasonNotCompliantRFC indicates protocol or RFC violations.
	//
	// Detectable signals:
	// - Invalid SMTP commands, malformed headers, bad addresses
	//
	// Reliable inferences:
	// - Usually permanent
	// - Protocol-level non-compliance
	BounceReasonNotCompliantRFC Reason = "NotCompliantRFC"

	// BounceReasonSyntaxError indicates syntax errors in SMTP dialogue or addresses.
	//
	// Detectable signals:
	// - SMTP syntax error responses
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Protocol syntax issue
	BounceReasonSyntaxError Reason = "SyntaxError"

	// BounceReasonNoRelaying indicates the remote server refused to relay the message.
	//
	// Detectable signals:
	// - "Relaying denied" diagnostics
	//
	// Reliable inferences:
	// - Permanent failure
	// - Relay policy rejection
	BounceReasonNoRelaying Reason = "NoRelaying"

	// BounceReasonRateLimited indicates throttling or rate enforcement.
	//
	// Detectable signals:
	// - SMTP replies like 421 or diagnostics mentioning rate limits
	//
	// Reliable inferences:
	// - Temporary failure
	// - Volume or concurrency throttling
	BounceReasonRateLimited Reason = "RateLimited"

	// BounceReasonSystemFull indicates remote system capacity exhaustion.
	//
	// Detectable signals:
	// - Queue full, disk full diagnostics
	//
	// Reliable inferences:
	// - Usually temporary
	// - Remote system capacity issue
	BounceReasonSystemFull Reason = "SystemFull"

	// BounceReasonSystemError indicates a generic remote system error.
	//
	// Detectable signals:
	// - Non-specific "system error" responses
	//
	// Reliable inferences:
	// - Usually temporary
	// - Infrastructure-level failure
	BounceReasonSystemError Reason = "SystemError"

	// BounceReasonNetworkError indicates connectivity or routing failures.
	//
	// Detectable signals:
	// - Connection timeouts
	// - TCP/IP errors
	//
	// Reliable inferences:
	// - Temporary failure
	// - Network-layer issue
	BounceReasonNetworkError Reason = "NetworkError"

	// BounceReasonExpired indicates delivery attempts were retried and eventually abandoned.
	//
	// Detectable signals:
	// - Explicit "expired" or retry timeout indicators from the MTA
	//
	// Reliable inferences:
	// - Final failure after temporary retries
	BounceReasonExpired Reason = "Expired"

	// BounceReasonFiltered indicates the message was filtered without explicit spam classification.
	//
	// Detectable signals:
	// - Diagnostics mentioning filtering or moderation
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Filtering decision, exact cause unclear
	BounceReasonFiltered Reason = "Filtered"

	// BounceReasonRejected indicates a generic rejection without detailed classification.
	//
	// Detectable signals:
	// - "Rejected" diagnostics with no additional context
	//
	// Reliable inferences:
	// - Failure occurred
	// - Cause unknown
	BounceReasonRejected Reason = "Rejected"

	// BounceReasonNotAccept indicates refusal to accept the message.
	//
	// Detectable signals:
	// - "Not accepted" diagnostics
	//
	// Reliable inferences:
	// - Failure occurred
	// - Cause unclear
	BounceReasonNotAccept Reason = "NotAccept"

	// BounceReasonMailerError indicates a generic mailer failure.
	//
	// Detectable signals:
	// - Mail system error messages without detail
	//
	// Reliable inferences:
	// - Failure occurred
	// - Classification uncertain
	BounceReasonMailerError Reason = "MailerError"

	// BounceReasonSecurityError indicates rejection due to security controls.
	//
	// Detectable signals:
	// - Diagnostics mentioning security, abuse, or protection mechanisms
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Security-related rejection
	BounceReasonSecurityError Reason = "SecurityError"

	// BounceReasonSuspend indicates the sender or recipient account is suspended.
	//
	// Detectable signals:
	// - Account disabled or suspended diagnostics
	//
	// Reliable inferences:
	// - Usually permanent
	// - Account-level restriction
	BounceReasonSuspend Reason = "Suspend"

	// BounceReasonFeedback indicates enforcement based on feedback loop data.
	//
	// Detectable signals:
	// - ESP-side suppression tied to complaints
	//
	// Reliable inferences:
	// - Not an SMTP bounce
	// - Pre-delivery suppression
	BounceReasonFeedback Reason = "Feedback"

	// BounceReasonSuppressed indicates the message was suppressed before delivery attempt.
	//
	// Detectable signals:
	// - ESP or platform reports suppression without SMTP attempt
	//
	// Reliable inferences:
	// - Not an SMTP bounce
	// - Intentional non-delivery
	BounceReasonSuppressed Reason = "Suppressed"

	// BounceReasonUndefined indicates the reason could not be determined.
	//
	// Detectable signals:
	// - Missing or unparsable diagnostics
	//
	// Reliable inferences:
	// - Failure occurred
	// - No further classification possible
	BounceReasonUndefined Reason = "Undefined"
)

// ====================================
// Drop Reasons (Pre-SMTP, Intentional Non-Delivery)
// ====================================
//
// These reasons explain why a message was intentionally not sent
// by the originating platform or delivery gateway.
//
// They are NOT SMTP bounces and must never be treated as mailbox failures.
type DropReason Reason

const (
	// DropReasonSuppressed indicates the message was suppressed due to an
	// internal suppression list or prior decision.
	//
	// Actionability:
	// - No new suppression
	// - Confirms existing suppression enforcement
	DropReasonSuppressed DropReason = "Suppressed"

	// DropReasonPolicyViolation indicates the message violated platform
	// or gateway policy (not recipient mailbox policy).
	//
	// Actionability:
	// - Block sending until configuration or content changes
	// - Escalate to operator
	DropReasonPolicyViolation DropReason = "PolicyViolation"

	// DropReasonRateLimited indicates the message was dropped due to
	// rate, concurrency, or throughput limits.
	//
	// Actionability:
	// - Retry later
	// - No suppression
	DropReasonRateLimited DropReason = "RateLimited"

	// DropReasonFeedback indicates the message was dropped due to prior
	// spam complaints or feedback loop enforcement.
	//
	// Actionability:
	// - Global suppression
	// - Compliance escalation
	DropReasonFeedback DropReason = "Feedback"

	// DropReasonSecurity indicates the message was dropped due to security
	// concerns (account abuse, compromised credentials, etc.).
	//
	// Actionability:
	// - Immediate stop
	// - Manual review required
	DropReasonSecurity DropReason = "Security"

	// DropReasonConfiguration indicates the platform could not send the
	// message due to missing or invalid configuration.
	//
	// Actionability:
	// - Fix integration, auth, domain, or sender setup
	// - No suppression
	DropReasonConfiguration DropReason = "Configuration"
)

// ====================================
// Unsubscribe Reasons
// ====================================
type UnsubscribeReason Reason

const (
	// UnsubscribeReasonUserInitiated indicates the recipient explicitly
	// opted out (UI, link, reply).
	//
	// Actionability:
	// - Permanent suppression
	// - No rehabilitation
	UnsubscribeReasonUserInitiated UnsubscribeReason = "UserInitiated"

	// UnsubscribeReasonOneClick indicates RFC-compliant one-click unsubscribe.
	//
	// Actionability:
	// - Permanent suppression
	// - Compliance-grade signal
	UnsubscribeReasonOneClick UnsubscribeReason = "OneClick"

	// UnsubscribeReasonAdmin indicates an operator or admin unsubscribed
	// the recipient.
	//
	// Actionability:
	// - Suppression applies
	// - Rehabilitation possible if reversed
	UnsubscribeReasonAdmin UnsubscribeReason = "AdminAction"

	// UnsubscribeReasonCompliance indicates automatic unsubscribe due to
	// legal or regulatory enforcement.
	//
	// Actionability:
	// - Permanent suppression
	// - No override
	UnsubscribeReasonCompliance UnsubscribeReason = "Compliance"
)

// ====================================
// Spam Report Reasons
// ====================================
type SpamReportReason Reason

const (
	// SpamReportReasonFeedbackLoop indicates ISP feedback loop report.
	//
	// Actionability:
	// - Immediate global suppression
	// - Reputation impact
	SpamReportReasonFeedbackLoop SpamReportReason = "FeedbackLoop"

	// SpamReportReasonUserAction indicates the recipient explicitly
	// marked the message as spam.
	//
	// Actionability:
	// - Immediate suppression
	// - High confidence signal
	SpamReportReasonUserAction SpamReportReason = "UserAction"

	// SpamReportReasonProviderHeuristic indicates the provider inferred
	// spam reporting behavior without explicit user action.
	//
	// Actionability:
	// - Suppression recommended
	// - Lower confidence
	SpamReportReasonProviderHeuristic SpamReportReason = "ProviderHeuristic"
)

// ====================================
// Failure Reasons (Non-SMTP, Operational)
// ====================================
type FailureReason Reason

const (
	// FailureReasonConfiguration indicates invalid or missing setup.
	//
	// Actionability:
	// - Fix integration or credentials
	FailureReasonConfiguration FailureReason = "Configuration"

	// FailureReasonAuth indicates authentication or authorization failure.
	//
	// Actionability:
	// - Rotate credentials
	// - Fix OAuth / API keys
	FailureReasonAuth FailureReason = "Auth"

	// FailureReasonSystem indicates internal platform or gateway failure.
	//
	// Actionability:
	// - Retry
	// - Monitor incident
	FailureReasonSystem FailureReason = "System"

	// FailureReasonTimeout indicates upstream or downstream timeout.
	//
	// Actionability:
	// - Retry
	// - No suppression
	FailureReasonTimeout FailureReason = "Timeout"
)
