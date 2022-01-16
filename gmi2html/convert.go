package gmi2html

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

//var heading1Regexp = regexp.MustCompile("^# (.*)$")
//var heading2Regexp = regexp.MustCompile("^## (.*)$")
//var heading3Regexp = regexp.MustCompile("^### (.*)$")
//var linkRegexp = regexp.MustCompile("^=> ([^\\s]+) ?(.+)?$")
var blockquoteRegexp = regexp.MustCompile("^> (.*)$")
var preRegexp = regexp.MustCompile("^```.*$")
var bulletRegexp = regexp.MustCompile(`^\* ?(.*)$`)
var urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
var imgRegexp = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
var linkRegexp = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

// [Duck Duck Go](https://duckduckgo.com)

func clearLinkMode(linkMode *bool, rv *[]string) {
	//if *linkMode {
	//	*rv = append(*rv, "</p>")
	//	*linkMode = false
	//}
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
			//case heading1Regexp.MatchString(l):
			//	clearUlMode(&ulMode, &rv)
			//	clearLinkMode(&linkMode, &rv)
			//	matches := heading1Regexp.FindStringSubmatch(l)
			//	rv = append(rv, "<h1>"+sanitize(matches[1])+"</h1>")
			//case heading2Regexp.MatchString(l):
			//	clearUlMode(&ulMode, &rv)
			//	clearLinkMode(&linkMode, &rv)
			//	matches := heading2Regexp.FindStringSubmatch(l)
			//	rv = append(rv, "<h2>"+sanitize(matches[1])+"</h2>")
			//case heading3Regexp.MatchString(l):
			//	clearUlMode(&ulMode, &rv)
			//	clearLinkMode(&linkMode, &rv)
			//	matches := heading3Regexp.FindStringSubmatch(l)
			//	rv = append(rv, "<h3>"+sanitize(matches[1])+"</h3>")
			case blockquoteRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				clearLinkMode(&linkMode, &rv)
				matches := blockquoteRegexp.FindStringSubmatch(l)
				rv = append(rv, "<blockquote>> "+sanitize(matches[1])+"</blockquote>")
			//case linkRegexp.MatchString(l):
			//	clearUlMode(&ulMode, &rv)
			//	matches := linkRegexp.FindStringSubmatch(l)
			//	if len(matches[2]) == 0 {
			//		matches[2] = matches[1]
			//	}
			//	if strings.HasSuffix(matches[1], ".png") || strings.HasSuffix(matches[1], ".PNG") || strings.HasSuffix(matches[1], ".jpg") || strings.HasSuffix(matches[1], ".JPG") || strings.HasSuffix(matches[1], ".jpeg") || strings.HasSuffix(matches[1], ".gif") || strings.HasSuffix(matches[1], ".GIF") {
			//		rv = append(rv, "<img src=\""+sanitize(matches[1])+"\"/>")
			//		continue
			//	}
			//	if linkMode {
			//		rv = append(rv, "<a href=\""+sanitize(matches[1])+"\">"+sanitize(matches[2])+"</a><br/>")
			//		continue
			//	}
			//	rv = append(rv, "<p><a href=\""+sanitize(matches[1])+"\">"+sanitize(matches[2])+"</a><br/>")
			//	linkMode = true
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
				sane := sanitize(l)
				//if urlRegexp.MatchString(sane) {
				//	matches := urlRegexp.FindAllStringSubmatch(l, -1)
				//	for _, m := range matches {
				//		url := m[0]
				//		ext := strings.ToLower(filepath.Ext(url))
				//		if ext == ".gif" || ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				//			sane = strings.Replace(sane, url, fmt.Sprintf("<img src=\"%s\"/>", url), 1)
				//		} else {
				//			sane = strings.Replace(sane, url, fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, url), 1)
				//		}
				//	}
				//}
				if imgRegexp.MatchString(sane) || linkRegexp.MatchString(sane) {
					matches := imgRegexp.FindAllStringSubmatch(sane, -1)
					for _, m := range matches {
						sane = strings.Replace(sane, m[0], fmt.Sprintf("<img src=\"%s\" alt=\"%s\"/>", m[2], m[1]), 1)
					}
					matches = linkRegexp.FindAllStringSubmatch(sane, -1)
					for _, m := range matches {
						sane = strings.Replace(sane, m[0], fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", m[2], m[1]), 1)
					}
				} else if urlRegexp.MatchString(sane) {
					matches := urlRegexp.FindAllStringSubmatch(l, -1)
					for _, m := range matches {
						url := m[0]
						sane = strings.Replace(sane, url, fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, url), 1)
					}
				}
				if len(sane) != 0 {
					rv = append(rv, "<p>"+sane+"</p>")
				}
			}
		}
	}
	clearUlMode(&ulMode, &rv)
	clearLinkMode(&linkMode, &rv)
	return strings.Join(rv, "\n")
}
