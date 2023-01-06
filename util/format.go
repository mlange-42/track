package util

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mlange-42/track/tree"
)

const (
	// DateFormat for date formatting
	DateFormat = "2006-01-02"
	// TimeFormat for time formatting
	TimeFormat = "15:04"
	// FileTimeFormat for file name time formatting
	FileTimeFormat = "15-04"
	// FileDateTimeFormat for date and time from paths
	FileDateTimeFormat = "2006-01-02 15-04"
	// DateTimeFormat for date and time formatting
	DateTimeFormat = "2006-01-02 15:04"
	// NoTimeString string representation for zero end time
	NoTimeString = " now "
)

const (
	// PrevDayPrefix is a prefix for a time on the previous day
	PrevDayPrefix = "<"
	// NextDaySuffix is a suffix for a time on the next day
	NextDaySuffix = ">"
)

// BlockRunes are utf8 8th blocks from empty to full
var BlockRunes = [9]rune{'·', 9601, 9602, 9603, 9604, 9605, 9606, 9607, 9608}

// FloatToBlock returns utf8 8th blocks for values between 0 and 1
func FloatToBlock(value float64, space *rune) rune {
	idx := int(value * float64(len(BlockRunes)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(BlockRunes) {
		return BlockRunes[len(BlockRunes)-1]
	}
	if idx == 0 && space != nil {
		return *space
	}
	return BlockRunes[idx]
}

// FormatDuration formats a duration
func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(d.Hours()), int(d.Minutes())%60)
}

// FormatTimeWithOffset formats a time with day offset indicators
func FormatTimeWithOffset(t time.Time, reference time.Time) string {
	if t.IsZero() {
		return "?"
	}
	timeStr := t.Format(TimeFormat)
	if ToDate(t).After(reference) {
		timeStr = timeStr + NextDaySuffix
		return timeStr
	}
	if ToDate(reference).After(t) {
		timeStr = PrevDayPrefix + timeStr
		return timeStr
	}
	return timeStr
}

// Format formats a string with named placeholders.
//
// Example:
// s := Format("foo {name} bar", map[string]string{"name": "baz"})
func Format(str string, repl map[string]string) string {
	format := "{%s}"
	for k, v := range repl {
		str = strings.ReplaceAll(str, fmt.Sprintf(format, k), v)
	}
	return str
}

// TreeFormatter formats trees
type TreeFormatter[T any] struct {
	NameFunc     func(t *tree.MapNode[T], indent int) string
	Indent       int
	prefixNone   string
	prefixEmpty  string
	prefixNormal string
	prefixLast   string
}

// NewTreeFormatter creates a new TreeFormatter
func NewTreeFormatter[T any](
	nameFunc func(t *tree.MapNode[T], indent int) string,
	indent int,
) TreeFormatter[T] {
	return TreeFormatter[T]{
		NameFunc:     nameFunc,
		Indent:       indent,
		prefixNone:   strings.Repeat(" ", indent),
		prefixEmpty:  "│" + strings.Repeat(" ", indent-1),
		prefixNormal: "├" + strings.Repeat("─", indent-1),
		prefixLast:   "└" + strings.Repeat("─", indent-1),
	}
}

// FormatTree formats a tree
func (f *TreeFormatter[T]) FormatTree(t *tree.MapTree[T]) string {
	sb := strings.Builder{}
	f.formatTree(&sb, t.Root, 0, false, "")
	return sb.String()
}

func (f *TreeFormatter[T]) formatTree(sb *strings.Builder, t *tree.MapNode[T], depth int, last bool, prefix string) {
	pref := prefix
	if depth > 0 {
		pref = prefix + f.createPrefix(last)
	}
	fmt.Fprint(sb, pref)
	fmt.Fprintf(sb, "%s", f.NameFunc(t, utf8.RuneCountInString(pref)))
	fmt.Fprint(sb, "\n")

	if depth > 0 {
		pref = prefix + f.createPrefixEmpty(last)
	}

	names := make([]string, 0, len(t.Children))
	for name := range t.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	for i, name := range names {
		last := i == len(names)-1
		f.formatTree(sb, t.Children[name], depth+1, last, pref)
	}
}

func (f *TreeFormatter[T]) createPrefix(last bool) string {
	if last {
		return f.prefixLast
	}
	return f.prefixNormal
}

func (f *TreeFormatter[T]) createPrefixEmpty(last bool) string {
	if last {
		return f.prefixNone
	}
	return f.prefixEmpty
}
