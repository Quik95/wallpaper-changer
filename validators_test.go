package wallpaperchanger

import (
	"fmt"
	"testing"
)

func TestValidateCategory(t *testing.T) {
	var tests = []struct {
		data string
		//is error
		wantsError bool
	}{
		{"110", false},
		{"001", false},
		{"elo", true},
		{"1000", true},
		{"0", true},
		{"090", true},
	}

	for _, tt := range tests {
		testname := tt.data
		t.Run(testname, func(t *testing.T) {
			ans := validateCategory(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}

func TestValidateResolution(t *testing.T) {
	var tests = []struct {
		data string
		//is error
		wantsError bool
	}{
		{"1920x1080", false},
		{"16x9", false},
		{"1920x1080,1276x726", false},
		{"4x3", false},
		{"4", true},
		{"x", true},
		{",", true},
	}

	for _, tt := range tests {
		testname := tt.data
		t.Run(testname, func(t *testing.T) {
			ans := validateResolution(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}

func TestValidateSorting(t *testing.T) {
	var tests = []struct {
		data string
		//is error
		wantsError bool
	}{
		{"date_added", false},
		{"relevance", false},
		{"random", false},
		{"views", false},
		{"favorites", false},
		{"toplist", false},
		{"asasa", true},
		{"topList", true},
		{"dateadded", true},
		{"110", true},
	}

	for _, tt := range tests {
		testname := tt.data
		t.Run(testname, func(t *testing.T) {
			ans := validateSorting(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}

func TestValidateTimeRange(t *testing.T) {
	var tests = []struct {
		data string
		//is error
		wantsError bool
	}{
		{"1d", false},
		{"3d", false},
		{"1w", false},
		{"1M", false},
		{"3M", false},
		{"6M", false},
		{"1y", false},
		{"1", true},
		{"3m", true},
	}

	for _, tt := range tests {
		testname := tt.data
		t.Run(testname, func(t *testing.T) {
			ans := validateTimeRange(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}

func TestValidateOrder(t *testing.T) {
	var tests = []struct {
		data string
		//is error
		wantsError bool
	}{
		{"desc", false},
		{"asc", false},
		{"descending", true},
		{"ascending", true},
		{"as", true},
	}

	for _, tt := range tests {
		testname := tt.data
		t.Run(testname, func(t *testing.T) {
			ans := validateOrder(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}

func TestValidatePages(t *testing.T) {
	var tests = []struct {
		data int
		//is error
		wantsError bool
	}{
		{1, false},
		{10, false},
		{-1, true},
		{0, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprint(tt.data)
		t.Run(testname, func(t *testing.T) {
			ans := validatePages(tt.data)
			if ans == nil && tt.wantsError == true {
				t.Errorf("Got <nil>, expected error")
			}
			if ans != nil && tt.wantsError == false {
				t.Errorf("Got `%v`, expected <nil>", ans)
			}
		})
	}
}
