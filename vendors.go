package mxevents

// This file contains two distinct vendor concepts:
//  1) SourceVendor: the system that generated the event (ESP webhook, CRM, etc.)
//  2) MailboxVendorURI: the recipient mailbox provider behavior bucket (Gmail, Outlook, ...)
//
// This package intentionally distinguishes between the event source vendor and
// the recipient mailbox vendor.

// SourceVendor is the vendor/platform that generated the event.
type SourceVendor string

// MailboxVendorURI is a canonical string key for mailbox provider behavior.
// Values look like domains but are NOT the recipient's domain.
type MailboxVendorURI string

// SourceVendor* constants identify the event source (ESP/CRM/etc.).
// This list is not exhaustive; callers may use custom strings.
const (
	// SourceVendorUnknown is the fallback when you cannot identify the event source.
	SourceVendorUnknown SourceVendor = "unknown"

	// SourceVendorSendGrid identifies SendGrid as the event source.
	SourceVendorSendGrid SourceVendor = "sendgrid"

	// SourceVendorHubSpot identifies HubSpot as the event source.
	SourceVendorHubSpot SourceVendor = "hubspot"
)

// MailboxVendorURI* constants represent mailbox provider behavior buckets.
const (
	// MailboxVendorURIGmail represents Google-backed mailbox behavior (Gmail + Google Workspace hosted domains).
	MailboxVendorURIGmail MailboxVendorURI = "gmail.com"

	// MailboxVendorURIOutlook represents Microsoft-backed mailbox behavior (Outlook/Hotmail/Live/MSN + Microsoft 365 hosted domains).
	MailboxVendorURIOutlook MailboxVendorURI = "outlook.com"

	// MailboxVendorURIYahoo represents Yahoo-backed mailbox behavior (Yahoo + AOL, and related).
	MailboxVendorURIYahoo MailboxVendorURI = "yahoo.com"

	// MailboxVendorURIiCloud represents Apple-backed mailbox behavior (iCloud + legacy Me/Mac).
	MailboxVendorURIiCloud MailboxVendorURI = "icloud.com"

	// MailboxVendorURIProton represents Proton-backed mailbox behavior.
	MailboxVendorURIProton MailboxVendorURI = "proton.me"

	// MailboxVendorURIFastmail represents Fastmail-backed mailbox behavior.
	MailboxVendorURIFastmail MailboxVendorURI = "fastmail.com"

	// MailboxVendorURIZoho represents Zoho-backed mailbox behavior (Zoho Mail hosted domains).
	MailboxVendorURIZoho MailboxVendorURI = "zoho.com"

	// MailboxVendorURIYandex represents Yandex-backed mailbox behavior.
	MailboxVendorURIYandex MailboxVendorURI = "yandex.com"

	// MailboxVendorURIGMX represents GMX / mail.com behavior (United Internet).
	MailboxVendorURIGMX MailboxVendorURI = "gmx.com"

	// MailboxVendorURITuta represents Tuta (Tutanota) behavior.
	MailboxVendorURITuta MailboxVendorURI = "tuta.com"

	// MailboxVendorURINaver represents Naver (Korea) behavior.
	MailboxVendorURINaver MailboxVendorURI = "naver.com"

	// MailboxVendorURIDaum represents Daum / Kakao (Korea) behavior.
	MailboxVendorURIDaum MailboxVendorURI = "daum.net"

	// MailboxVendorURIMailRu represents Mail.ru (VK) behavior.
	MailboxVendorURIMailRu MailboxVendorURI = "mail.ru"

	// MailboxVendorURIComcast represents Comcast behavior (big US ISP bucket; optional).
	MailboxVendorURIComcast MailboxVendorURI = "comcast.net"

	// MailboxVendorURIUnknown is the default bucket when you cannot infer provider behavior.
	MailboxVendorURIUnknown MailboxVendorURI = "unknown"
)
