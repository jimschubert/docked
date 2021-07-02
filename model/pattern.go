package model

import (
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
)

type RegexEngine int

const (
	RegexpEngine RegexEngine = iota
	Regexp2Engine
)

type regexEngine interface {
	MatchString(pattern string, value string) bool
}

type goRegexpEngine struct{}

func (g goRegexpEngine) MatchString(pattern string, value string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

type regexp2Engine struct{}

func (r regexp2Engine) MatchString(pattern string, value string) bool {
	re := regexp2.MustCompile(pattern, 0)
	re.MatchTimeout = time.Second * 1
	isMatch, _ := re.MatchString(value)
	return isMatch
}

var (
	engine regexEngine = goRegexpEngine{}
)

type Pattern interface {
	Matches(value string) bool
}

type StandardPattern struct {
	pattern string
}

func (s StandardPattern) Matches(value string) bool {
	return engine.MatchString(s.pattern, value)
}

func NewPattern(value string) Pattern {
	return StandardPattern{
		pattern: value,
	}
}

func SetRegexEngine(e RegexEngine) {
	if e == RegexpEngine {
		engine = goRegexpEngine{}
	} else if e == Regexp2Engine {
		engine = regexp2Engine{}
	} else {
		panic("unsupported regex engine")
	}
}
