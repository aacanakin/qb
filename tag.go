package qb

import (
	"fmt"
	"strings"
)

// Tag is the base abstraction of qbit tag
type Tag struct {

	// contains default, null, notnull, unique, primary_key, foreign_key(table.column), check(condition > 0)
	Constraints []string

	// contains type(size) or type parameters
	Type string

	// true if it is ignored
	Ignore bool
}

// ParseTag parses raw qbit tag and builds a Tag object
func ParseTag(rawTag string) (*Tag, error) {

	rawTag = strings.Trim(rawTag, " ")

	tag := &Tag{
		Constraints: []string{},
	}

	if rawTag == "" {
		return tag, nil
	}

	tags := strings.Split(rawTag, ";")
	for _, t := range tags {

		tagKeyVal := strings.Split(t, ":")

		if tagKeyVal[0] == "index" {

			tag.Constraints = append(tag.Constraints, t)
			continue

		} else if tagKeyVal[0] == "-" {
			tag.Ignore = true
			return tag, nil
		}

		if len(tagKeyVal) != 2 {
			return nil, fmt.Errorf("Invalid tag key length, tag: %v", tag)
		}

		if tagKeyVal[0] == "type" {
			tag.Type = tagKeyVal[1]
		} else if tagKeyVal[0] == "constraints" || tagKeyVal[0] == "constraint" {
			for _, c := range strings.Split(tagKeyVal[1], ",") {
				if c != "" {
					tag.Constraints = append(tag.Constraints, c)
				}
			}
		}
	}

	return tag, nil
}
