package maimaiclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL          = "https://maimaidx.jp/maimai-mobile"
	LoginEndpoint    = "/submit"
	AimeListEndpoint = "/aimeList/submit/?idx=0"
	ErrorEndpoint    = "/error/"
)

type Client struct {
	HTTPClient *http.Client
}

func New() *Client {
	cookiejar, _ := cookiejar.New(nil)

	c := &Client{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.HTTPClient = &http.Client{
		Transport:     tr,
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           cookiejar,
	}

	return c
}

func (m *Client) Login(segaId, password string) error {
	// Fetch the login page to get the CSRF token
	req, err := http.NewRequest("GET", BaseURL, nil)
	if err != nil {
		return err
	}

	// Add headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", BaseURL)

	res, err := m.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	csrfToken, extractCSRFTokenErr := extractCSRFToken(res)
	if extractCSRFTokenErr != nil {
		return extractCSRFTokenErr
	}

	// submit login form
	values := url.Values{
		"segaId":   {segaId},
		"password": {password},
		"token":    {csrfToken},
	}
	formReq, err := http.NewRequest("POST", BaseURL+LoginEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	// Add headers to mimic a browser
	formReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	formReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	formReq.Header.Set("Accept", "application/x-www-form-urlencoded")
	formReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formReq.Header.Set("Referer", BaseURL)

	loginRes, err := m.HTTPClient.Do(formReq)
	if err != nil {
		return err
	}
	defer loginRes.Body.Close()

	// check if login was successful
	if loginRes.Request.URL.String() == BaseURL+ErrorEndpoint {
		return fmt.Errorf("login failed: incorrect Sega ID or password")
	}

	// redirect to set cookies
	aimeReq, err := http.NewRequest("GET", BaseURL+AimeListEndpoint, nil)
	if err != nil {
		return err
	}

	// Add headers to mimic a browser
	aimeReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	aimeReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	aimeReq.Header.Set("Referer", BaseURL+LoginEndpoint)

	aimeRes, err := m.HTTPClient.Do(aimeReq)
	if err != nil {
		return err
	}
	defer aimeRes.Body.Close()

	return nil
}

// extractCSRFToken parses the CSRF token from the HTML document.
func extractCSRFToken(res *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	csrfToken := doc.Find(`.black input[type="hidden"]`).AttrOr("value", "")
	if csrfToken == "" {
		return "", errors.New("CSRF token not found")
	}

	return csrfToken, nil
}
