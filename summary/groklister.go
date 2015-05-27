package summary

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/gemsi/grok"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"log"
	"strings"
	"sync"
)

type GrokLister struct {
	grok                      *grok.Grok
	patternsByListValueByName map[string]map[string]string
	stringsPerPatternName     map[string][]string
	waitGroup                 sync.WaitGroup
	predicate                 logentry.Predicate
	collator                  *collate.Collator
}

func NewGrokLister(patternsByListValueByName map[string]map[string]string, predicate logentry.Predicate) (gl *GrokLister) {
	return NewGrokListerWithLocale("en-US", patternsByListValueByName, predicate)
}

func NewGrokListerWithLocale(locale string, patternsByListValueByName map[string]map[string]string, predicate logentry.Predicate) (gl *GrokLister) {
	return NewGrokListerWithLanguageTag(language.Make(locale), patternsByListValueByName, predicate)
}

func NewGrokListerWithLanguageTag(locale language.Tag, patternsByListValueByName map[string]map[string]string, predicate logentry.Predicate) (gl *GrokLister) {
	gl = new(GrokLister)
	gl.stringsPerPatternName = make(map[string][]string)
	gl.patternsByListValueByName = patternsByListValueByName
	gl.grok = grok.New()
	gl.predicate = predicate
	gl.collator = collate.New(locale, collate.Numeric)

	return
}

func (gl *GrokLister) SummarizeAsync(entries <-chan logentry.LogEntry) {
	gl.waitGroup.Add(1)
	go func() {
		gl.Summarize(entries)
		gl.waitGroup.Done()
	}()
}

func (gl *GrokLister) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		if gl.predicate.Applies(&entry) {
			gl.summarizeEntry(&entry)
		}
	}
}

func (gl *GrokLister) summarizeEntry(entry *logentry.LogEntry) {
	for name, patternsByListValue := range gl.patternsByListValueByName {
		for listValue, pattern := range patternsByListValue {
			matches, err := gl.grok.Parse(pattern, entry.Message)
			if err != nil {
				log.Printf("Error in GrokLister.summarizeEntry: %s [entry=%s]", err, entry)
			} else if len(matches) > 0 {
				for matchName, matchContent := range matches {
					// TODO: Add escape possibility
					name = strings.Replace(name, "%{"+matchName+"}", matchContent, len(name))
					listValue = strings.Replace(listValue, "%{"+matchName+"}", matchContent, len(listValue))
				}
				gl.addStringToPatternName(name, listValue)
			}
		}
	}
}

func (gl *GrokLister) addStringToPatternName(patternName string, newString string) {
	allStrings := gl.stringsPerPatternName[patternName]
	gl.stringsPerPatternName[patternName] = append(allStrings, newString)
}

func (gl *GrokLister) StringAfterSummarizeAsyncCompleted() string {
	gl.waitGroup.Wait()
	return gl.String()
}

func (gl *GrokLister) String() string {
	var result stringList

	for patternName, allStrings := range gl.stringsPerPatternName {
		result = result.Append(patternName + ": " + strings.Join(allStrings, ", "))
	}

	gl.collator.Sort(result)
	return result.Join("\n")
}
