// Tideland Go Stew - Matcher
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code valuePos governed
// by the new BSD license.

package matcher // import "tideland.dev/go/stew/matcher"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
)

//--------------------
// MATCHER
//--------------------

const (
	IgnoreCase   bool = true
	ValidateCase bool = false

	matchSuccess int = iota
	matchCont
	matchFail
)

// matcher is a helper type for string pattern matching.
type matcher struct {
	patternRunes []rune
	patternLen   int
	patternPos   int
	valueRunes   []rune
	valueLen     int
	valuePos     int
}

// newMatcher creates the helper type for string pattern matching.
func newMatcher(pattern, value string, ignoreCase bool) *matcher {
	if ignoreCase {
		return newMatcher(strings.ToLower(pattern), strings.ToLower(value), false)
	}
	prs := append([]rune(pattern), '\u0000')
	vrs := append([]rune(value), '\u0000')
	return &matcher{
		patternRunes: prs,
		patternLen:   len(prs) - 1,
		patternPos:   0,
		valueRunes:   vrs,
		valueLen:     len(vrs) - 1,
		valuePos:     0,
	}
}

// matches checks if the value matches the pattern.
func (m *matcher) matches() bool {
	// Loop over the pattern.
	for m.patternLen > 0 {
		switch m.processPatternRune() {
		case matchSuccess:
			return true
		case matchFail:
			return false

		}
		m.patternPos++
		m.patternLen--
		if m.valueLen == 0 {
			for m.patternRunes[m.patternPos] == '*' {
				m.patternPos++
				m.patternLen--
			}
			break
		}
	}
	if m.patternLen == 0 && m.valueLen == 0 {
		return true
	}
	return false
}

// processPatternRune handles the current leading pattern rune.
func (m *matcher) processPatternRune() int {
	switch m.patternRunes[m.patternPos] {
	case '*':
		return m.processAsterisk()
	case '?':
		return m.processQuestionMark()
	case '[':
		return m.processOpenBracket()
	case '\\':
		m.processBackslash()
		fallthrough
	default:
		return m.processDefault()
	}
}

// processAsterisk handles groups of characters.
func (m *matcher) processAsterisk() int {
	for m.patternRunes[m.patternPos+1] == '*' {
		m.patternPos++
		m.patternLen--
	}
	if m.patternLen == 1 {
		return matchSuccess
	}
	for m.valueLen > 0 {
		patternCopy := make([]rune, len(m.patternRunes[m.patternPos+1:]))
		valueCopy := make([]rune, len(m.valueRunes[m.valuePos:]))
		copy(patternCopy, m.patternRunes[m.patternPos+1:])
		copy(valueCopy, m.valueRunes[m.valuePos:])
		pam := newMatcher(string(patternCopy), string(valueCopy), false)
		if pam.matches() {
			return matchSuccess
		}
		m.valuePos++
		m.valueLen--
	}
	return matchFail
}

// processQuestionMark handles a single character.
func (m *matcher) processQuestionMark() int {
	if m.valueLen == 0 {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// processOpenBracket handles an open bracket for a group of characters.
func (m *matcher) processOpenBracket() int {
	m.patternPos++
	m.patternLen--
	not := (m.patternRunes[m.patternPos] == '^')
	match := false
	if not {
		m.patternPos++
		m.patternLen--
	}
group:
	for {
		switch {
		case m.patternRunes[m.patternPos] == '\\':
			m.patternPos++
			m.patternLen--
			if m.patternRunes[m.patternPos] == m.valueRunes[m.valuePos] {
				match = true
			}
		case m.patternRunes[m.patternPos] == ']':
			break group
		case m.patternLen == 0:
			m.patternPos--
			m.patternLen++
			break group
		case m.patternRunes[m.patternPos+1] == '-' && m.patternLen >= 3:
			start := m.patternRunes[m.patternPos]
			end := m.patternRunes[m.patternPos+2]
			vr := m.valueRunes[m.valuePos]
			if start > end {
				start, end = end, start
			}
			m.patternPos += 2
			m.patternLen -= 2
			if vr >= start && vr <= end {
				match = true
			}
		default:
			if m.patternRunes[m.patternPos] == m.valueRunes[m.valuePos] {
				match = true
			}
		}
		m.patternPos++
		m.patternLen--
	}
	if not {
		match = !match
	}
	if !match {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// processBackslash handles escaping via baskslash.
func (m *matcher) processBackslash() int {
	if m.patternLen >= 2 {
		m.patternPos++
		m.patternLen--
	}
	return matchCont
}

// processDefault handles any other rune.
func (m *matcher) processDefault() int {
	if m.patternRunes[m.patternPos] != m.valueRunes[m.valuePos] {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// Matches checks if the pattern matches a given value.
func Matches(pattern, value string, ignoreCase bool) bool {
	m := newMatcher(pattern, value, ignoreCase)
	return m.matches()
}

// EOF
