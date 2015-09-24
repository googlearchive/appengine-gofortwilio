package fortwilio

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	sep = "%\n"
)

type FortuneJar []string

func (fj *FortuneJar) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open fortune file %q: %v", filename, err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read from fortune file %q: %v", filename, err)
	}
	*fj = strings.Split(strings.TrimRight(string(bs), sep), sep)
	return nil
}

func (fj FortuneJar) Get() string {
	return choose(fj)
}

func choose(s []string) string {
	return s[rand.Intn(len(s))]
}

func init() {
	rand.Seed(time.Now().Unix())
}
