package tvwsclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type TVHttpClient struct {
	baseURL     string
	deviceToken string
	sessionID   string
	sessionSign string
	httpClient  *http.Client
}

func NewTVHttpClient(baseURL string, deviceToken string, sessionID string, sessionSign string) *TVHttpClient {
	return &TVHttpClient{
		baseURL:     baseURL,
		deviceToken: deviceToken,
		sessionID:   sessionID,
		sessionSign: sessionSign,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *TVHttpClient) GetQuoteToken() (string, error) {
	// Create request
	req, err := http.NewRequest("POST", c.baseURL+"/quote_token/", nil)
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://www.tradingview.com")
	req.Header.Set("Referer", "https://www.tradingview.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Language", "en")

	// Set cookies with the complete format
	req.Header.Set("Cookie", "cookiePrivacyPreferenceBannerProduction=notApplicable; "+
		"cookiesSettings={\"analytics\":true,\"advertising\":true}; "+
		"device_t="+c.deviceToken+"; "+
		"sessionid="+c.sessionID+"; "+
		"sessionid_sign="+c.sessionSign)

	// Add form data
	req.URL.RawQuery = "grabSession=true"

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Remove quotes from the response if present
	token := strings.Trim(string(body), "\"")

	return token, nil
}
