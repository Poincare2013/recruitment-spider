package main

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	f, err := ioutil.ReadFile("./config/liepin.json")
	if err != nil {
		log.Fatal("read fail", err)
	}
	var cookies []MySetCookieParams
	err = json.Unmarshal(f, &cookies)
	if err != nil {
		log.Fatal("unmarshal fail", err)
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	go chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			if ev.Type == "XHR" {
				resp := ev.Response
				if resp.URL != "" {
					log.Println(resp.URL)
				}
			}
		}
	})

	if err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.ActionFunc(func(cxt context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => false, });").Do(cxt)
			if err != nil {
				return err
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, cookie := range cookies {
				_, err := genSetCookieParamsFromCookie(cookie).Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Navigate("https://lpt.liepin.com/cvsearch/showcondition/"),
	); err != nil {
		log.Fatal(err)
	}
	sync := make(chan struct{})
	<-sync
}

type MySetCookieParams struct {
	Name     string              `json:"name"`               // Cookie name.
	Value    string              `json:"value"`              // Cookie value.
	URL      string              `json:"url,omitempty"`      // The request-URI to associate with the setting of the cookie. This value can affect the default domain and path values of the created cookie.
	Domain   string              `json:"domain,omitempty"`   // Cookie domain.
	Path     string              `json:"path,omitempty"`     // Cookie path.
	Secure   bool                `json:"secure,omitempty"`   // True if cookie is secure.
	HTTPOnly bool                `json:"httpOnly,omitempty"` // True if cookie is http-only.
	Expires  *cdp.TimeSinceEpoch `json:"expires,omitempty"`  // Cookie expiration date, session cookie if not set
}

func genSetCookieParamsFromCookie(cookie MySetCookieParams) *network.SetCookieParams {
	return &network.SetCookieParams{
		Name:     cookie.Name,
		Value:    cookie.Value,
		URL:      cookie.URL,
		Domain:   cookie.Domain,
		Path:     cookie.Path,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HTTPOnly,
		Expires:  cookie.Expires,
	}
}
