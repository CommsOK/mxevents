package mxevents

import "context"

type Classifier interface {
	Classify(ctx *context.Context, facts *EventFacts) (*ClassificationResult, error)
}
