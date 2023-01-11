package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
	"golang.org/x/exp/maps"
)

// TextRenderer renders a week or day schedule as colored text
type TextRenderer struct {
	Weekly        bool
	BlocksPerHour int
}

// Render renders the schedule
func (r *TextRenderer) Render(t *core.Track, reporter *core.Reporter, startDate time.Time) (string, error) {
	bph := r.BlocksPerHour

	spaceSym := []rune(t.Config.EmptyCell)[0]
	pauseSym := []rune(t.Config.PauseCell)[0]
	recordSym := []rune(t.Config.RecordCell)[0]

	projects := maps.Keys(reporter.Projects)
	sort.Strings(projects)
	indices := make(map[string]int, len(projects))
	symbols := make([]rune, len(projects)+1, len(projects)+1)
	colors := make([]color.Style256, len(projects)+1, len(projects)+1)
	symbols[0] = spaceSym
	colors[0] = *color.S256(15, 0)
	for i, p := range projects {
		indices[p] = i + 1
		symbols[i+1] = []rune(reporter.Projects[p].Symbol)[0]
		colors[i+1] = reporter.Projects[p].Render
	}

	numDays := 1
	if r.Weekly {
		numDays = 7
	}

	timeline := make([]int, 24*numDays*bph, 24*numDays*bph)
	paused := make([]bool, 24*numDays*bph, 24*numDays*bph)
	record := make([]int, 24*numDays*bph, 24*numDays*bph)

	now := time.Now()

	for recIdx, rec := range reporter.Records {
		startIdx, endIdx, ok := toIndexRange(rec.Start, rec.End, startDate, bph, numDays)
		if !ok {
			continue
		}
		index := indices[rec.Project]
		for i := startIdx; i <= endIdx; i++ {
			timeline[i] = index
			record[i] = recIdx + 1
		}
		for _, p := range rec.Pause {
			startIdx, endIdx, ok := toIndexRange(p.Start, p.End, startDate, bph, numDays)
			if !ok {
				continue
			}
			for i := startIdx; i <= endIdx; i++ {
				paused[i] = true
			}
		}
	}

	nowIdx := int(now.Sub(startDate).Hours() * float64(bph))

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "      |Day %s : %s/cell\n",
		startDate.Format(util.DateFormat),
		time.Duration(1e9*(int(time.Hour)/(bph*1e9))).String(),
	)

	fmt.Fprint(&sb, "      ")
	for weekday := 0; weekday < numDays; weekday++ {
		date := startDate.Add(time.Duration(weekday * 24 * int(time.Hour)))
		str := fmt.Sprintf(
			"%s %02d %s",
			date.Weekday().String()[:2],
			date.Day(),
			date.Month().String()[:3],
		)
		if len(str) > bph {
			fmt.Fprintf(&sb, "|%s", str[:bph])
		} else {
			fmt.Fprintf(&sb, "|%s%s", str, strings.Repeat(" ", bph-len(str)))
		}
	}
	fmt.Fprintln(&sb, "|")

	lastRecord := -1
	idxRecord := 0
	currNote := []rune{}
	currName := []rune{}
	for hour := 0; hour < 24; hour++ {
		fmt.Fprintf(&sb, "%02d:00 ", hour)
		for weekday := 0; weekday < numDays; weekday++ {
			s := (weekday*24 + hour) * bph
			fmt.Fprint(&sb, "|")
			for i := s; i < s+bph; i++ {
				rec := record[i]
				pr := timeline[i]
				pause := paused[i]

				if rec > 0 && rec != lastRecord {
					lastRecord = rec
					idxRecord = 0
					if pr == 0 {
						currNote = []rune{}
						currName = []rune{}
					} else {
						currNote = []rune(reporter.Records[rec-1].Note)
						currName = []rune(reporter.Records[rec-1].Project)
					}
				} else {
					if !pause {
						idxRecord++
					}
				}

				sym := symbols[pr]
				col := colors[pr]
				if pause {
					sym = pauseSym
				}
				if !r.Weekly && !pause && pr > 0 {
					nameLen := len(currName)
					noteLen := len(currNote)
					if idxRecord == 0 {
						sym = ' '
					} else if idxRecord-1 < nameLen {
						sym = currName[idxRecord-1]
					} else if idxRecord-1 == nameLen {
						sym = ':'
					} else if idxRecord-1 == nameLen+1 {
						sym = ' '
					} else if idxRecord-3-nameLen < noteLen {
						sym = currNote[idxRecord-3-nameLen]
						if sym == '\n' || sym == '\r' {
							sym = ' '
						}
					} else if idxRecord-3-nameLen == noteLen {
						sym = ' '
					} else {
						sym = recordSym
					}
				}
				if i == nowIdx {
					sym = '@'
				}
				fmt.Fprint(&sb, col.Sprintf("%c", sym))
			}
		}
		fmt.Fprintln(&sb, "|")
	}

	totalWidth := 7 + numDays*(bph+1)
	lineWidth := 0

	line1 := ""
	line2 := ""
	for i, p := range projects {
		col := colors[i+1]
		width := utf8.RuneCountInString(p)
		if width < 3 {
			width = 3
		}
		if lineWidth > 0 && lineWidth+width+4 > totalWidth {
			lineWidth = 0
			fmt.Fprintln(&sb, line1)
			fmt.Fprintln(&sb, line2)
			line1 = ""
			line2 = ""
		}

		line1 += col.Sprintf(" %c:%3s ", symbols[indices[p]], p)
		line2 += col.Sprintf(" %*s ", width+2, util.FormatDuration(reporter.TotalTime[p], false))
		lineWidth += width + 4
	}
	if len(line1) > 0 {
		fmt.Fprintln(&sb, line1)
		fmt.Fprintln(&sb, line2)
	}

	return sb.String(), nil
}

func toIndexRange(start, end, startDate time.Time, bph int, days int) (int, int, bool) {
	if end.IsZero() {
		end = time.Now()
	}
	if end.Before(startDate) {
		return -1, -1, false
	}

	startIdx := int(start.Sub(startDate).Hours() * float64(bph))
	endIdx := int(end.Sub(startDate).Hours()*float64(bph)) - 1
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx >= bph*24*days {
		endIdx = bph*24*days - 1
	}
	return startIdx, endIdx, true
}
