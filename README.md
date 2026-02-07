# mxevents

A canonical taxonomy and toolkit for email event classification in Go.

## Overview

mxevents provides a standardized schema for classifying email delivery events across the entire email pipeline—from origin platforms (CRMs) through gateways (ESPs) to mailbox delivery and recipient engagement. It enables consistent event handling regardless of the underlying email service provider.

## Features

- **Canonical Event Types**: Standardized event taxonomy covering origin, gateway, mailbox, engagement, and status events
- **Detailed Reason Codes**: Granular failure reasons with sender/recipient attribution
- **Pluggable Architecture**: Extensible enrichers and classifiers
- **Bounce Classification**: Built-in bounce classification using [Sisimai](https://libsisimai.org/)
- **Versioned Taxonomy**: Semantic versioning for backwards compatibility

## Installation

```bash
go get github.com/commsok/mxevents
```

Requires Go 1.24 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/commsok/mxevents"
    "github.com/commsok/mxevents/toolkit"
)

func main() {
    // Create a classifier with default enrichers and classifiers
    classifier := toolkit.NewDefaultEventClassifier()

    // Create event facts from your email provider's webhook
    facts := &mxevents.EventFacts{
        SMTPResponse:       "550 5.1.1 User unknown",
        SMTPCode:           "550",
        SMTPDeliveryStatus: "5.1.1",
        Sender: mxevents.SenderFacts{
            Vendor:    "sendgrid",
            EventName: "bounce",
        },
        Recipient: mxevents.RecipientFacts{
            RecipientDomain: "example.com",
        },
    }

    // Classify the event
    ctx := context.Background()
    result, err := classifier.Classify(&ctx, facts, 0) // 0 = latest taxonomy version
    if err != nil {
        panic(err)
    }

    fmt.Printf("Event Type: %s\n", result.EventType)
    fmt.Printf("Reason: %s\n", result.Reason)
    fmt.Printf("Confidence: %.2f\n", result.Confidence)
}
```

## Event Taxonomy

Events are organized into five categories representing stages of the email delivery pipeline:

### Origin Events (CRM/Vendor Platform)
| Event | Description |
|-------|-------------|
| `origin-success` | Message accepted into vendor's pipeline |
| `origin-failed` | Message rejected before reaching ESP |
| `origin-dropped` | Message intentionally discarded (suppression, policy) |

### Gateway Events (ESP Layer)
| Event | Description |
|-------|-------------|
| `gateway-accepted` | ESP acknowledged and queued the message |
| `gateway-success` | Message prepared for SMTP delivery |
| `gateway-failed` | Gateway could not process the message |
| `gateway-dropped` | Gateway chose not to attempt delivery |

### Mailbox Events (Recipient Mail Server)
| Event | Description |
|-------|-------------|
| `mailbox-attempt` | SMTP delivery attempted |
| `mailbox-success` | Message accepted (2xx response) |
| `mailbox-failed` | Delivery failed (unknown cause) |
| `mailbox-tempfail` | Temporary failure (4xx), retry may succeed |
| `mailbox-permfail` | Permanent failure (5xx), unknown attribution |
| `mailbox-sender-permfail` | Permanent failure, sender-side cause |
| `mailbox-recipient-permfail` | Permanent failure, recipient-side cause |
| `mailbox-quarantined` | Message placed in spam/quarantine |

### Engagement Events (Recipient Actions)
| Event | Description |
|-------|-------------|
| `engagement-open` | Email opened (tracking pixel fired) |
| `engagement-click` | Link clicked |
| `engagement-engaged` | Custom engagement signal |

### Status Events (Preferences/Compliance)
| Event | Description |
|-------|-------------|
| `status-subscribed` | Recipient globally subscribed |
| `status-unsubscribed` | Recipient opted out globally |
| `status-group-subscribed` | Recipient opted into a topic/list |
| `status-group-unsubscribed` | Recipient opted out of a topic/list |
| `status-spam-reported` | Spam complaint received |
| `status-spam-cleared` | Spam complaint status cleared |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    EventClassifier                       │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │  Enricher   │ -> │  Enricher   │ -> │  Enricher   │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│         │                                               │
│         v                                               │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │ Classifier  │    │ Classifier  │    │ Classifier  │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│         │                 │                  │          │
│         └────────────────┬──────────────────┘          │
│                          v                              │
│              Best Classification Result                 │
└─────────────────────────────────────────────────────────┘
```

### Enrichers

Enrichers augment `EventFacts` with additional context before classification:

- **SMTPEnricher**: Parses SMTP response codes and delivery status
- **CommonVendorURIEnricher**: Extracts domain information from vendor URIs

### Classifiers

Classifiers analyze enriched facts and return classification results:

- **SisimaiClassifier**: Uses Sisimai library for bounce classification

## Extending

### Custom Enricher

```go
type MyEnricher struct{}

func (e *MyEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
    // Add custom enrichment logic
    return nil
}
```

### Custom Classifier

```go
type MyClassifier struct{}

func (c *MyClassifier) Classify(ctx *context.Context, facts *mxevents.EventFacts, version int) (*mxevents.ClassificationResult, error) {
    // Add custom classification logic
    return &mxevents.ClassificationResult{
        TaxonomyVersion: version,
        EventType:       mxevents.EventMailboxPermFail,
        Reason:          mxevents.ReasonUnknown,
        Confidence:      0.8,
        Facts:           facts,
    }, nil
}
```

### Using Custom Components

```go
classifier := toolkit.NewEventClassifier(
    append(
        []mxevents.Enricher{&MyEnricher{}},
        toolkit.DefaultEnrichers...,
    ),
    append(
        []mxevents.Classifier{&MyClassifier{}},
        toolkit.DefaultClassifiers...,
    ),
)
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
