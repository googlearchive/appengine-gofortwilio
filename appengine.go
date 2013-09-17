// +build appengine

package fortwilio

import (
	"fmt"

	"appengine"
	"appengine/datastore"

	"github.com/nf/twilio"
)

var allVoices = []string{
	"man",
	"woman",
	"alice",
}

var enLangs = []string{
	"en", "en-gb",
}

type Fortwilio struct {
	Fortunes FortuneJar
	Langs    []string
	Voices   []string
}

func (ft Fortwilio) Say() string {
	f := ft.Fortunes.Get()
	v, l := choose(ft.Voices), choose(ft.Langs)
	twiml := fmt.Sprintf("<Say voice=%q language=%q>%v</Say>", v, l, f)
	twiml += fmt.Sprintf("<Sms>%v</Sms>", f)
	return twiml
}

var fts = map[string]Fortwilio{
	"fortunes":     {Voices: allVoices, Langs: enLangs, Fortunes: FortuneJar{}},
	"startrek":     {Voices: allVoices, Langs: enLangs, Fortunes: FortuneJar{}},
	"bofh-excuses": {Voices: allVoices, Langs: enLangs, Fortunes: FortuneJar{}},
	"proverbes":    {Voices: allVoices, Langs: []string{"fr-FR", "fr-CA"}, Fortunes: FortuneJar{}},
}

type LastFortune struct {
	Number string `datastore:",noindex"`
	Twiml  []byte `datastore:",noindex"`
}

func (ft Fortwilio) Handle(tw twilio.Context) {
	c := appengine.NewContext(tw.Request())
	twiml := ft.Say()
	c.Debugf(twiml)
	tw.Response(twiml)
	lf := &LastFortune{
		tw.Value("To"),
		[]byte(twiml),
	}
	if _, err := datastore.Put(c, datastore.NewKey(c, "LastFortune", tw.Value("From"), 0, nil), lf); err != nil {
		c.Errorf("error recording last fortune for %q: %v", tw.Value("From"), err)
		return
	}
}

func init() {
	for f, ft := range fts {
		ft.Fortunes.Load(f + ".u8")
		twilio.Handle("/"+f, ft.Handle)
	}
}
