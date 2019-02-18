// Copyright 2019 Dmitry A. Mottl. All rights reserved.
// Use of this source code is governed by MIT license
// that can be found in the LICENSE file.

// Package ctimefmt provides functionality for formating and parsing time.Time
// struct using legacy ctime-like syntax such as: "%Y-%m-%d %H:%M:%S %Z"
package ctimefmt

import (
	"regexp"
	"time"
)

var ctimeRegexp, decimalsRegexp *regexp.Regexp
var ctimeSubstitutes map[string]string

func init() {
	// ctime format -> Go format conversion
	ctimeSubstitutes = make(map[string]string)
	ctimeSubstitutes["%Y"] = "2006"
	ctimeSubstitutes["%y"] = "06"
	ctimeSubstitutes["%m"] = "01"
	ctimeSubstitutes["%b"] = "Jan"
	ctimeSubstitutes["%h"] = "Jan"
	ctimeSubstitutes["%B"] = "January"
	ctimeSubstitutes["%d"] = "02"
	ctimeSubstitutes["%e"] = "_2"
	ctimeSubstitutes["%a"] = "Mon"
	ctimeSubstitutes["%A"] = "Monday"
	ctimeSubstitutes["%H"] = "15"
	ctimeSubstitutes["%I"] = "03"
	ctimeSubstitutes["%p"] = "PM"
	ctimeSubstitutes["%M"] = "04"
	ctimeSubstitutes["%S"] = "05"
	ctimeSubstitutes["%f"] = "999999"
	ctimeSubstitutes["%z"] = "-0700"
	ctimeSubstitutes["%Z"] = "MST"

	ctimeSubstitutes["%D"] = "01/02/2006"
	ctimeSubstitutes["%x"] = "01/02/2006"
	ctimeSubstitutes["%F"] = "2006-01-02"
	ctimeSubstitutes["%T"] = "15:04:05"
	ctimeSubstitutes["%X"] = "15:04:05"
	ctimeSubstitutes["%r"] = "03:04:05 pm"
	ctimeSubstitutes["%R"] = "15:04"
	ctimeSubstitutes["%n"] = "\n"
	ctimeSubstitutes["%t"] = "\t"
	ctimeSubstitutes["%%"] = "%"
	ctimeSubstitutes["%c"] = "Mon Jan 02 15:04:05 2006"

	ctimeRegexp = regexp.MustCompile(`%.`)
	decimalsRegexp = regexp.MustCompile(`\d`)
}

// Format returns a textual representation of the time value formatted
// according to ctime-like format string. Possible directives are:
//   %Y - Year, zero-padded (0001, 0002, ..., 2019, 2020, ..., 9999)
//   %y - Year, last two digits, zero-padded (01, ..., 99)
//   %m - Month as a decimal number (01, 02, ..., 12)
//   %b, %h - Abbreviated month name (Jan, Feb, ...)
//   %B - Full month name (January, February, ...)
//   %d - Day of the month, zero-padded (01, 02, ..., 31)
//   %e - Day of the month, space-padded ( 1, 2, ..., 31)
//   %a - Abbreviated weekday name (Sun, Mon, ...)
//   %A - Full weekday name (Sunday, Monday, ...)
//   %H - Hour (24-hour clock) as a zero-padded decimal number (00, ..., 24)
//   %I - Hour (12-hour clock) as a zero-padded decimal number (00, ..., 12)
//   %p - Locale’s equivalent of either AM or PM
//   %M - Minute, zero-padded (00, 01, ..., 59)
//   %S - Second as a zero-padded decimal number (00, 01, ..., 59)
//   %f - Microsecond as a decimal number, zero-padded on the left (00, 01, ..., 59)
//   %z - UTC offset in the form ±HHMM[SS[.ffffff]] or empty(+0000, -0400)
//   %Z - Timezone name or abbreviation or empty (UTC, EST, CST)
//   %D, %x - Short MM/DD/YY date, equivalent to %m/%d/%y
//   %F - Short YYYY-MM-DD date, equivalent to %Y-%m-%d
//   %T, %X - ISO 8601 time format (HH:MM:SS), equivalent to %H:%M:%S
//   %r - 12-hour clock time (02:55:02 pm)
//   %R - 24-hour HH:MM time, equivalent to %H:%M
//   %n - New-line character ('\n')
//   %t - Horizontal-tab character ('\t')
//   %% - A % sign
//   %c - Date and time representation (Mon Jan 02 15:04:05 2006)
func Format(format string, t time.Time) string {
	return t.Format(ToNative(format))
}

// Parse parses a ctime-like formatted string (e.g. "%Y-%m-%d ...") and returns
// the time value it represents.
//
// Refer to Format() function documentation for possible directives.
func Parse(format, value string) (time.Time, error) {
	return time.Parse(ToNative(format), value)
}

// ToNative converts ctime's format string to Go native layout
// (which is used by time.Time.Format() and time.Parse() functions).
func ToNative(format string) string {
	if match := decimalsRegexp.FindString(format); match != "" {
		panic("Format string should not contain decimals")
	}

	replaceFunc := func(directive string) string {
		if subst, ok := ctimeSubstitutes[directive]; ok {
			return subst
		} else {
			panic("Unsupported ctimefmt.ToNative() directive: " + directive)
		}
	}

	return ctimeRegexp.ReplaceAllStringFunc(format, replaceFunc)
}