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

func (c *EventClassifier) Classify(ctx *context.Context, facts *mxevents.EventFacts) (*mxevents.ClassificationResult, error) {
	// Enrich facts first
	for _, enricher := range c.enrichers {
		if err := enricher.Enrich(ctx, facts); err != nil {
			return nil, err
		}
	}

	// Get all classifications and return the one with the highest confidence
	var bestResult *mxevents.ClassificationResult
	for _, classifier := range c.classifiers {
		result, err := classifier.Classify(ctx, facts)
		if err != nil {
			return nil, err
		}
		if bestResult == nil || result.Confidence > bestResult.Confidence {
			bestResult = result
		}
	}

	return bestResult, nil
}
