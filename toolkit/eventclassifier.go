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
	&bounce.SisimaiClassifier{},
}

func NewEventClassifier(enrichers []mxevents.Enricher, classifiers []mxevents.Classifier) *EventClassifier {
	return &EventClassifier{enrichers: enrichers, classifiers: classifiers}
}

func NewDefaultEventClassifier() *EventClassifier {
	return &EventClassifier{enrichers: DefaultEnrichers, classifiers: DefaultClassifiers}
}

// Classify classifies the event and returns event details according to the taxonomy version requested.
// Parameters:
//   - ctx: context for cancellation and tracing
//   - facts: event facts to classify
//   - taxonomyVersion: version of the taxonomy to use for classification. If 0, uses the latest version.
func (c *EventClassifier) Classify(ctx *context.Context, facts *mxevents.EventFacts, taxonomyVersion int) (*mxevents.ClassificationResult, error) {
	// Enrich facts first
	for _, enricher := range c.enrichers {
		if err := enricher.Enrich(ctx, facts); err != nil {
			return nil, err
		}
	}

	// Get all classifications and return the one with the highest confidence
	var bestResult *mxevents.ClassificationResult
	for _, classifier := range c.classifiers {
		result, err := classifier.Classify(ctx, facts, taxonomyVersion)
		if err != nil {
			return nil, err
		}
		if result != nil && (bestResult == nil || result.Confidence > bestResult.Confidence) {
			bestResult = result
		}
	}

	return bestResult, nil
}
