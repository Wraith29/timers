package main

import (
	"fmt"
	"unicode"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func validateTimeInput(s string) error {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("invalid time %c", c)
		}
	}

	return nil
}

type TimeEntry struct {
	widget.Entry
}

func NewTimeEntry(ds binding.String) *TimeEntry {
	entry := &TimeEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Bind(ds)
	entry.SetText("0")
	entry.Validator = validateTimeInput

	return entry
}
