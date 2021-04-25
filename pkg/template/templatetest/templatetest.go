// Package templatetest provides utilities for template testing
package templatetest

import (
	"fmt"
	"log"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
)

// Date is a fixed date - changing this value will affect some tests
var Date = time.Date(2020, 12, 20, 1, 1, 1, 1, time.UTC)

// GetOpts returns a configuration struct for tests - changing these values will affect some tests
func GetOpts() config.Opts {
	opts := config.Opts{
		AppDir: "path/to/app/dir",
		Header: config.HeaderOpts{
			Prefix:           "-^-",
			Suffix:           "-v-",
			TrailingNewlines: 1,
			TimeFormat:       "[Mon] 02 Jan 2006",
		},
		Section: config.SectionOpts{
			Prefix:           "_p_",
			Suffix:           "_q_",
			TrailingNewlines: 3,
			Names: []string{
				"TestSection1",
				"TestSection2",
				"TestSection3",
			},
		},
		File: config.FileOpts{
			Ext:        "txt",
			TimeFormat: "2006-01-02",
			CursorLine: 1,
		},
		Archive: config.ArchiveOpts{
			AfterDays:                7,
			FilePrefix:               "archive-",
			HeaderPrefix:             "ARCHIVEPREFIX ",
			HeaderSuffix:             " ARCHIVESUFFIX",
			SectionContentPrefix:     "[",
			SectionContentSuffix:     "]",
			SectionContentTimeFormat: "2006-01-02",
			MonthTimeFormat:          "Jan2006",
		},
		Cli: config.CliOpts{
			TimeFormat: "2006-01-02",
		},
		TemplateFileCountThresh: 90,
	}

	err := config.ValidateOpts(opts)
	if err != nil {
		log.Fatal(err)
	}
	return opts
}

// MakeItemHeader is a helper to construct a header property of a contentItem struct
func MakeItemHeader(date time.Time, opts config.Opts) string {
	return fmt.Sprintf("%s%s%s",
		opts.Archive.SectionContentPrefix,
		date.Format(opts.Archive.SectionContentTimeFormat),
		opts.Archive.SectionContentSuffix,
	)
}
