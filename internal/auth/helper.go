package auth

import (
	"net/http"
)

// SetAuthCookies sets the access and refresh tokens as secure cookies.
// accessToken: Short-lived (e.g., 15 mins)
// refreshToken: Long-lived (e.g., 7 days)
func (h *Handler) SetAuthCookies(w http.ResponseWriter, accessToken string, accessCookieExpiry int, refreshToken string, refreshCoookieExpiry int) {
	// 1. Configure the Access Token Cookie
	accessCookie := &http.Cookie{
		Name:     h.accessCookieName,
		Value:    accessToken,
		Path:     h.accessCookiePath,   // Cookie is valid for all paths
		MaxAge:   accessCookieExpiry,   // Expires in 15 minutes (in seconds)
		HttpOnly: true,                 // XSS protection: JS cannot read this
		Secure:   true,                 // Send only over HTTPS
		SameSite: http.SameSiteLaxMode, // CSRF protection
	}

	// 2. Configure the Refresh Token Cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",                  // Or set to specific refresh path like "/auth/refresh"
		MaxAge:   refreshCoookieExpiry, // Expires in 7 days
		HttpOnly: true,                 // XSS protection
		Secure:   true,                 // HTTPS only
		SameSite: http.SameSiteLaxMode,
	}

	// 3. Set the cookies on the ResponseWriter
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

func ClearAuthCookies(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Instant expiration
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	// Repeat for refresh_token...
}
