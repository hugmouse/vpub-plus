//go:build !blackfriday
// +build !blackfriday

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
var urlRegexp = regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)`)
var imgRegexp = regexp.MustCompile(`!\[(.*?)]\((.*?)\)`)
var linkRegexp = regexp.MustCompile(`\[(.*?)]\((.*?)\)`)
var boldRegexp = regexp.MustCompile(`\*\*(.*?)\*\*`)
var italicsRegexp = regexp.MustCompile(`\*(.*?)\*`)
var tableLikeHeader = regexp.MustCompile(`^\|\s.+\s\|$`)
var tableSeparator = regexp.MustCompile(`(:?-.-+:?)`)
var codeRegexp = regexp.MustCompile("`(.*)`")
var strikethroughRegexp = regexp.MustCompile(`~~(.*?)~~`)

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

func processCode(input string) string {
	if codeRegexp.MatchString(input) {
		matches := codeRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			input = strings.Replace(input, m[0], fmt.Sprintf("<code>%s</code>", m[1]), 1)
		}
	}
	return input
}

func processStrikethrough(input string) string {
	if strikethroughRegexp.MatchString(input) {
		matches := strikethroughRegexp.FindAllStringSubmatch(input, -1)
		for _, m := range matches {
			input = strings.Replace(input, m[0], fmt.Sprintf("<s>%s</s>", m[1]), 1)
		}
	}
	return input
}

// Returns a sanitized output
func processDecoration(input string) string {
	sane := processLinks(input)
	sane = processBold(sane)
	sane = processItalics(sane)
	sane = processCode(sane)
	sane = processStrikethrough(sane)
	return sane
}

func Convert(gmi string, wrap bool) string {
	var rv []string

	// Table parser logic
	var tableHeaderTmp string
	var tableHeaderAlreadyBuilt = false
	var tableBuilder strings.Builder
	var tableCenteredRows []int
	var tableRightAlignedRows []int
	var inCurrentTableTBodyIsAlreadyExists = false
	tableMode := false

	preMode := false
	ulMode := false

	// Remove \r from existence
	gmi = strings.ReplaceAll(gmi, "\r\n", "\n")
	separatedGmi := strings.Split(gmi, "\n")

	for index, l := range separatedGmi {
		if len(l) >= 4 && l[0:4] == "----" {
			rv = append(rv, "<hr>")
			continue
		}

		// Because of the semantics of the current document,
		// there a <h1> that already exists with the topic
		// name, and we can't use multiple <h1> in one document.
		//
		// So # means <h2>, ## = <h3> and so on.
		//
		// Also this implementation allows this syntax: "#header",
		// which is not typically allowed.
		if len(l) >= 1 && l[0:1] == "#" {
			c := 0
			for i, char := range l {
				if char == '#' {
					continue
				} else {
					c = i
					break
				}
			}
			if c >= 1 && c <= 5 {
				rv = append(rv, fmt.Sprintf("<h%d>%v</h%d>", c+1, l[c:], c+1))
			}
			continue
		}

		// If tableMode detected, then we sure hope that the current string
		// is either a header separator, or a continuous table
		if tableMode {
			tmpIsThisAHeaderSeparator := false
			// Remembering what rows should be centered or right aligned
			for i, match := range tableSeparator.FindAllStringSubmatch(l, -1) {
				tmpIsThisAHeaderSeparator = true
				if match[0][0:1] == ":" && match[0][len(match[0])-2:] == "-:" {
					tableCenteredRows = append(tableCenteredRows, i)
				} else if match[0][0:1] == "-" && match[0][len(match[0])-2:] == "-:" {
					tableRightAlignedRows = append(tableRightAlignedRows, i)
				}
			}

			if !tableHeaderAlreadyBuilt {
				tableBuilder.WriteString("<table><thead><tr>")
				sep := strings.Split(tableHeaderTmp, " | ")
				for i := 0; i < len(sep); i++ {
					sep[i] = strings.Trim(sep[i], "|")
					tmpAligned := false
					for _, centered := range tableCenteredRows {
						if i == centered {
							tableBuilder.WriteString("<td align=\"center\">" + sep[i] + "</td>")
							tmpAligned = true
							break
						}
					}
					for _, rightAlign := range tableRightAlignedRows {
						if i == rightAlign {
							tableBuilder.WriteString("<td align=\"right\">" + sep[i] + "</td>")
							tmpAligned = true
							break
						}
					}
					if !tmpAligned {
						tableBuilder.WriteString("<td>" + sep[i] + "</td>")
						tmpAligned = false
					}
				}
				tableBuilder.WriteString("</tr></thead>")
				tableHeaderAlreadyBuilt = true
			}

			// Usually after the table ends there is an empty string with a newline in it
			if !tmpIsThisAHeaderSeparator && len(l) > 2 {
				if !inCurrentTableTBodyIsAlreadyExists {
					tableBuilder.WriteString("<tbody>")
					inCurrentTableTBodyIsAlreadyExists = true
				}

				sep := strings.Split(l, " | ")
				tableBuilder.WriteString("<tr>")
				for i := 0; i < len(sep); i++ {
					sep[i] = strings.Trim(sep[i], "|")
					tmpAligned := false
					for _, centered := range tableCenteredRows {
						if i == centered {
							tableBuilder.WriteString("<td align=\"center\">" + sep[i] + "</td>")
							tmpAligned = true
							break
						}
					}
					for _, rightAlign := range tableRightAlignedRows {
						if i == rightAlign {
							tableBuilder.WriteString("<td align=\"right\">" + sep[i] + "</td>")
							tmpAligned = true
							break
						}
					}

					if !tmpAligned {
						tableBuilder.WriteString("<td>" + sep[i] + "</td>")
						tmpAligned = false
					}
				}
				tableBuilder.WriteString("</tr>")
			}

			// Sometimes tables ends without an \n at the end!
			// This stuff checks if we are on the end of our soulless existence
			// and if we are (at the end of the string without an \n) - we set
			// our table mode to false
			if len(separatedGmi) == index+1 {
				tableMode = false
			}

			// This triggers on the empty string with a newline on the end
			// Also this means that we have to close our tbody and table tags
			if !tmpIsThisAHeaderSeparator && !tableLikeHeader.MatchString(l) || !tableMode {
				tableMode = false
				inCurrentTableTBodyIsAlreadyExists = false
				tableHeaderAlreadyBuilt = false
				tableBuilder.WriteString("</tbody>")
				tableBuilder.WriteString("</table>")
				rv = append(rv, tableBuilder.String())
				tableBuilder.Reset()
			}

			continue
		}
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
			case tableLikeHeader.MatchString(l):
				tableMode = true
				tableHeaderTmp = l
			case blockquoteRegexp.MatchString(l):
				clearUlMode(&ulMode, &rv)
				matches := blockquoteRegexp.FindStringSubmatch(l)
				rv = append(rv, "<blockquote>"+sanitize(matches[1])+"</blockquote>")
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
