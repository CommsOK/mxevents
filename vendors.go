package mxevents

// VendorURI constants represent behavior bucket identifiers.
// They look like domains but are NOT the recipient's domain.
// They are canonical keys you normalize to, to hint expected mailbox behavior.
const (
	// VendorURIGmail represents Google-backed mailbox behavior (Gmail + Google Workspace hosted domains).
	VendorURIGmail = "gmail.com"

	// VendorURIOutlook represents Microsoft-backed mailbox behavior (Outlook/Hotmail/Live/MSN + Microsoft 365 hosted domains).
	VendorURIOutlook = "outlook.com"

	// VendorURIYahoo represents Yahoo-backed mailbox behavior (Yahoo + AOL, and related).
	VendorURIYahoo = "yahoo.com"

	// VendorURIiCloud represents Apple-backed mailbox behavior (iCloud + legacy Me/Mac).
	VendorURIiCloud = "icloud.com"

	// VendorURIProton represents Proton-backed mailbox behavior.
	VendorURIProton = "proton.me"

	// VendorURIFastmail represents Fastmail-backed mailbox behavior.
	VendorURIFastmail = "fastmail.com"

	// VendorURIZoho represents Zoho-backed mailbox behavior (Zoho Mail hosted domains).
	VendorURIZoho = "zoho.com"

	// VendorURIYandex represents Yandex-backed mailbox behavior.
	VendorURIYandex = "yandex.com"

	// VendorURIGMX represents GMX / mail.com behavior (United Internet).
	VendorURIGMX = "gmx.com"

	// VendorURITuta represents Tuta (Tutanota) behavior.
	VendorURITuta = "tuta.com"

	// VendorURINaver represents Naver (Korea) behavior.
	VendorURINaver = "naver.com"

	// VendorURIDaum represents Daum / Kakao (Korea) behavior.
	VendorURIDaum = "daum.net"

	// VendorURIMailRu represents Mail.ru (VK) behavior.
	VendorURIMailRu = "mail.ru"

	// VendorURIComcast represents Comcast behavior (big US ISP bucket; optional).
	VendorURIComcast = "comcast.net"

	// VendorURIUnknown is the default bucket when you cannot infer provider behavior.
	VendorURIUnknown = "unknown"
)
