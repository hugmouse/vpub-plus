package syntax

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var blockquoteRegexp = regexp.MustCompile("^> (.*)$")
var preRegexp = regexp.MustCompile("^```.*$")
var bulletRegexp = regexp.MustCompile(`^\* (.*)$`)
var urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
var imgRegexp = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
var linkRegexp = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
var boldRegexp = regexp.MustCompile(`\*\*(.*?)\*\*`)
var italicsRegexp = regexp.MustCompile(`\*(.*?)\*`)

func clearUlMode(ulMode *bool, rv *[]string) {
	if *ulMode {
		*rv = append(*rv, "</ul>")
		*ulMode = false
	}
}

func sanitize(input string) string {
	return html.EscapeString(input)
}

func processLinks(input string) string {
	sane := html.EscapeString(input)
	if imgRegexp.MatchString(input) || linkRegexp.MatchString(input) {
		matches := imgRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			sane = strings.Replace(sane, m[0], fmt.Sprintf("<img src=\"%s\" alt=\"%s\"/>", m[2], m[1]), 1)
		}
		matches = linkRegexp.FindAllStringSubmatch(sane, -1)
		for _, m := range matches {
			sane = strings.Replace(sane, m[0], fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", m[2], m[1]), 1)
		}
	} else if urlRegexp.MatchString(input) {
		matches := urlRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			url := m[0]
			sane = strings.Replace(sane, url, fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, url), 1)
		}
	}
	return sane
}

func processBold(input string) string {
	if boldRegexp.MatchString(input) {
		matches := boldRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			input = strings.Replace(input, m[0], fmt.Sprintf("<b>%s</b>", m[1]), 1)
		}
	}
	return input
}

func processItalics(input string) string {
	if italicsRegexp.MatchString(input) {
		matches := italicsRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			input = strings.Replace(input, m[0], fmt.Sprintf("<i>%s</i>", m[1]), 1)
		}
	}
	return input
}

// Returns a sanitized output
func processDecoration(input string) string {
	sane := processLinks(input)
	sane = processBold(sane)
	sane = processItalics(sane)
	return sane
}

func Convert(gmi string, wrap bool) string {
	var rv []string
	preMode := false
	ulMode := false
	for _, l := range strings.Split(gmi, "\n") {
		l = strings.TrimRight(l, "\r")
		if preMode {
			switch {
			case preRegexp.MatchString(l):
				rv = append(rv, "</pre>")
				preMode = false
			default:
				rv = append(rv, sanitize(l))
			}
		} else {
			switch {
			case blockquoteRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				matches := blockquoteRegexp.FindStringSubmatch(l)
				rv = append(rv, "<blockquote>> "+sanitize(matches[1])+"</blockquote>")
			case preRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				rv = append(rv, "<pre>")
				preMode = true
			case bulletRegexp.MatchString(l):
				matches := bulletRegexp.FindStringSubmatch(l)
				sane := processDecoration(matches[1])
				if ulMode {
					rv = append(rv, "<li>"+sane+"</li>")
					continue
				}
				rv = append(rv, "<ul>\n<li>"+sane+"</li>")
				ulMode = true
			default:
				clearUlMode(&ulMode, &rv)
				sane := processDecoration(l)
				if len(l) != 0 {
					if wrap {
						rv = append(rv, "<p>"+sane+"</p>")
					} else {
						rv = append(rv, sane)
					}
				}
			}
		}
	}
	clearUlMode(&ulMode, &rv)
	return strings.Join(rv, "\n")
}
