package config

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsEmpty(t *testing.T) {

	predicate := (&PredicateConfig{Field: "Message", Is: "Empty"}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestIsNotEmpty(t *testing.T) {

	predicate := (&PredicateConfig{Field: "Message", Is: "Not Empty"}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}
