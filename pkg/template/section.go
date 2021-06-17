package template

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/pkg/errors"
)

// section is a named section of a Template
type section struct {
	name     string
	contents []contentItem
}

// newSection constructs a Section
func newSection(name string, items ...contentItem) *section {
	return &section{
		name:     name,
		contents: items,
	}
}

func (s *section) deleteContents() {
	s.contents = []contentItem{}
}

func (s *section) sortContents() {
	// stable sort to preserve order for empty header case
	sort.SliceStable(s.contents, func(i, j int) bool {
		return s.contents[i].header < s.contents[j].header
	})
}

func (s *section) isEmpty() bool {
	for _, content := range s.contents {
		if !content.isEmpty() {
			return false
		}
	}
	return true
}

func (s *section) getNameString(prefix string, suffix string) string {
	return fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
}

func (s *section) getContentString() string {
	str := ""
	for _, content := range s.contents {
		txt := content.string()
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		str += txt
	}
	return str
}

type contentItem struct {
	header string
	text   string
}

func (ci contentItem) string() string {
	if ci.header != "" {
		return fmt.Sprintf("%s\n%s", ci.header, ci.text)
	}
	return ci.text
}

func (ci contentItem) isEmpty() bool {
	// exclude trailing newlines for empty content check
	strippedTxt := strings.Replace(ci.text, "\n", "", -1)
	return len(strippedTxt) == 0
}

func parseSection(text string, opts config.Opts) (*section, error) {
	if len(text) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	lines := strings.Split(text, "\n")
	name := stripPrefixSuffix(lines[0], opts.Section.Prefix, opts.Section.Suffix)
	contents := parseSectionContents(
		lines[1:],
		opts.Archive.SectionContentPrefix,
		opts.Archive.SectionContentSuffix,
		opts.File.TimeFormat,
	)

	// return section populated with contents if any contentItem is non-empty
	for _, content := range contents {
		if !content.isEmpty() {
			return newSection(name, contents...), nil
		}
	}
	// all contents are empty so return unpopulated section
	return newSection(name), nil
}

func parseSectionContents(lines []string, prefix string, suffix string, format string) []contentItem {
	contents := []contentItem{}
	if len(lines) == 0 {
		return contents
	}

	// parse first line
	line := lines[0]
	header := ""
	body := []string{}
	if isArchiveItemHeader(line, prefix, suffix, format) {
		header = line
	} else {
		body = append(body, line)
	}

	for _, line := range lines[1:] {
		// if the line is a header it indicates new contents, so "flush" (append) the current
		// header/body and start tracking the new contents
		if isArchiveItemHeader(line, prefix, suffix, format) {
			contents = append(contents, contentItem{
				header: header,
				text:   strings.Join(body, "\n"),
			})

			header = line
			body = []string{}
			continue
		}

		body = append(body, line)
	}

	// ensure remaining content is appended
	if len(body) != 0 || header != "" {
		contents = append(contents, contentItem{
			header: header,
			text:   strings.Join(body, "\n"),
		})
	}
	return contents
}

func stripPrefixSuffix(line string, prefix string, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(line, suffix), prefix)
}

func getSectionNameRegex(prefix string, suffix string) (*regexp.Regexp, error) {
	sectionPattern := fmt.Sprintf("%s.*%s", prefix, suffix)
	sectionNameRegex, err := regexp.Compile(sectionPattern)
	if err != nil {
		return sectionNameRegex, fmt.Errorf("invalid section prefix [%s] or suffix [%s]", prefix, suffix)
	}
	return sectionNameRegex, nil
}
