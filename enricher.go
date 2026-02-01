package mxevents

import "context"

type Enricher interface {
	Enrich(ctx *context.Context, facts *EventFacts) error
}
