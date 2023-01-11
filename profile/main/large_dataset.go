package main

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
)

func main() {
	dir, err := os.MkdirTemp("", "track-test")
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(dir)

	track, err := core.NewTrack(&dir)
	if err != nil {
		panic(err.Error())
	}

	out.Print("Generating dataset\n")
	generateDataset(
		&track,
		util.Date(1990, 1, 1),
		2*time.Hour,
		10000,
	)

	f, err := os.Create("large_dataset.pprof")
	if err != nil {
		panic(err.Error())
	}

	out.Print("Profiling\n")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	fn, results, _ := track.AllRecords()
	go fn()
	for res := range results {
		_ = res.Record
	}
}

func generateDataset(t *core.Track, start time.Time, step time.Duration, records int) error {
	currTime := start
	for i := 0; i < records; i++ {
		rec := timedRecord(currTime, step/2, rand.Intn(3), rand.Intn(3))
		if err := t.SaveRecord(&rec, false); err != nil {
			return err
		}
		currTime = currTime.Add(step)
	}
	return nil
}

var noteLines = strings.Split(
	`Lorem ipsum dolor sit +amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut +labore et dolore magna aliqua.
A diam sollicitudin +tempor id eu nisl nunc mi ipsum.
Ullamcorper sit amet +risus +nullam eget felis eget.
Cursus euismod quis viverra nibh. Nunc sed blandit libero volutpat sed.
At augue eget arcu dictum varius +duis at consectetur lorem.
Velit euismod in pellentesque +massa placerat duis ultricies.
Risus nec +feugiat in fermentum posuere urna.
Id diam vel quam +elementum +pulvinar +etiam non quam.
Adipiscing bibendum est ultricies +integer quis auctor elit sed.
Sagittis orci a scelerisque purus +semper eget duis.
Elementum tempus egestas sed sed risus pretium.
Sit amet +mattis vulputate enim nulla aliquet porttitor lacus luctus.
Dignissim +sodales ut eu sem integer vitae justo eget.
Posuere urna nec tincidunt +praesent semper feugiat nibh sed pulvinar.
Dignissim cras tincidunt lobortis +feugiat vivamus at augue eget.
Viverra nam libero justo laoreet sit +amet cursus.
Venenatis +lectus +magna +fringilla +urna porttitor rhoncus dolor.
Aliquam id +diam maecenas ultricies.`, "\n")

func timedRecord(start time.Time, duration time.Duration, pauses int, lines int) core.Record {
	pause := make([]core.Pause, pauses, pauses)
	note := make([]string, lines, lines)

	if pauses > 0 {
		step := time.Duration(int(duration) / (pauses + 2))
		pauseStart := start.Add(step / 2)
		for i := 0; i < pauses; i++ {
			pause[i] = core.Pause{
				Start: pauseStart,
				End:   pauseStart.Add(step / 2),
				Note:  "Pause comment",
			}
			pauseStart = pauseStart.Add(step)
		}
	}

	for i := 0; i < lines; i++ {
		note[i] = noteLines[rand.Intn(len(noteLines))]
	}

	return core.Record{
		Project: "test",
		Start:   start,
		End:     start.Add(duration),
		Pause:   pause,
		Note:    strings.Join(note, "\n"),
		Tags:    map[string]string{},
	}
}
