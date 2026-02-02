package toolkit

import (
	"context"
	"errors"
	"testing"

	"github.com/commsok/mxevents"
)

// MockEnricher is a test double for mxevents.Enricher
type MockEnricher struct {
	shouldError bool
}

func (m *MockEnricher) Enrich(ctx *context.Context, facts *mxevents.EventFacts) error {
	if m.shouldError {
		return errors.New("enricher error")
	}
	return nil
}

// MockClassifier is a test double for mxevents.Classifier
type MockClassifier struct {
	result      *mxevents.ClassificationResult
	shouldError bool
}

func (m *MockClassifier) Classify(ctx *context.Context, facts *mxevents.EventFacts, taxonomyVersion int) (*mxevents.ClassificationResult, error) {
	if m.shouldError {
		return nil, errors.New("classifier error")
	}
	return m.result, nil
}

func TestEventClassifier_Classify_NilResult(t *testing.T) {
	// Test case: classifier returns nil result (no classification possible)
	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{
			&MockClassifier{result: nil, shouldError: false},
		},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{
		SMTPCode: "550",
	}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
}

func TestEventClassifier_Classify_MultipleClassifiers_SelectsHighestConfidence(t *testing.T) {
	// Test case: multiple classifiers return results, should select highest confidence
	lowConfidenceResult := &mxevents.ClassificationResult{
		TaxonomyVersion: mxevents.CurrentTaxonomyVersion,
		EventType:       mxevents.EventMailboxPermFail,
		Reason:          mxevents.BounceReasonUndefined,
		Confidence:      0.3,
	}

	highConfidenceResult := &mxevents.ClassificationResult{
		TaxonomyVersion: mxevents.CurrentTaxonomyVersion,
		EventType:       mxevents.EventMailboxRecipientPermFail,
		Reason:          mxevents.BounceReasonUserUnknown,
		Confidence:      0.9,
	}

	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{
			&MockClassifier{result: lowConfidenceResult, shouldError: false},
			&MockClassifier{result: highConfidenceResult, shouldError: false},
		},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{
		SMTPCode: "550",
	}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Confidence != 0.9 {
		t.Errorf("expected confidence 0.9, got %f", result.Confidence)
	}

	if result.Reason != mxevents.BounceReasonUserUnknown {
		t.Errorf("expected reason %s, got %s", mxevents.BounceReasonUserUnknown, result.Reason)
	}
}

func TestEventClassifier_Classify_MixedNilAndNonNilResults(t *testing.T) {
	// Test case: some classifiers return nil, others return results
	// Should handle nil results gracefully and select the best non-nil result
	validResult := &mxevents.ClassificationResult{
		TaxonomyVersion: mxevents.CurrentTaxonomyVersion,
		EventType:       mxevents.EventMailboxPermFail,
		Reason:          mxevents.BounceReasonUserUnknown,
		Confidence:      0.8,
	}

	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{
			&MockClassifier{result: nil, shouldError: false},
			&MockClassifier{result: validResult, shouldError: false},
			&MockClassifier{result: nil, shouldError: false},
		},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{
		SMTPCode: "550",
	}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Confidence != 0.8 {
		t.Errorf("expected confidence 0.8, got %f", result.Confidence)
	}
}

func TestEventClassifier_Classify_EnricherError(t *testing.T) {
	// Test case: enricher returns error
	classifier := NewEventClassifier(
		[]mxevents.Enricher{
			&MockEnricher{shouldError: true},
		},
		[]mxevents.Classifier{},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err == nil {
		t.Fatal("expected error from enricher")
	}

	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}

func TestEventClassifier_Classify_ClassifierError(t *testing.T) {
	// Test case: classifier returns error
	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{
			&MockClassifier{shouldError: true},
		},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err == nil {
		t.Fatal("expected error from classifier")
	}

	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}

func TestEventClassifier_Classify_AllClassifiersReturnNil(t *testing.T) {
	// Test case: all classifiers return nil results
	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{
			&MockClassifier{result: nil, shouldError: false},
			&MockClassifier{result: nil, shouldError: false},
		},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result when all classifiers return nil, got %+v", result)
	}
}

func TestEventClassifier_Classify_NoClassifiers(t *testing.T) {
	// Test case: no classifiers configured
	classifier := NewEventClassifier(
		[]mxevents.Enricher{},
		[]mxevents.Classifier{},
	)

	ctx := context.Background()
	facts := &mxevents.EventFacts{}

	result, err := classifier.Classify(&ctx, facts, mxevents.CurrentTaxonomyVersion)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result with no classifiers, got %+v", result)
	}
}
