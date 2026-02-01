package toolkit

import (
	"context"

	"github.com/commsok/mxevents"
	"github.com/commsok/mxevents/classifier/bounce"
	"github.com/commsok/mxevents/enrichment"
)

type EventClassifier struct {
	enrichers   []mxevents.Enricher
	classifiers []mxevents.Classifier
}

var DefaultEnrichers = []mxevents.Enricher{
	&enrichment.SMTPEnricher{},
}

var DefaultClassifiers = []mxevents.Classifier{
	&bounce.Classifier{},
}

func NewEventClassifier(ctx *context.Context, enrichers []mxevents.Enricher, classifiers []mxevents.Classifier) *EventClassifier {
	return &EventClassifier{enrichers: enrichers, classifiers: classifiers}
}

func NewDefaultEventClassifier(ctx *context.Context) *EventClassifier {
	return &EventClassifier{enrichers: DefaultEnrichers, classifiers: DefaultClassifiers}
}
