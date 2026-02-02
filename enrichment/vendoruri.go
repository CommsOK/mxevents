package enrichment

import (
	"context"
	"strings"

	"github.com/commsok/mxevents"
)

type CommonVendorURIEnricher struct {
}

func (e *CommonVendorURIEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
	if facts.Recipient.VendorURI != "" {
		return nil
	}

	if facts.Recipient.RecipientDomain == "" {
		return nil
	}

	// Normalize domain to lowercase for case-insensitive matching
	domain := strings.ToLower(facts.Recipient.RecipientDomain)
	if vendorURI, ok := CommonDomainToVendorURI[domain]; ok {
		facts.Recipient.VendorURI = vendorURI
	}

	return nil
}

// CommonDomainToVendorURI maps common recipient domains to their corresponding VendorURI buckets.
// This mapping is based on domain name alone and includes only the "obvious" mappings.
var CommonDomainToVendorURI = map[string]string{
	// Google
	"gmail.com":      mxevents.VendorURIGmail,
	"googlemail.com": mxevents.VendorURIGmail,

	// Microsoft
	"outlook.com": mxevents.VendorURIOutlook,
	"hotmail.com": mxevents.VendorURIOutlook,
	"live.com":    mxevents.VendorURIOutlook,
	"msn.com":     mxevents.VendorURIOutlook,

	// Yahoo / AOL
	"yahoo.com":      mxevents.VendorURIYahoo,
	"ymail.com":      mxevents.VendorURIYahoo,
	"rocketmail.com": mxevents.VendorURIYahoo,
	"aol.com":        mxevents.VendorURIYahoo,

	// Apple
	"icloud.com": mxevents.VendorURIiCloud,
	"me.com":     mxevents.VendorURIiCloud,
	"mac.com":    mxevents.VendorURIiCloud,

	// Proton
	"proton.me":      mxevents.VendorURIProton,
	"protonmail.com": mxevents.VendorURIProton,
	"pm.me":          mxevents.VendorURIProton,

	// Fastmail
	"fastmail.com": mxevents.VendorURIFastmail,

	// Zoho
	"zoho.com":     mxevents.VendorURIZoho,
	"zohomail.com": mxevents.VendorURIZoho,

	// Yandex
	"yandex.com": mxevents.VendorURIYandex,
	"yandex.ru":  mxevents.VendorURIYandex,
	"ya.ru":      mxevents.VendorURIYandex,

	// GMX / Mail.com
	"gmx.com":  mxevents.VendorURIGMX,
	"gmx.de":   mxevents.VendorURIGMX,
	"gmx.net":  mxevents.VendorURIGMX,
	"mail.com": mxevents.VendorURIGMX,

	// Tuta (Tutanota)
	"tuta.com":     mxevents.VendorURITuta,
	"tutanota.com": mxevents.VendorURITuta,
	"tutanota.de":  mxevents.VendorURITuta,

	// Korea
	"naver.com":   mxevents.VendorURINaver,
	"daum.net":    mxevents.VendorURIDaum,
	"hanmail.net": mxevents.VendorURIDaum,

	// Mail.ru (VK)
	"mail.ru":  mxevents.VendorURIMailRu,
	"inbox.ru": mxevents.VendorURIMailRu,
	"list.ru":  mxevents.VendorURIMailRu,
	"bk.ru":    mxevents.VendorURIMailRu,

	// ISP example (optional)
	"comcast.net": mxevents.VendorURIComcast,
}
