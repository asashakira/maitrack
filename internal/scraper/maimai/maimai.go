package maimai

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	URL = "https://maimaidx.jp/maimai-mobile"
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
	res, err := m.HTTPClient.Get(URL)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var csrfToken string
	doc.Find(`.black input[type="hidden"]`).Each(func(i int, s *goquery.Selection) {
		csrfToken = s.AttrOr("value", "CSRF token not found")
	})

	// login
	values := make(url.Values)
	values.Set("segaId", segaId)
	values.Set("password", password)
	values.Set("token", csrfToken)
	_, err = m.HTTPClient.PostForm(URL+"/submit", values)
	if err != nil {
		return err
	}

	// redirect to set cookies
	_, err = m.HTTPClient.Get(URL + "/aimeList/submit/?idx=0")
	if err != nil {
		return err
	}

	return nil
}
