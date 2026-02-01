package enrichment

import (
	"context"

	"github.com/commsok/mxevents"
)

type SMTPEnricher struct {
}

func (e *SMTPEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
	return nil
}
