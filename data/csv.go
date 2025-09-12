package data

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"preselect/business"
)

// CSVLoader retrieves CSV data entry by entry.
type CSVLoader struct {
	reader  *bufio.Reader
	delim   rune
	quote   rune
	row     int
	col     int
	inQuote bool
	token   []rune
	eof     bool
}

// NewCSVLoader creates a new CSVLoader. If r is nil, an empty reader is used.
// Delim defaults to comma and quote to double quote when zero values are provided.
func NewCSVLoader(r io.Reader, delim, quote rune) *CSVLoader {
	if r == nil {
		r = strings.NewReader("")
	}
	if delim == 0 {
		delim = ','
	}
	if quote == 0 {
		quote = '"'
	}
	return &CSVLoader{
		reader: bufio.NewReader(r),
		delim:  delim,
		quote:  quote,
		row:    1,
	}
}

// Next returns the next cell from the CSV source.
func (l *CSVLoader) Next() (business.Entry, error) {
	for {
		if l.eof {
			if l.token != nil {
				l.col++
				entry := business.Entry{
					Value: string(l.token),
					Path:  []string{strconv.Itoa(l.row), strconv.Itoa(l.col)},
				}
				l.token = nil
				return entry, nil
			}
			return business.Entry{}, io.EOF
		}

		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.eof = true
				continue
			}
			return business.Entry{}, err
		}

		if l.inQuote {
			if r == l.quote {
				next, _, err := l.reader.ReadRune()
				if err == nil {
					if next == l.quote {
						l.token = append(l.token, l.quote)
						continue
					}
					l.reader.UnreadRune()
				} else if err == io.EOF {
					l.eof = true
				} else {
					return business.Entry{}, err
				}
				l.inQuote = false
				continue
			}
			l.token = append(l.token, r)
			continue
		}

		switch r {
		case l.quote:
			l.inQuote = true
		case l.delim:
			l.col++
			entry := business.Entry{
				Value: string(l.token),
				Path:  []string{strconv.Itoa(l.row), strconv.Itoa(l.col)},
			}
			l.token = nil
			return entry, nil
		case '\n':
			l.col++
			entry := business.Entry{
				Value: string(l.token),
				Path:  []string{strconv.Itoa(l.row), strconv.Itoa(l.col)},
			}
			l.token = nil
			l.row++
			l.col = 0
			return entry, nil
		case '\r':
			next, _, err := l.reader.ReadRune()
			if err == nil {
				if next != '\n' {
					l.reader.UnreadRune()
				}
			} else if err != io.EOF {
				return business.Entry{}, err
			} else {
				l.eof = true
			}
			l.col++
			entry := business.Entry{
				Value: string(l.token),
				Path:  []string{strconv.Itoa(l.row), strconv.Itoa(l.col)},
			}
			l.token = nil
			l.row++
			l.col = 0
			return entry, nil
		default:
			l.token = append(l.token, r)
		}
	}
}
