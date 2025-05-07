//go:build blackfriday
// +build blackfriday

package syntax

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	imageProxyURLBase = "/image-proxy?url="
)

// ImageProxyRenderer is a blackfriday.Renderer that modifies image URLs
// to route them through an image proxy
type ImageProxyRenderer struct {
	blackfriday.Renderer
}

// RenderNode customizes the rendering of image nodes to prepend an image proxy URL
func (r *ImageProxyRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.Image && entering {
		originalURL := string(node.LinkData.Destination)
		escapedURL := url.QueryEscape(originalURL)
		node.LinkData.Destination = []byte(fmt.Sprintf("%s%s", imageProxyURLBase, escapedURL))
	}
	return r.Renderer.RenderNode(w, node, entering)
}

func Convert(gmiContent string, wrap bool) string {
	normalizedInput := strings.ReplaceAll(gmiContent, "\r\n", "\n")

	htmlRendererParams := blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	}

	baseRenderer := blackfriday.NewHTMLRenderer(htmlRendererParams)
	proxyRenderer := &ImageProxyRenderer{Renderer: baseRenderer}

	unsafeHTMLBytes := blackfriday.Run(
		[]byte(normalizedInput),
		blackfriday.WithRenderer(proxyRenderer),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)

	sanitizationPolicy := bluemonday.UGCPolicy()
	sanitizedHTMLBytes := sanitizationPolicy.SanitizeBytes(unsafeHTMLBytes)

	return string(sanitizedHTMLBytes)
}
