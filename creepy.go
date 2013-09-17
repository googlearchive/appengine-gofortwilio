// +build appengine

package fortwilio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"github.com/nf/twilio"
)

func Reshare(tw twilio.Context) {
	c := appengine.NewContext(tw.Request())
	from := tw.Value("From")
	var lf LastFortune
	if err := datastore.Get(c, datastore.NewKey(c, "LastFortune", from, 0, nil), &lf); err != nil {
		c.Errorf("error getting last fortune for %q: %v", from, err)
		return
	}
	c.Debugf("LastFortune(%q): %v", from, lf)
	client := urlfetch.Client(c)
	friend := strings.TrimSpace(tw.Value("Body"))
	data := url.Values{
		"From": []string{lf.Number},
		"To": []string{friend},
		"Url": []string{fmt.Sprintf("%s/%s", repeatUrl, from)},
		"Method": []string{"GET"},
	}
	c.Debugf("POST(%s): %q", callsApiUrl, data.Encode())
	req, err := http.NewRequest("POST", callsApiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		c.Errorf("error creating API request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(accountSid, authToken)
	resp, err := client.Do(req)
	if err != nil {
		c.Errorf("error posting outbound call: %v", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Errorf("error reading API response body: %v", err)
		return
	}
	c.Debugf("API response body: %q", string(body))
}

func Repeat(tw twilio.Context) {
	req := tw.Request()
	c := appengine.NewContext(req)
	parts := strings.Split(req.URL.Path, "/")
	from := parts[len(parts)-1]
	var lf LastFortune
	if err := datastore.Get(c, datastore.NewKey(c, "LastFortune", from, 0, nil), &lf); err != nil {
		c.Errorf("error getting last fortune for %q: %v", from, err)
		return
	}
	c.Debugf("LastFortune(%q): %v", from, lf)
	tw.Response(string(lf.Twiml))
}

func init() {
	twilio.Handle("/reshare", Reshare)
	twilio.Handle("/repeat/", Repeat)
}
