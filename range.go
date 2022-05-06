package prism

import (
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	End   int
}

func (r *Range) Parse(lines string) error {
	bounds := strings.Split(lines, "-")

	start, err := strconv.Atoi(bounds[0])
	if err != nil {
		return err
	}

	end, err := strconv.Atoi(bounds[1])
	if err != nil {
		return err
	}

	r.Start = start
	r.End = end

	return nil
}

func substr(code string, lineRange Range) (string, error) {
	lines := strings.Split(code, "\n")

	if lineRange.Start < 1 || lineRange.End > len(lines) {
		return "", fmt.Errorf("Range outside of bounds")
	}

	return strings.Join(lines[lineRange.Start-1:lineRange.End], "\n"), nil
}
