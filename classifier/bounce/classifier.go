package bounce

import (
	"context"

	"github.com/commsok/mxevents"
)

type Classifier struct {
}

func (c *Classifier) Classify(ctx *context.Context, facts *mxevents.EventFacts) (*mxevents.ClassificationResult, error) {
	return nil, nil
}
