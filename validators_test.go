package wallpaperchanger

import (
	"flag"
	"fmt"
	"testing"

	"github.com/urfave/cli/v2"
)

func getEmptyConfig() (*cli.Context, *flag.FlagSet) {
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)
	command := cli.Command{Name: "mycommand"}
	ctx.Command = &command

	return ctx, set
}

// returns true when ans doesn't match wantsError
func checkResult(ans error, wantsError bool) bool {
	isError := false
	if ans != nil {
		isError = true
	}
	if isError != wantsError {
		return true
	}

	return false
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		category, value string
		wantsError      bool
	}{
		// categories and purity settings
		{"categories", "111-", true},
		{"categories", "", false},
		{"categories", "000", true},
		{"categories", "010", false},
		{"categories", "101", false},

		// don't test with NSFW setting 1 to avoid checking api key for now
		{"purity", "111-", true},
		{"purity", "", false},
		{"purity", "000", true},
		{"purity", "010", false},
		{"purity", "100", false},

		//resolutions and ratios
		{"resolutions", "", false},
		{"resolutions", "1920x1080", false},
		{"resolutions", "1920x1080,1280x720", false},
		{"resolutions", "19201080", true},
		{"resolutions", "1920x1080,", true},
		{"resolutions", "1920xx1080,", true},
		{"resolutions", "xx1080,", true},

		{"ratios", "", false},
		{"ratios", "16x9", false},
		{"ratios", "16x9,16x10", false},
		{"ratios", "169", true},
		{"ratios", "16x9,", true},
		{"ratios", "16xx9", true},
		{"ratios", "xx9", true},

		// sorting types
		{"sorting", "", false},
		{"sorting", "date_added", false},
		{"sorting", "relevance", false},
		{"sorting", "random", false},
		{"sorting", "views", false},
		{"sorting", "favorites", false},
		{"sorting", "toplist", false},
		{"sorting", "aaa", true},
		{"sorting", "AAAAAA", true},
		{"sorting", "top-list", true},
		{"sorting", "favourites", true},

		//order types
		{"order", "", false},
		{"order", "asc", false},
		{"order", "desc", false},
		{"order", "ascending", true},
		{"order", "descending", true},

		//query
		{"query", "", false},
		{"query", "cyberpunk", false},
		{"query", "-cyberpunk", false},
		{"query", "+cyberpunk", false},
		{"query", "-cyberpunk +starwars", false},
		{"query", "@Quik95", false},
		{"query", "id:4445", false},
		{"query", "id:1", false},
		{"query", "id:a", true},
		{"query", "type:png", false},
		{"query", "type:jpg", false},
		{"query", "type:jpeg", true},
		{"query", "like:slfk24", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Category: %s, Value: %s", tt.category, tt.value)

		ctx, set := getEmptyConfig()
		set.String(tt.category, tt.value, "")

		//workaround for empty pages test
		set.Int("pages", 1, "")

		t.Run(testname, func(t *testing.T) {
			ans := ValidateArgs(ctx)
			if checkResult(ans, tt.wantsError) {
				t.Errorf("got: %v, want error: %t", ans, tt.wantsError)
			}
		})
	}
}

func TestValidatePairedOptions(t *testing.T) {
	tests := []struct {
		keys, values []string
		wantsError   bool
	}{
		// apikey and purity
		{[]string{"purity", "api-key"}, []string{"111", "asdasda"}, false},
		{[]string{"purity", "api-key"}, []string{"110", "asdasda"}, false},
		{[]string{"purity", "api-key"}, []string{"111", ""}, true},

		//time range and sorting toplist
		{[]string{"sorting", "top-range"}, []string{"toplist", "1w"}, false},
		{[]string{"sorting", "top-range"}, []string{"random", "1w"}, true},
		{[]string{"sorting", "top-range"}, []string{"toplist", ""}, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Testing paired options: %#v, with values: %#v", tt.keys, tt.values)

		ctx, set := getEmptyConfig()

		for i, key := range tt.keys {
			set.String(key, tt.values[i], "")
		}

		t.Run(testname, func(t *testing.T) {
			ans := validatePairedOptions(ctx)
			if checkResult(ans, tt.wantsError) {
				t.Errorf("got: %v, want error: %t", ans, tt.wantsError)
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
		testname := fmt.Sprintf("Testing timerange: %s", tt.data)
		t.Run(testname, func(t *testing.T) {
			ans := validateTimeRange(tt.data)
			if checkResult(ans, tt.wantsError) {
				t.Fatalf("Got %v, wanted error: %t", ans, tt.wantsError)
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
		{10, false},
		{-10, true},
		{0, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Testing page number: %d", tt.data)
		t.Run(testname, func(t *testing.T) {
			ans := validatePages(tt.data)
			if checkResult(ans, tt.wantsError) {
				t.Fatalf("Got %v, wanted error: %t", ans, tt.wantsError)
			}
		})
	}
}
