package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const trackingSeparator = "|"

// EncodeTrackingToken encodes a URL + user public-ID pair for use in a tracking link.
func EncodeTrackingToken(originalURL, userPublicId string) string {
	raw := originalURL + trackingSeparator + userPublicId
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// DecodeTrackingToken reverses EncodeTrackingToken.
func DecodeTrackingToken(token string) (originalURL, userPublicId string, ok bool) {
	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(b), trackingSeparator, 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

var hrefRegex = regexp.MustCompile(`(?i)(href=["'])([^"']+)(["'])`)

// RewriteLinksForTracking replaces every href in the HTML body with a
// tracking redirect URL so that clicks are recorded.
func RewriteLinksForTracking(html, userPublicId, baseURL string) string {
	if baseURL == "" || userPublicId == "" {
		return html
	}
	return hrefRegex.ReplaceAllStringFunc(html, func(match string) string {
		parts := hrefRegex.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		originalURL := parts[2]
		if strings.HasPrefix(originalURL, "mailto:") || strings.HasPrefix(originalURL, "#") {
			return match
		}
		token := EncodeTrackingToken(originalURL, userPublicId)
		trackingURL := strings.TrimRight(baseURL, "/") + "/api/track/click/" + token
		return parts[1] + trackingURL + parts[3]
	})
}

// HandleClickTracking decodes the tracking token and redirects to the original URL.
// In the monolith this also fires an internal webhook; here it's a stub that
// just performs the redirect.
func HandleClickTracking(c *gin.Context) {
	token := c.Param("token")

	originalURL, _, ok := DecodeTrackingToken(token)
	if !ok {
		log.Printf("invalid tracking token: %s", token)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// TODO: Look up user, fire email.clicked webhook, record click event

	c.Redirect(http.StatusFound, originalURL)
}
