package formats

import "time"

type Line struct {
	Source    string
	Timestamp time.Time
}

type LineReader interface {
	ParseLine(line string) (Line, error)
}
