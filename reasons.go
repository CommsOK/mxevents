package mxevents

import "strings"

// ====================================
// Reason Taxonomy (namespace.code)
// ====================================
//
// This package defines a canonical Reason taxonomy for email event classification.
//
// Design goals:
//
//   - Reasons clarify an EventType; they are not intended to be interpreted without the EventType.
//   - The Reason string uses a structured "namespace.code" format to keep entropy balanced:
//   - namespace: low-cardinality bucket indicating the broad mechanism (smtp, network, suppression, etc.)
//   - code: UpperCamelCase token indicating the specific classification within that namespace
//   - Detailed vendor diagnostics belong in a separate free-form field (e.g., reason_message / diagnostics),
//     not in the canonical taxonomy.
//   - The taxonomy is designed to map cleanly to common ESP/CRM signals and to Sisimai-style bounce parsing.
//
// Format rules:
//
//   - Reason is a string literal: "<namespace>.<Code>"
//   - namespace is lowercase ASCII (e.g., "smtp", "network").
//   - Code is UpperCamelCase ASCII (e.g., "UserUnknown", "TlsFailure").
//   - '.' is the only separator.
//   - Reasons must be stable once published.
//
// Examples:
//
//   - EventMailboxRecipientPermFail + "smtp.UserUnknown"
//   - EventMailboxFailed           + "network.Timeout"
//   - EventGatewayDropped          + "suppression.Suppressed"
//   - EventStatusSpamReported      + "complaint.FeedbackLoop"
//
// Notes on semantics:
//
//   - Namespaces are the primary unit for high-level actionability.
//   - Codes provide additional clarity, but the long tail remains in diagnostics.
//   - Same Code tokens MAY exist across different namespaces; (namespace, code) is the identity.
//     Do not implement behavior keyed only on the Code token.
//
// ====================================
// Reason Namespace Catalog
// ====================================
//
// smtp.*
//
//	SMTP / mailbox-side rejection or deferral. Usually sourced from SMTP replies,
//	enhanced status codes (when present), and diagnostic parsing (often via Sisimai).
//
// network.*
//
//	DNS, connectivity, timeout, routing, and TLS negotiation failures where an SMTP-level
//	response may be missing or unreliable.
//
// suppression.*
//
//	Intentional non-delivery due to known recipient state, policy enforcement, or explicit
//	operator decisions. These are not mailbox bounces and must never be treated as SMTP failures.
//
// policy.*
//
//	Policy and reputation enforcement where the mechanism is policy-based but may not be
//	safe to interpret as a specific SMTP reason (or is reported as policy at a platform layer).
//
// config.*
//
//	Missing/invalid configuration that prevents processing or delivery attempts.
//
// auth.*
//
//	Authentication/authorization failures with integrations, APIs, tokens, or sender auth.
//
// system.*
//
//	Platform or gateway operational errors (internal failures, transient unavailability, backpressure).
//
// complaint.*
//
//	Spam complaints and feedback loop signals. Not SMTP bounces.
//
// unsubscribe.*
//
//	Unsubscribe mechanisms and sources. Not SMTP bounces.
//
// unknown.*
//
//	Fallback bucket when no reliable classification is possible.
//
// ====================================
// SMTP Reasons (smtp.*)
// ====================================
//
// These Reasons describe why a remote system rejected or deferred an SMTP delivery attempt,
// as inferred from SMTP reply codes, enhanced status codes, and diagnostic messages.
//
// Important constraints:
//
//   - smtp.* Reasons describe *why* a delivery attempt failed or was deferred.
//   - Temporary vs permanent is determined primarily via SMTP 4xx/5xx or a hardbounce signal
//     (when available), not by the Reason name alone.
//   - Sender/recipient attribution is represented by EventType (sender-permfail vs recipient-permfail),
//     not by the Reason alone.
//
// Docs describe what common providers and parsers (including Sisimai-style parsing) can realistically detect.
const (
	// ReasonSMTPUserUnknown indicates the remote server reported that the recipient mailbox
	// does not exist or is not recognized.
	//
	// Detectable signals:
	// - SMTP replies like 550 5.1.1
	// - Diagnostics containing "user unknown", "no such user", etc.
	//
	// Reliable inferences:
	// - Often permanent (hardbounce)
	// - Recipient address invalid
	ReasonSMTPUserUnknown Reason = "smtp.UserUnknown"

	// ReasonSMTPMailboxFull indicates the recipient mailbox is over quota.
	//
	// Detectable signals:
	// - SMTP replies like 452 4.2.2 or quota wording
	//
	// Reliable inferences:
	// - Usually temporary (softbounce)
	// - Recipient-side capacity issue
	ReasonSMTPMailboxFull Reason = "smtp.MailboxFull"

	// ReasonSMTPHasMoved indicates the recipient address is no longer valid and has moved.
	//
	// Detectable signals:
	// - Relocation / forwarding-only wording
	//
	// Reliable inferences:
	// - Usually permanent
	// - Recipient-side condition
	ReasonSMTPHasMoved Reason = "smtp.HasMoved"

	// ReasonSMTPSpamDetected indicates the remote server classified the message as spam/bulk/unsolicited.
	//
	// Detectable signals:
	// - SMTP diagnostics referencing spam, bulk, content scoring
	//
	// Reliable inferences:
	// - May be temporary or permanent
	// - Content/reputation related rejection
	// - Scope (global vs recipient-specific) is not guaranteed
	ReasonSMTPSpamDetected Reason = "smtp.SpamDetected"

	// ReasonSMTPPolicyViolation indicates rejection due to remote policy enforcement other than explicit spam classification.
	//
	// Detectable signals:
	// - Diagnostics referencing policy rules/restrictions
	//
	// Reliable inferences:
	// - May be temporary or permanent
	// - Policy-based rejection, scope unclear
	ReasonSMTPPolicyViolation Reason = "smtp.PolicyViolation"

	// ReasonSMTPAuthFailure indicates authentication-related rejection during SMTP evaluation.
	//
	// Detectable signals:
	// - Diagnostics mentioning SPF, DKIM, DMARC, AUTH, rDNS requirements tied to auth posture
	//
	// Reliable inferences:
	// - May be temporary or permanent
	// - Sender-side authentication posture problem
	ReasonSMTPAuthFailure Reason = "smtp.AuthFailure"

	// ReasonSMTPBlocked indicates the message was blocked/refused/denied by the remote server.
	//
	// Detectable signals:
	// - SMTP replies containing "blocked", "denied", "refused"
	//
	// Reliable inferences:
	// - May be temporary or permanent
	// - Policy/security-based rejection, scope unclear
	ReasonSMTPBlocked Reason = "smtp.Blocked"

	// ReasonSMTPRateLimited indicates throttling or rate enforcement by the remote system.
	//
	// Detectable signals:
	// - SMTP replies like 421/451 with rate limit wording
	// - Diagnostics mentioning "too many connections" or "try again later"
	//
	// Reliable inferences:
	// - Usually temporary
	// - Volume/concurrency throttling
	ReasonSMTPRateLimited Reason = "smtp.RateLimited"

	// ReasonSMTPMessageTooLarge indicates message size exceeded remote limits.
	//
	// Detectable signals:
	// - SMTP replies like 552 5.3.4 or "message too large"
	//
	// Reliable inferences:
	// - Usually permanent until content changes
	ReasonSMTPMessageTooLarge Reason = "smtp.MessageTooLarge"

	// ReasonSMTPVirusDetected indicates rejection due to malware/virus detection.
	//
	// Detectable signals:
	// - Diagnostics mentioning virus/malware/infected content
	//
	// Reliable inferences:
	// - Usually permanent
	// - Content security rejection
	ReasonSMTPVirusDetected Reason = "smtp.VirusDetected"

	// ReasonSMTPSyntaxError indicates SMTP dialogue or address syntax error.
	//
	// Detectable signals:
	// - SMTP syntax error responses
	// - Diagnostics mentioning malformed addresses/commands
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Protocol/address validity issue
	ReasonSMTPSyntaxError Reason = "smtp.SyntaxError"

	// ReasonSMTPNoRelaying indicates the remote server refused to relay the message.
	//
	// Detectable signals:
	// - "Relaying denied" diagnostics
	//
	// Reliable inferences:
	// - Usually permanent
	// - Relay policy rejection
	ReasonSMTPNoRelaying Reason = "smtp.NoRelaying"

	// ReasonSMTPExpired indicates delivery was retried and eventually abandoned.
	//
	// Detectable signals:
	// - Explicit "expired" indicators
	// - Retry timeout indicators from MTA/ESP
	//
	// Reliable inferences:
	// - Final failure after temporary retries
	ReasonSMTPExpired Reason = "smtp.Expired"

	// ReasonSMTPContentError indicates the message was rejected due to malformed
	// or non-compliant message content.
	//
	// Detectable signals:
	// - Invalid or malformed MIME structure
	// - Missing or invalid headers
	// - RFC compliance violations reported by the remote server
	//
	// Reliable inferences:
	// - Usually permanent until message construction is corrected
	ReasonSMTPContentError Reason = "smtp.ContentError"

	// ReasonNetworkError indicates a delivery attempt failed due to a network-related
	// issue that could not be reliably classified further.
	//
	// Detectable signals:
	// - Generic network failure diagnostics
	// - Sisimai reason "networkerror"
	//
	// Reliable inferences:
	// - Usually temporary
	// - Retry may succeed
	ReasonNetworkError Reason = "network.Error"
)

// ====================================
// Network Reasons (network.*)
// ====================================
//
// These Reasons describe delivery attempt failures that occur before a reliable SMTP outcome
// can be obtained, or where the failure is clearly transport-related.
//
// They are commonly used with EventMailboxFailed, and sometimes with temporary failures
// when the provider signals transport issues.
const (
	// ReasonNetworkDnsFailure indicates DNS resolution failure for recipient domain or mail host.
	//
	// Detectable signals:
	// - NXDOMAIN/SERVFAIL for domain or MX
	// - Diagnostics like "host unknown", "domain does not exist", "no MX"
	//
	// Reliable inferences:
	// - Often permanent, sometimes temporary (depends on failure type)
	// - Recipient domain / DNS problem
	ReasonNetworkDnsFailure Reason = "network.DnsFailure"

	// ReasonNetworkConnectFailure indicates connection failure to remote host.
	//
	// Detectable signals:
	// - Connection refused/reset/unreachable
	// - TCP-level errors
	//
	// Reliable inferences:
	// - Usually temporary
	// - Transport/connectivity issue
	ReasonNetworkConnectFailure Reason = "network.ConnectFailure"

	// ReasonNetworkTimeout indicates timeout during dial/read/write while attempting delivery.
	//
	// Detectable signals:
	// - Dial timeout / read timeout diagnostics
	//
	// Reliable inferences:
	// - Usually temporary
	// - Network-layer delay/outage
	ReasonNetworkTimeout Reason = "network.Timeout"

	// ReasonNetworkTlsFailure indicates TLS/STARTTLS negotiation failure.
	//
	// Detectable signals:
	// - STARTTLS failure messages
	// - TLS handshake/certificate errors
	//
	// Reliable inferences:
	// - Temporary or permanent
	// - Transport security negotiation failure
	ReasonNetworkTlsFailure Reason = "network.TlsFailure"

	// ReasonNetworkRoutingError indicates routing issues (network unreachable, route errors).
	//
	// Detectable signals:
	// - "network unreachable", routing diagnostics
	//
	// Reliable inferences:
	// - Usually temporary
	// - Network routing issue
	ReasonNetworkRoutingError Reason = "network.RoutingError"
)

// ====================================
// Suppression Reasons (suppression.*)
// ====================================
//
// These Reasons indicate intentional non-delivery. They are not SMTP bounces and must never
// be treated as mailbox failures.
//
// They commonly attach to EventOriginDropped or EventGatewayDropped.
const (
	// ReasonSuppressionSuppressed indicates the message was suppressed due to an existing suppression
	// list entry or prior decision.
	//
	// Actionability:
	// - No new suppression
	// - Confirms enforcement of an existing suppression state
	ReasonSuppressionSuppressed Reason = "suppression.Suppressed"

	// ReasonSuppressionUnsubscribed indicates the message was suppressed due to an unsubscribe state.
	//
	// Actionability:
	// - Permanent suppression
	// - No rehabilitation unless unsubscribe is explicitly reversed (if supported)
	ReasonSuppressionUnsubscribed Reason = "suppression.Unsubscribed"

	// ReasonSuppressionSpamComplaint indicates the message was suppressed due to prior spam complaint enforcement.
	//
	// Actionability:
	// - Immediate global suppression
	// - Compliance workflow recommended
	ReasonSuppressionSpamComplaint Reason = "suppression.SpamComplaint"

	// ReasonSuppressionPolicyBlocked indicates the message was suppressed due to local platform/gateway policy.
	//
	// Actionability:
	// - Block sending until policy/config changes
	// - Manual review may be required
	ReasonSuppressionPolicyBlocked Reason = "suppression.PolicyBlocked"

	// ReasonSuppressionAdminBlocked indicates an operator/admin explicitly blocked delivery.
	//
	// Actionability:
	// - Suppression applies
	// - May be reversible depending on operator policy
	ReasonSuppressionAdminBlocked Reason = "suppression.AdminBlocked"

	// ReasonSuppressionCleared indicates the recipient was removed from a suppression list.
	// This is used for "unsuppression" events like UNBOUNCE where an address is reactivated.
	//
	// Actionability:
	// - Suppression removed
	// - Recipient is eligible for sending again
	// - Commonly used with EventStatusSubscribed to indicate reactivation
	ReasonSuppressionCleared Reason = "suppression.Cleared"
)

// ====================================
// Policy Reasons (policy.*)
// ====================================
//
// These Reasons represent policy or reputation enforcement. They may be emitted by providers/ESPs/CRMs
// without sufficient detail to safely classify as a specific SMTP diagnostic.
const (
	// ReasonPolicyBlocked indicates policy enforcement or deny-list style blocking.
	//
	// Detectable signals:
	// - Vendor status like "blocked" / "policy" with limited diagnostics
	// - Provider signal that a policy rule prevented acceptance
	//
	// Reliable inferences:
	// - Blocking is policy-based
	// - Scope may be unclear (global vs recipient-specific)
	ReasonPolicyBlocked Reason = "policy.Blocked"

	// ReasonPolicyBadReputation indicates rejection associated with sender reputation signals.
	//
	// Detectable signals:
	// - Diagnostics mentioning reputation/trust/history
	//
	// Reliable inferences:
	// - Reputation-based enforcement
	// - Scope and permanence may vary
	ReasonPolicyBadReputation Reason = "policy.BadReputation"

	// ReasonPolicyRestricted indicates content/recipient/sender category restrictions.
	//
	// Detectable signals:
	// - "restricted", "not permitted", "not allowed" type diagnostics
	//
	// Reliable inferences:
	// - Policy restriction exists
	ReasonPolicyRestricted Reason = "policy.Restricted"
)

// ====================================
// Configuration Reasons (config.*)
// ====================================
//
// These Reasons indicate missing or invalid configuration preventing processing or delivery attempts.
// They commonly attach to EventOriginFailed, EventOriginDropped, EventGatewayFailed, or EventGatewayDropped.
const (
	// ReasonConfigMissing indicates required configuration is absent.
	//
	// Examples:
	// - Missing sender identity, domain verification, template, routing config
	//
	// Actionability:
	// - Fix configuration
	// - No suppression
	ReasonConfigMissing Reason = "config.Missing"

	// ReasonConfigInvalid indicates configuration exists but is invalid.
	//
	// Examples:
	// - Invalid sender domain, invalid payload fields, incompatible settings
	//
	// Actionability:
	// - Fix configuration
	// - No suppression
	ReasonConfigInvalid Reason = "config.Invalid"

	// ReasonConfigInvalidRecipient indicates the recipient address is malformed or unparseable.
	//
	// Examples:
	// - Malformed email address syntax
	// - Missing or invalid to/cc/bcc fields
	//
	// Actionability:
	// - Fix recipient address in source data
	// - No suppression
	ReasonConfigInvalidRecipient Reason = "config.InvalidRecipient"

	// ReasonConfigInvalidMessage indicates the message structure is invalid.
	//
	// Examples:
	// - Missing required fields (subject, body)
	// - Invalid content type or encoding
	// - Malformed attachments
	//
	// Actionability:
	// - Fix message construction in source system
	// - No suppression
	ReasonConfigInvalidMessage Reason = "config.InvalidMessage"
)

// ====================================
// Authentication Reasons (auth.*)
// ====================================
//
// These Reasons indicate API/integration authentication or authorization failures (non-SMTP).
// They commonly attach to EventOriginFailed or EventGatewayFailed.
const (
	// ReasonAuthUnauthorized indicates authorization failure (insufficient scopes/permissions, forbidden).
	//
	// Actionability:
	// - Fix permissions/scopes
	// - No suppression
	ReasonAuthUnauthorized Reason = "auth.Unauthorized"

	// ReasonAuthInvalid indicates invalid credentials or authentication token.
	//
	// Actionability:
	// - Rotate/refresh credentials
	// - Fix OAuth / API key / token validity
	ReasonAuthInvalid Reason = "auth.Invalid"
)

// ====================================
// System Reasons (system.*)
// ====================================
//
// These Reasons indicate platform or gateway operational failures. They are not mailbox semantics.
// They may attach to EventOriginFailed, EventGatewayFailed, or EventMailboxFailed (when failure is internal).
const (
	// ReasonSystemError indicates a generic internal error.
	//
	// Actionability:
	// - Retry may succeed
	// - Monitor incident and logs
	ReasonSystemError Reason = "system.Error"

	// ReasonSystemUnavailable indicates an upstream dependency or service is unavailable.
	//
	// Actionability:
	// - Retry later
	// - Monitor upstream health
	ReasonSystemUnavailable Reason = "system.Unavailable"

	// ReasonSystemRateLimited indicates local throttling/backpressure at the platform/gateway layer.
	//
	// Actionability:
	// - Retry later
	// - No suppression
	ReasonSystemRateLimited Reason = "system.RateLimited"

	// ReasonSystemTimeout indicates internal or upstream timeout not clearly attributable to network SMTP transport.
	//
	// Actionability:
	// - Retry
	// - No suppression
	ReasonSystemTimeout Reason = "system.Timeout"

	// ReasonSystemBackpressure indicates queue/throughput pressure (queue full, resource constraints).
	//
	// Actionability:
	// - Retry later
	// - Consider scaling/limits tuning
	ReasonSystemBackpressure Reason = "system.Backpressure"
)

// ====================================
// Complaint Reasons (complaint.*)
// ====================================
//
// These Reasons represent spam complaint signals. They are not SMTP bounces and should be treated as
// compliance/suppression signals.
const (
	// ReasonComplaintFeedbackLoop indicates ISP feedback loop complaint.
	//
	// Actionability:
	// - Immediate global suppression
	// - Reputation impact tracking
	ReasonComplaintFeedbackLoop Reason = "complaint.FeedbackLoop"

	// ReasonComplaintUserAction indicates a direct user action marking a message as spam.
	//
	// Actionability:
	// - Immediate suppression recommended
	// - High confidence complaint signal
	ReasonComplaintUserAction Reason = "complaint.UserAction"

	// ReasonComplaintProviderHeuristic indicates provider/platform inferred spam complaint behavior
	// without a direct user action signal.
	//
	// Actionability:
	// - Suppression recommended
	// - Lower confidence than explicit complaint
	ReasonComplaintProviderHeuristic Reason = "complaint.ProviderHeuristic"
)

// ====================================
// Unsubscribe Reasons (unsubscribe.*)
// ====================================
//
// These Reasons represent unsubscribe mechanisms/sources. They are not SMTP bounces.
const (
	// ReasonUnsubscribeOneClick indicates RFC-compliant one-click unsubscribe.
	//
	// Actionability:
	// - Permanent suppression
	// - Compliance-grade signal
	ReasonUnsubscribeOneClick Reason = "unsubscribe.OneClick"

	// ReasonUnsubscribeUserInitiated indicates recipient explicitly opted out via link/UI/reply.
	//
	// Actionability:
	// - Permanent suppression
	// - No rehabilitation unless explicitly reversed (if supported)
	ReasonUnsubscribeUserInitiated Reason = "unsubscribe.UserInitiated"

	// ReasonUnsubscribeAdminAction indicates an operator/admin unsubscribed the recipient.
	//
	// Actionability:
	// - Suppression applies
	// - Rehabilitation possible if reversed
	ReasonUnsubscribeAdminAction Reason = "unsubscribe.AdminAction"

	// ReasonUnsubscribeCompliance indicates automatic unsubscribe due to legal/regulatory enforcement.
	//
	// Actionability:
	// - Permanent suppression
	// - No override
	ReasonUnsubscribeCompliance Reason = "unsubscribe.Compliance"
)

// ====================================
// Unknown Reasons (unknown.*)
// ====================================
//
// Use when classification is not possible or diagnostics are missing/unreliable.
const (
	// ReasonUnknown indicates the reason could not be determined.
	//
	// Detectable signals:
	// - Missing or unparsable diagnostics
	//
	// Reliable inferences:
	// - A failure/rejection occurred
	// - No further classification possible
	ReasonUnknown Reason = "unknown.Unknown"
)

// ====================================
// Reason Compatibility (EventType -> Allowed Namespaces)
// ====================================
//
// Implementations should validate that a Reason belongs to an allowed namespace for a given EventType.
// This keeps "switch" logic manageable and prevents invalid combinations.
//
// Notes:
//   - Success events generally SHOULD NOT carry Reasons.
//   - EventMailboxFailed is reserved for cases without reliable SMTP-level outcome;
//     it commonly pairs with network.* and system.* reasons.
type ReasonNamespace string

const (
	NSMTP         ReasonNamespace = "smtp"
	NSNetwork     ReasonNamespace = "network"
	NSSuppression ReasonNamespace = "suppression"
	NSPolicy      ReasonNamespace = "policy"
	NSConfig      ReasonNamespace = "config"
	NSAuth        ReasonNamespace = "auth"
	NSSystem      ReasonNamespace = "system"
	NSComplaint   ReasonNamespace = "complaint"
	NSUnsubscribe ReasonNamespace = "unsubscribe"
	NSUnknown     ReasonNamespace = "unknown"
)

var AllowedReasonNamespacesByEventType = map[EventType][]ReasonNamespace{
	// Origin platform
	EventOriginFailed:  {NSConfig, NSAuth, NSSystem, NSPolicy, NSUnknown},
	EventOriginDropped: {NSSuppression, NSPolicy, NSConfig, NSSystem, NSUnknown},

	// Gateway (ESP layer)
	EventGatewayFailed:  {NSConfig, NSAuth, NSSystem, NSPolicy, NSUnknown},
	EventGatewayDropped: {NSSuppression, NSPolicy, NSSystem, NSUnknown},

	// Mailbox
	EventMailboxFailed:            {NSNetwork, NSSystem, NSUnknown},
	EventMailboxTempFail:          {NSMTP, NSNetwork, NSSystem, NSPolicy, NSUnknown},
	EventMailboxPermFail:          {NSMTP, NSPolicy, NSUnknown},
	EventMailboxSenderPermFail:    {NSMTP, NSPolicy, NSAuth, NSConfig, NSUnknown},
	EventMailboxRecipientPermFail: {NSMTP, NSUnknown},

	// Status
	EventStatusSubscribed:        {NSSuppression, NSSystem, NSUnknown}, // suppression.Cleared for reactivation
	EventStatusGroupSubscribed:   {NSSuppression, NSSystem, NSUnknown}, // suppression.Cleared for reactivation
	EventStatusUnsubscribed:      {NSUnsubscribe, NSSuppression, NSUnknown},
	EventStatusGroupUnsubscribed: {NSUnsubscribe, NSSuppression, NSUnknown},
	EventStatusSpamReported:      {NSComplaint, NSSuppression, NSUnknown},
	EventStatusSpamCleared:       {NSComplaint, NSSystem, NSUnknown},
}

func ParseReason(r Reason) (namespace string, code string, ok bool) {
	s := string(r)
	i := strings.IndexByte(s, '.')
	if i <= 0 || i >= len(s)-1 {
		return "", "", false
	}
	return s[:i], s[i+1:], true
}

func ValidateReasonForEventType(et EventType, r Reason) bool {
	ns, _, ok := ParseReason(r)
	if !ok {
		return false
	}
	allowed := AllowedReasonNamespacesByEventType[et]
	for _, a := range allowed {
		if ns == string(a) {
			return true
		}
	}
	return false
}
