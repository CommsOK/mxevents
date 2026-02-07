package enrichment

import (
	"context"
	"strings"

	"github.com/commsok/mxevents"
)

// CommonMailboxVendorURIEnricher is a simple reference enricher that sets the MailboxVendorURI field based on the recipient domain.
// It is not intended to be a complete solution, but rather a starting point for custom enrichers.
type CommonMailboxVendorURIEnricher struct{}

func (e *CommonMailboxVendorURIEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
	if facts.Recipient.MailboxVendorURI != "" {
		return nil
	}

	if facts.Recipient.RecipientDomain == "" {
		return nil
	}

	// Normalize domain to lowercase for case-insensitive matching
	domain := strings.ToLower(facts.Recipient.RecipientDomain)
	if vendorURI, ok := CommonDomainToMailboxVendorURI[domain]; ok {
		facts.Recipient.MailboxVendorURI = vendorURI
	}

	return nil
}

// CommonDomainToMailboxVendorURI maps common recipient domains to their corresponding mailbox vendor URI buckets.
// This mapping is based on domain name alone and includes only the "obvious" mappings.
var CommonDomainToMailboxVendorURI = map[string]mxevents.MailboxVendorURI{
	// Google
	"gmail.com":      mxevents.MailboxVendorURIGmail,
	"googlemail.com": mxevents.MailboxVendorURIGmail,

	// Microsoft
	"outlook.com": mxevents.MailboxVendorURIOutlook,
	"hotmail.com": mxevents.MailboxVendorURIOutlook,
	"live.com":    mxevents.MailboxVendorURIOutlook,
	"msn.com":     mxevents.MailboxVendorURIOutlook,

	// Yahoo / AOL
	"yahoo.com":      mxevents.MailboxVendorURIYahoo,
	"ymail.com":      mxevents.MailboxVendorURIYahoo,
	"rocketmail.com": mxevents.MailboxVendorURIYahoo,
	"aol.com":        mxevents.MailboxVendorURIYahoo,

	// Apple
	"icloud.com": mxevents.MailboxVendorURIiCloud,
	"me.com":     mxevents.MailboxVendorURIiCloud,
	"mac.com":    mxevents.MailboxVendorURIiCloud,

	// Proton
	"proton.me":      mxevents.MailboxVendorURIProton,
	"protonmail.com": mxevents.MailboxVendorURIProton,
	"pm.me":          mxevents.MailboxVendorURIProton,

	// Fastmail
	"fastmail.com": mxevents.MailboxVendorURIFastmail,

	// Zoho
	"zoho.com":     mxevents.MailboxVendorURIZoho,
	"zohomail.com": mxevents.MailboxVendorURIZoho,

	// Yandex
	"yandex.com": mxevents.MailboxVendorURIYandex,
	"yandex.ru":  mxevents.MailboxVendorURIYandex,
	"ya.ru":      mxevents.MailboxVendorURIYandex,

	// GMX / Mail.com
	"gmx.com":  mxevents.MailboxVendorURIGMX,
	"gmx.de":   mxevents.MailboxVendorURIGMX,
	"gmx.net":  mxevents.MailboxVendorURIGMX,
	"mail.com": mxevents.MailboxVendorURIGMX,

	// Tuta (Tutanota)
	"tuta.com":     mxevents.MailboxVendorURITuta,
	"tutanota.com": mxevents.MailboxVendorURITuta,
	"tutanota.de":  mxevents.MailboxVendorURITuta,

	// Korea
	"naver.com":   mxevents.MailboxVendorURINaver,
	"daum.net":    mxevents.MailboxVendorURIDaum,
	"hanmail.net": mxevents.MailboxVendorURIDaum,

	// Mail.ru (VK)
	"mail.ru":  mxevents.MailboxVendorURIMailRu,
	"inbox.ru": mxevents.MailboxVendorURIMailRu,
	"list.ru":  mxevents.MailboxVendorURIMailRu,
	"bk.ru":    mxevents.MailboxVendorURIMailRu,

	// ISP example (optional)
	"comcast.net": mxevents.MailboxVendorURIComcast,
}
