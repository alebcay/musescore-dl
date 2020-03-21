package msdl

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchMIDI(id string, secret string, dest string) error {
	x := id[len(id)-1:]
	y := id[len(id)-2 : len(id)-1]
	z := id[len(id)-3 : len(id)-2]

	midi_url := fmt.Sprintf("https://musescore.com/static/musescore/scoredata/gen/%s/%s/%s/%s/%s/score.mid", x, y, z, id, secret)
	resp, err := http.Get(midi_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	midi_file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer midi_file.Close()

	_, err = io.Copy(midi_file, resp.Body)

	return err
}
