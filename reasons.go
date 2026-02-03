package mxevents

import "strings"

// SMTP / mailbox-side classification codes.
// Use when the failure is attributable to the remote mail system response
// (or the SMTP session outcome) rather than your local platform decision.
const (
	ReasonSMTPUserUnknown     Reason = "smtp.UserUnknown"     // 5.1.1, no such user, mailbox not found
	ReasonSMTPMailboxFull     Reason = "smtp.MailboxFull"     // 4.2.2, over quota
	ReasonSMTPHasMoved        Reason = "smtp.HasMoved"        // address moved/forward-only
	ReasonSMTPSpamDetected    Reason = "smtp.SpamDetected"    // message flagged as spam, bulk, content scoring
	ReasonSMTPPolicyViolation Reason = "smtp.PolicyViolation" // remote policy rejection, not clearly spam
	ReasonSMTPAuthFailure     Reason = "smtp.AuthFailure"     // SPF/DKIM/DMARC/auth related rejection
	ReasonSMTPBlocked         Reason = "smtp.Blocked"         // blocked/denied/refused, scope unclear
	ReasonSMTPRateLimited     Reason = "smtp.RateLimited"     // 4xx throttle, too many connections, deferrals
	ReasonSMTPMessageTooLarge Reason = "smtp.MessageTooLarge" // size limit exceeded
	ReasonSMTPVirusDetected   Reason = "smtp.VirusDetected"   // malware/virus rejection
	ReasonSMTPSyntaxError     Reason = "smtp.SyntaxError"     // bad address syntax, protocol syntax issues
	ReasonSMTPNoRelaying      Reason = "smtp.NoRelaying"      // relaying denied
	ReasonSMTPExpired         Reason = "smtp.Expired"         // retried then abandoned (final failure)
)

const (
	ReasonNetworkDnsFailure     Reason = "network.DnsFailure"     // NXDOMAIN, SERVFAIL, host unknown, no MX
	ReasonNetworkConnectFailure Reason = "network.ConnectFailure" // connection refused, reset, unreachable
	ReasonNetworkTimeout        Reason = "network.Timeout"        // dial timeout, read timeout
	ReasonNetworkTlsFailure     Reason = "network.TlsFailure"     // STARTTLS/handshake/cert failures
	ReasonNetworkRoutingError   Reason = "network.RoutingError"   // network unreachable, routing issues
)

const (
	ReasonSuppressionSuppressed    Reason = "suppression.Suppressed"    // already suppressed (list/decision)
	ReasonSuppressionUnsubscribed  Reason = "suppression.Unsubscribed"  // recipient opted out
	ReasonSuppressionSpamComplaint Reason = "suppression.SpamComplaint" // prior complaint enforcement
	ReasonSuppressionPolicyBlocked Reason = "suppression.PolicyBlocked" // local policy prevents send
	ReasonSuppressionAdminBlocked  Reason = "suppression.AdminBlocked"  // manual block/suppression
)

const (
	ReasonPolicyBlocked       Reason = "policy.Blocked"       // blocked by policy, deny list, enforcement
	ReasonPolicyBadReputation Reason = "policy.BadReputation" // reputation-based rejection
	ReasonPolicyRestricted    Reason = "policy.Restricted"    // not allowed recipient/sender/content category
)

const (
	ReasonConfigMissing    Reason = "config.Missing"    // missing config (domain, sender, template, settings)
	ReasonConfigInvalid    Reason = "config.Invalid"    // invalid config values
	ReasonAuthUnauthorized Reason = "auth.Unauthorized" // authz failure, revoked, insufficient scopes
	ReasonAuthInvalid      Reason = "auth.Invalid"      // bad credentials, invalid token, signature mismatch
)

const (
	ReasonSystemError        Reason = "system.Error"        // generic internal error
	ReasonSystemUnavailable  Reason = "system.Unavailable"  // upstream down, maintenance
	ReasonSystemRateLimited  Reason = "system.RateLimited"  // local queue/concurrency limit
	ReasonSystemTimeout      Reason = "system.Timeout"      // internal or upstream timeout
	ReasonSystemBackpressure Reason = "system.Backpressure" // queue full, resource pressure
)

const (
	ReasonComplaintFeedbackLoop      Reason = "complaint.FeedbackLoop"      // ISP FBL
	ReasonComplaintUserAction        Reason = "complaint.UserAction"        // recipient marked as spam
	ReasonComplaintProviderHeuristic Reason = "complaint.ProviderHeuristic" // inferred by provider/platform

	ReasonUnsubscribeOneClick      Reason = "unsubscribe.OneClick"      // one-click unsubscribe
	ReasonUnsubscribeUserInitiated Reason = "unsubscribe.UserInitiated" // link, reply, UI
	ReasonUnsubscribeAdminAction   Reason = "unsubscribe.AdminAction"   // admin/operator action
	ReasonUnsubscribeCompliance    Reason = "unsubscribe.Compliance"    // legal/regulatory enforcement
)

const (
	ReasonUnknown Reason = "unknown.Unknown"
)

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
	// Origin
	EventOriginFailed:  {NSConfig, NSAuth, NSSystem, NSPolicy, NSUnknown},
	EventOriginDropped: {NSSuppression, NSPolicy, NSConfig, NSSystem, NSUnknown},

	// Gateway
	EventGatewayFailed:  {NSConfig, NSAuth, NSSystem, NSPolicy, NSUnknown},
	EventGatewayDropped: {NSSuppression, NSPolicy, NSSystem, NSUnknown},

	// Mailbox
	EventMailboxFailed:            {NSNetwork, NSSystem, NSUnknown},
	EventMailboxTempFail:          {NSMTP, NSNetwork, NSSystem, NSPolicy, NSUnknown},
	EventMailboxPermFail:          {NSMTP, NSPolicy, NSUnknown},
	EventMailboxSenderPermFail:    {NSMTP, NSPolicy, NSAuth, NSConfig, NSUnknown},
	EventMailboxRecipientPermFail: {NSMTP, NSUnknown},

	// Status / preference
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
