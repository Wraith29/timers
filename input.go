package main

import (
	"fmt"
	"strconv"
	"unicode"

	"fyne.io/fyne/v2/widget"
)

func timeInputValidator(s string) error {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("invalid time %c", c)
		}
	}

	return nil
}

type TimeInput struct {
	widget.Entry
}

func NewTimeInput() *TimeInput {
	entry := &TimeInput{}
	entry.ExtendBaseWidget(entry)
	entry.SetText("0")
	entry.Validator = timeInputValidator

	return entry
}

func (t *TimeInput) Get() int {
	num, err := strconv.Atoi(t.Text)
	if err != nil {
		fmt.Printf("ERR: Failed to convert %s to integer\n", t.Text)
		panic(err)
	}

	return num
}
