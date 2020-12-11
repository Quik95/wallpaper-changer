package wallpaperchanger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// ValidateArgs takes cli.Context struct and validates arguments defined by user
func ValidateArgs(c *cli.Context) error {
	if err := validateCategory(c.String("categories")); err != nil {
		return err
	}
	if err := validateCategory(c.String("purity")); err != nil {
		return err
	}
	if err := validateResolution(c.String("resolutions")); err != nil {
		return err
	}
	if err := validateResolution(c.String("ratios")); err != nil {
		return err
	}
	if err := validateResolution(c.String("atleast")); err != nil {
		return err
	}
	if err := validateSorting(c.String("sorting")); err != nil {
		return err
	}
	if err := validateOrder(c.String("order")); err != nil {
		return err
	}
	if err := validateTimeRange(c.String("top-range")); err != nil {
		return err
	}
	if err := validatePages(c.Int("pages")); err != nil {
		return err
	}

	if err := validatePairedOptions(c); err != nil {
		return err
	}

	return nil
}

func validatePairedOptions(c *cli.Context) error {
	purity := c.String("purity")
	key := c.String("api-key")

	if purity[2] == '1' && len(key) == 0 {
		return fmt.Errorf("When using purity setting NSFW providing api key is required")
	}

	sorting := c.String("sorting")
	topRange := c.String("top-range")

	if len(topRange) != 0 && sorting != "toplist" {
		return fmt.Errorf("When using top-range option sorting must be set to toplist")
	}

	return nil
}

func validateCategory(c string) error {
	if len(c) != 3 {
		return fmt.Errorf("%s is not a valid category", c)
	}
	possible := map[string]bool{
		"100": true,
		"110": true,
		"111": true,
		"011": true,
		"001": true,
		"101": true,
		"010": true,
	}
	if !possible[c] {
		return fmt.Errorf("%s is not a valid category", c)
	}
	return nil
}

func validateResolution(r string) error {
	if r == "" {
		return nil
	}

	parts := strings.Split(strings.ReplaceAll(r, ",", "x"), "x")
	if len(parts)%2 != 0 {
		return fmt.Errorf("%s is not a valid resolution", r)
	}

	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return fmt.Errorf("%s is not a valid number", p)
		}
	}

	return nil
}

func validateMap(m *map[string]bool, v string, message string) error {
	if !(*m)[v] {
		return fmt.Errorf(message)
	}

	return nil
}

func validateSorting(s string) error {
	valid := map[string]bool{
		"date_added": true,
		"relevance":  true,
		"random":     true,
		"views":      true,
		"favorites":  true,
		"toplist":    true,
	}

	return validateMap(&valid, s, fmt.Sprintf("%s is not a valid sorting type", s))
}

func validateOrder(o string) error {
	valid := map[string]bool{
		"desc": true,
		"asc":  true,
	}

	return validateMap(&valid, o, fmt.Sprintf("%s is not a valid sorting order", o))
}

func validateTimeRange(r string) error {
	valid := map[string]bool{
		"1d": true,
		"3d": true,
		"1w": true,
		"1M": true,
		"3M": true,
		"6M": true,
		"1y": true,
		"":   true,
	}

	return validateMap(&valid, r, fmt.Sprintf("%s is not a valid time range", r))
}

func validatePages(p int) error {
	if p <= 0 {
		return fmt.Errorf("Page number cannot be negative or zero")
	}

	return nil
}
