package tokenizer

import "regexp"

type Stream struct {
	pos       int
	str       string
	tabSize   int
	lineStart int
	start     int
}

func (s *Stream) Eol() bool {
	return s.pos >= len(s.str)
}

func (s *Stream) Sol() bool {
	return s.pos == s.lineStart
}

func (s *Stream) Current() string {
	if s.pos < len(s.str) {
		return s.str[s.start:s.pos]
	}

	return s.str[s.start:]
}

func (s *Stream) EatSpace() bool {
	start := s.pos
	rSpace := regexp.MustCompile(`[\s\xa0]`)
	for {

		if s.Eol() || !rSpace.Match([]byte{s.str[s.pos]}) {
			break
		}
		s.pos++
	}
	return s.pos > start
}

func (s *Stream) Next() string {
	if s.pos < len(s.str) {
		s.pos++
		return string(s.str[s.pos-1])
	}
	return ""
}

func (s *Stream) Match(pattern string, consume bool) string {
	regex := regexp.MustCompile(pattern)
	match := regex.FindString(s.str[s.pos:])
	if match != "" && consume {
		s.pos += len(match)
	}
	return match
}

func (s *Stream) EatString(match string) string {
	if s.pos >= len(s.str) {
		return ""
	}
	ch := string(s.str[s.pos])
	if ch == match {
		s.pos++
		return ch
	}
	return ""
}

func (s *Stream) EatRegex(regex *regexp.Regexp) string {
	if s.pos >= len(s.str) {
		return ""
	}
	ch := string(s.str[s.pos])
	if ch != "" && regex.MatchString(ch) {
		s.pos++
		return ch
	}
	return ""
}

func (s *Stream) SkipToEnd() {
	s.pos = len(s.str)
}

func (s *Stream) EatWhile(pattern *regexp.Regexp) bool {
	start := s.pos
	for {
		eated := s.EatRegex(pattern)
		if eated == "" {
			break
		}
	}
	return s.pos > start
}
