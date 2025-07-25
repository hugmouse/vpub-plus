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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"
	"vpub/syntax"
)

const cacheTTL = 5 * time.Minute
const cleanupInterval = 1 * time.Minute

var rendererCache sync.Map

type cacheEntry struct {
	data      string
	timestamp time.Time
}

func init() {
	go startCacheCleanup()
}

func startCacheCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	for range ticker.C {
		cleanupCache()
	}
}

func cleanupCache() {
	now := time.Now()
	rendererCache.Range(func(key, value any) bool {
		if entry, ok := value.(cacheEntry); ok {
			if now.Sub(entry.timestamp) > cacheTTL {
				rendererCache.Delete(key)
			}
		}
		return true
	})
}

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
	h := sha1.New()
	if _, err := h.Write([]byte(gmiContent)); err != nil {
		return fmt.Errorf("SHA-1 hash error: %w", err).Error()
	}
	sum := h.Sum(nil)

	cacheKey := hex.EncodeToString(sum)

	if cached, ok := rendererCache.Load(cacheKey); ok {
		entry := cached.(cacheEntry)
		if time.Since(entry.timestamp) < cacheTTL {
			return entry.data // Cache hit!
		}
	}

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

	rendererCache.Store(cacheKey, cacheEntry{
		data:      string(safe),
		timestamp: time.Now(),
	})
	return string(safe)
}

func (b *BlackfridayRenderer) Name() string {
	return "blackfriday"
}
