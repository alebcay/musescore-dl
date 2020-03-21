package cmd

import (
	"time"

	"github.com/alebcay/musescore-dl/internal/pkg"
	"github.com/briandowns/spinner"
)

func DownloadMIDI(url string, dest string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	s.Suffix = " Getting score information"
	id, secret := msdl.GetScoreIDSecret(url)
	if id == "" || secret == "" {
		panic("bad score ID/secret")
	}

	if dest == "" {
		dest = msdl.GetScoreTitle(url) + ".mid"
	}

	s.Suffix = " Downloading score MIDI"
	err := msdl.FetchMIDI(id, secret, dest)
	if err != nil {
		panic(err)
	}
}
