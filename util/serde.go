package util

import (
	"time"

	"gopkg.in/yaml.v3"
)

// Time is for somplified formatting in serialization/deserialization
type Time time.Time

// MarshalYAML converts a Time to YAML bytes
func (t Time) MarshalYAML() (interface{}, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return "nil", nil
	}
	return tm.Format(YAMLTimeFormat), nil
}

// UnmarshalYAML converts YAML bytes to a Time
func (t *Time) UnmarshalYAML(value *yaml.Node) (err error) {
	if value.Value == "nil" {
		*t = Time{}
		return
	}
	now, err := time.ParseInLocation(YAMLTimeFormat, value.Value, time.Local)
	*t = Time(now)
	return
}
