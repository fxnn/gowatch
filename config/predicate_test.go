package config

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestZeroPredicate(t *testing.T) {

	predicate := (&PredicateConfig{}).CreatePredicate()

	require.True(t, predicate.Applies(logentry.New()))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestIsEmpty(t *testing.T) {

	predicate := (&PredicateConfig{"message": map[string]interface{}{"is": "Empty"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestIsNotEmpty(t *testing.T) {

	predicate := (&PredicateConfig{"message": map[string]interface{}{"is": "Not Empty"}}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestCustomIsEmpty(t *testing.T) {

	predicate := (&PredicateConfig{"custom": map[string]interface{}{"is": "Empty"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Custom: map[string]string{"custom": "that's not empty"}}))

}

func TestCustomIsNotEmpty(t *testing.T) {

	predicate := (&PredicateConfig{"custom": map[string]interface{}{"is": "Not Empty"}}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Custom: map[string]string{"custom": "that's not empty"}}))

}

func TestNot(t *testing.T) {

	predicate := (&PredicateConfig{"not": messageIsEmpty()}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "that's not empty"}))

}

func TestContains(t *testing.T) {

	predicate := (&PredicateConfig{"message": PredicateConfig{"contains": " does "}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "this does contain our substring"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "this doesn't contain our substring"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{}))

}

func TestTagsContains(t *testing.T) {

	predicate := (&PredicateConfig{"tags": PredicateConfig{"contains": "mytag"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Tags: []string{"mytag", "another"}}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Tags: []string{"another", "yet another"}}))
	require.False(t, predicate.Applies(&logentry.LogEntry{}))

}

func TestMatches(t *testing.T) {

	predicate := (&PredicateConfig{"host": PredicateConfig{"matches": "%{IPV4}"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Host: "127.0.0.1"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Host: "localhost"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "127.0.0.1"}))

}

func TestAfter(t *testing.T) {

	predicate := (&PredicateConfig{"timestamp": PredicateConfig{"after": "2015-01-01T00:00:00Z"}}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("1970-01-01T00:00:00Z")}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2014-12-31T23:59:59Z")}))

	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2015-01-01T00:00:00Z")}))

	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2015-01-01T00:00:01Z")}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2031-01-01T00:00:00Z")}))

}

func TestBefore(t *testing.T) {

	predicate := (&PredicateConfig{"timestamp": PredicateConfig{"before": "2015-01-01T00:00:00Z"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("1970-01-01T00:00:00Z")}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2014-12-31T23:59:59Z")}))

	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2015-01-01T00:00:00Z")}))

	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2015-01-01T00:00:01Z")}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: timeFromRfc3339("2031-01-01T00:00:00Z")}))

}

func timeFromRfc3339(rfc3339FormattedTime string) time.Time {
	result, err := time.Parse(time.RFC3339, rfc3339FormattedTime)
	if err == nil {
		return result
	}
	log.Panic(err)
	return time.Time{} // actually never reached
}

func TestAllOf(t *testing.T) {

	predicate := (&PredicateConfig{
		"allof": []PredicateConfig{
			PredicateConfig{"message": PredicateConfig{"contains": "A"}},
			PredicateConfig{"host": PredicateConfig{"contains": "A"}},
		},
	}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestAllOf_DirectMap(t *testing.T) {

	predicate := (&PredicateConfig{
		"allof": PredicateConfig{
			"message": PredicateConfig{"contains": "A"},
			"host":    PredicateConfig{"contains": "A"},
		},
	}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestAllOf_Implicit(t *testing.T) {

	predicate := (&PredicateConfig{
		"message": PredicateConfig{"contains": "A"},
		"host":    PredicateConfig{"contains": "A"},
	}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestAnyOf(t *testing.T) {

	predicate := (&PredicateConfig{
		"anyof": []PredicateConfig{
			PredicateConfig{"message": PredicateConfig{"contains": "A"}},
			PredicateConfig{"host": PredicateConfig{"contains": "A"}},
		},
	}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestAnyOf_DirectMap(t *testing.T) {

	predicate := (&PredicateConfig{
		"anyof": PredicateConfig{
			"message": PredicateConfig{"contains": "A"},
			"host":    PredicateConfig{"contains": "A"},
		},
	}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestNoneOf(t *testing.T) {

	predicate := (&PredicateConfig{
		"noneof": []PredicateConfig{
			PredicateConfig{"message": PredicateConfig{"contains": "A"}},
			PredicateConfig{"host": PredicateConfig{"contains": "A"}},
		},
	}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestNoneOf_DirectMap(t *testing.T) {

	predicate := (&PredicateConfig{
		"noneof": PredicateConfig{
			"message": PredicateConfig{"contains": "A"},
			"host":    PredicateConfig{"contains": "A"},
		},
	}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xAx"}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Message: "xAx", Host: "xxx"}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Message: "xxx", Host: "xxx"}))

}

func TestYoungerThan(t *testing.T) {

	now := time.Now()
	predicate := (&PredicateConfig{"timestamp": PredicateConfig{"younger than": "24h"}}).CreatePredicate()

	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -2)}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -1)}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -1).Add(time.Microsecond)}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: now}))

}

func TestOlderThan(t *testing.T) {

	now := time.Now()
	predicate := (&PredicateConfig{"timestamp": PredicateConfig{"older than": "24h"}}).CreatePredicate()

	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -2)}))
	require.True(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -1).Add(-1 * time.Microsecond)}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: now.AddDate(0, 0, -1)}))
	require.False(t, predicate.Applies(&logentry.LogEntry{Timestamp: now}))

}

func messageIsEmpty() PredicateConfig {
	return PredicateConfig{"message": PredicateConfig{"is": "Empty"}}
}
