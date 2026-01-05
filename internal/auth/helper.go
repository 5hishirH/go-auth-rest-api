package auth

import (
	"net/http"
)

func GenerateCookieResponse(w http.ResponseWriter, cookieName string, cookiePath string, token string, expiry int, isSecure bool) {
	refreshCookie := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     cookiePath, // Or set to specific refresh path like "/auth/refresh"
		MaxAge:   expiry,     // Expires in 7 days
		HttpOnly: true,       // XSS protection
		Secure:   isSecure,   // HTTPS only
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, refreshCookie)
}

func GenerateClearCookieResponse(w http.ResponseWriter, cookieName string, cookiePath string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     cookiePath,
		MaxAge:   -1, // Instant expiration
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
