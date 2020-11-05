package formats

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type vCombinedLine struct {
	re *regexp.Regexp
}

func (v vCombinedLine) ParseLine(line string) (Line, error) {
	if v.re == nil {
		panic("no regexp available")
	}
	matchString := v.re.FindStringSubmatch(line)
	timeString := ""
	if len(matchString) > 0 {
		timeString = matchString[1]
	}
	parsedTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeString)
	if err != nil {
		return Line{}, errors.WithMessagef(err, "Failed to parse timestamp %q", timeString)
	}
	return Line{Source: line, Timestamp: parsedTime}, nil
}

func NewVCombinedParser() (LineReader, error) {
	const vCombinedRegexp = "^[\\w-.]+:\\d+ \\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3} - .+ \\[(.+)\\] \".*\" \\d{3} \\d+ \".*\"( \".*\")?( \".*\")?"
	compile := regexp.MustCompile(vCombinedRegexp)
	return vCombinedLine{re: compile}, nil
}
