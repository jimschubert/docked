package model

import (
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
)

// RegexEngine provides an enum type of supported engines
type RegexEngine int

const (
	// RegexpEngine is the go standard library RegEx engine
	RegexpEngine RegexEngine = iota
	// Regexp2Engine is the regexp2 engine, which is a port of the .NET Core RegEx engine (or, full RegEx support)
	Regexp2Engine
)

type regexEngine interface {
	MatchString(pattern string, value string) bool
}

type goRegexpEngine struct{}

// MatchString determines if value matches against pattern using the go regexp engine.
func (g goRegexpEngine) MatchString(pattern string, value string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

type regexp2Engine struct{}

// MatchString determines if value matches against pattern using the regexp2 engine.
func (r regexp2Engine) MatchString(pattern, value string) bool {
	re := regexp2.MustCompile(pattern, 0)
	re.MatchTimeout = time.Second * 1
	isMatch, _ := re.MatchString(value)
	return isMatch
}

var (
	engine regexEngine = goRegexpEngine{}
)

// Pattern defines the limited string matching interface to perform against a regex pattern
type Pattern interface {
	Matches(value string) bool
}

// StandardPattern is the default implementation abstracting over a pattern for use with the configurable regex engine.
// Use SetRegexEngine to change the underlying regex library if necessary.
type StandardPattern struct {
	pattern string
}

// Matches dispatches string matching to the configured regex engine
func (s StandardPattern) Matches(value string) bool {
	return engine.MatchString(s.pattern, value)
}

// NewPattern creates a new Pattern value from the provided value
// Note that this makes no assertions of the value, and go 'regexp' supports a subset of regex syntax.
// Use SetRegexEngine to change the underlying regex library if necessary.
func NewPattern(value string) Pattern {
	return StandardPattern{
		pattern: value,
	}
}

// SetRegexEngine globally applies the preferred regex library to use
func SetRegexEngine(e RegexEngine) {
	if e == RegexpEngine {
		engine = goRegexpEngine{}
	} else if e == Regexp2Engine {
		engine = regexp2Engine{}
	} else {
		panic("unsupported regex engine")
	}
}
