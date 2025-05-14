// Package customBlackfriday is a customized version of the Blackfriday Markdown renderer
// that adds support for the vpub image proxy service.
//
// To use the vpub image proxy service, you only need to replace the URL of an image.
// For example, a URL such as "https://i.imgur.com/404" will be replaced with
// "http://localhost:1337/image-proxy?url=https://i.imgur.com/404"
//
// To see source code of image proxy, check /web/handler/image_proxy.go
package customBlackfriday

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io"
	"net/url"
	"strings"
	"vpub/syntax"
)

type BlackfridayRenderer struct{}

type imageProxyRenderer struct {
	blackfriday.Renderer
}

func (r *imageProxyRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.Image && entering {
		originalURL := string(node.LinkData.Destination)
		escapedURL := url.QueryEscape(originalURL)
		node.LinkData.Destination = []byte(fmt.Sprintf("%s%s", syntax.ImageProxyURLBase, escapedURL))
	}
	return r.Renderer.RenderNode(w, node, entering)
}

func (b *BlackfridayRenderer) Convert(gmiContent string, wrap bool) string {
	input := strings.ReplaceAll(gmiContent, "\r\n", "\n")
	params := blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	}
	base := blackfriday.NewHTMLRenderer(params)
	proxy := &imageProxyRenderer{Renderer: base}

	unsafe := blackfriday.Run(
		[]byte(input),
		blackfriday.WithRenderer(proxy),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(safe)
}
