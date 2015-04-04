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

func TestNot(t *testing.T) {

	predicate := (&PredicateConfig{Not: messageIsEmpty()}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestContains(t *testing.T) {

	predicate := (&PredicateConfig{Field: "Message", Contains: " does "}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "this does contain our substring"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "this doesn't contain our substring"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{}))

}

func TestMatches(t *testing.T) {

	predicate := (&PredicateConfig{Field: "Host", Matches: "%{IPV4}"}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Host: "127.0.0.1"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Host: "localhost"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "127.0.0.1"}))

}

func messageIsEmpty() *PredicateConfig {
	return &PredicateConfig{Field: "Message", Is: "Empty"}
}
