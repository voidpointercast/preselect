package data

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"preselect/business"
)

// Loader retrieves file data entry by entry using token-based scanning.
type Loader struct {
	reader      *bufio.Reader
	delimiters  map[rune]struct{}
	token       []rune
	inToken     bool
	tokenStart  int64
	position    int64
	tokenNumber int
	line        int
	eof         bool
}

// NewLoader creates a new Loader. If r is nil, an empty reader is used. Delimiters
// default to space and newline when none are provided.
func NewLoader(r io.Reader, delimiters []rune) *Loader {
	if r == nil {
		r = strings.NewReader("")
	}
	if len(delimiters) == 0 {
		delimiters = []rune{' ', '\n'}
	}
	delimMap := make(map[rune]struct{}, len(delimiters))
	for _, d := range delimiters {
		delimMap[d] = struct{}{}
	}
	return &Loader{
		reader:     bufio.NewReader(r),
		delimiters: delimMap,
		line:       1,
	}
}

// Next returns the next entry from the source based on tokenization.
func (l *Loader) Next() (business.Entry, error) {
	for {
		if l.eof {
			if l.inToken {
				l.tokenNumber++
				entry := business.Entry{
					Value: string(l.token),
					Path: []string{
						strconv.Itoa(l.tokenNumber),
						strconv.FormatInt(l.tokenStart, 10),
						strconv.Itoa(l.line),
					},
				}
				l.token = nil
				l.inToken = false
				return entry, nil
			}
			return business.Entry{}, io.EOF
		}

		r, size, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.eof = true
				continue
			}
			return business.Entry{}, err
		}

		currentPos := l.position
		l.position += int64(size)

		if _, ok := l.delimiters[r]; ok {
			if l.inToken {
				l.tokenNumber++
				entry := business.Entry{
					Value: string(l.token),
					Path: []string{
						strconv.Itoa(l.tokenNumber),
						strconv.FormatInt(l.tokenStart, 10),
						strconv.Itoa(l.line),
					},
				}
				l.token = nil
				l.inToken = false
				if r == '\n' {
					l.line++
				}
				return entry, nil
			}
			if r == '\n' {
				l.line++
			}
			continue
		}

		if !l.inToken {
			l.inToken = true
			l.tokenStart = currentPos
		}
		l.token = append(l.token, r)
	}
}
