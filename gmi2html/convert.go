package gmi2html

import (
	"html"
	"regexp"
	"strings"
)

var heading1Regexp = regexp.MustCompile("^# (.*)$")
var heading2Regexp = regexp.MustCompile("^## (.*)$")
var heading3Regexp = regexp.MustCompile("^### (.*)$")
var linkRegexp = regexp.MustCompile("^=> ([^\\s]+) ?(.+)?$")
var blockquoteRegexp = regexp.MustCompile("^> (.*)$")
var preRegexp = regexp.MustCompile("^```.*$")
var bulletRegexp = regexp.MustCompile(`^\* ?(.*)$`)

func clearLinkMode(linkMode *bool, rv *[]string) {
	if *linkMode {
		*rv = append(*rv, "</p>")
		*linkMode = false
	}
}

func clearUlMode(ulMode *bool, rv *[]string) {
	if *ulMode {
		*rv = append(*rv, "</ul>")
		*ulMode = false
	}
}

func sanitize(input string) string {
	return html.EscapeString(input)
}

func Convert(gmi string) string {
	var rv []string
	preMode := false
	ulMode := false
	linkMode := false
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
			case heading1Regexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				matches := heading1Regexp.FindStringSubmatch(l)
				rv = append(rv, "<h1>"+sanitize(matches[1])+"</h1>")
			case heading2Regexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				matches := heading2Regexp.FindStringSubmatch(l)
				rv = append(rv, "<h2>"+sanitize(matches[1])+"</h2>")
			case heading3Regexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				matches := heading3Regexp.FindStringSubmatch(l)
				rv = append(rv, "<h3>"+sanitize(matches[1])+"</h3>")
			case blockquoteRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				matches := blockquoteRegexp.FindStringSubmatch(l)
				rv = append(rv, "<blockquote>> "+sanitize(matches[1])+"</blockquote>")
			case linkRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				matches := linkRegexp.FindStringSubmatch(l)
				if len(matches[2]) == 0 {
					matches[2] = matches[1]
				}
				if strings.HasSuffix(matches[1], ".png") || strings.HasSuffix(matches[1], ".PNG") || strings.HasSuffix(matches[1], ".jpg") || strings.HasSuffix(matches[1], ".JPG") || strings.HasSuffix(matches[1], ".jpeg") || strings.HasSuffix(matches[1], ".gif") || strings.HasSuffix(matches[1], ".GIF") {
					rv = append(rv, "<img src=\""+sanitize(matches[1])+"\"/>")
					continue
				}
				if linkMode {
					rv = append(rv, "<a href=\""+sanitize(matches[1])+"\">"+sanitize(matches[2])+"</a><br/>")
					continue
				}
				rv = append(rv, "<p><a href=\""+sanitize(matches[1])+"\">"+sanitize(matches[2])+"</a><br/>")
				linkMode = true
			case preRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				rv = append(rv, "<pre>")
				preMode = true
			case bulletRegexp.MatchString(l):
				clearLinkMode(&linkMode, &rv)
				matches := bulletRegexp.FindStringSubmatch(l)
				if ulMode {
					rv = append(rv, "<li>"+sanitize(matches[1])+"</li>")
					continue
				}
				rv = append(rv, "<ul>\n<li>"+sanitize(matches[1])+"</li>")
				ulMode = true
			default:
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				if len(l) != 0 {
					rv = append(rv, "<p>"+sanitize(l)+"</p>")
				}
			}
		}
	}
	clearUlMode(&ulMode, &rv)
	clearLinkMode(&linkMode, &rv)
	return strings.Join(rv, "\n")
}
