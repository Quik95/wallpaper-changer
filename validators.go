package wallpaperchanger

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateArgs takes SearchConfig struct and validates arguments defined by user
func ValidateArgs(c SearchConfig) error {
	if err := validateCategory(c.Categories); err != nil {
		return err
	}
	if err := validateCategory(c.Purity); err != nil {
		return err
	}
	if err := validateResolution(c.Resolutions); err != nil {
		return err
	}
	if err := validateResolution(c.Ratios); err != nil {
		return err
	}
	if err := validateResolution(c.Atleast); err != nil {
		return err
	}
	if err := validateSorting(c.Sorting); err != nil {
		return err
	}
	if err := validateOrder(c.Order); err != nil {
		return err
	}
	if err := validateTimeRange(c.TopRange); err != nil {
		return err
	}
	if err := validatePages(c.Pages); err != nil {
		return err
	}
	if err := validateQuery(c.Query); err != nil {
		return err
	}

	if err := validatePairedOptions(c); err != nil {
		return err
	}

	return nil
}

func validatePairedOptions(c SearchConfig) error {
	purity := c.Purity
	key := c.APIKey

	if len(purity) == 3 && purity[2] == '1' && len(key) == 0 {
		return fmt.Errorf("When using purity setting NSFW providing api key is required")
	}

	sorting := c.Sorting
	topRange := c.TopRange

	if len(topRange) != 0 && sorting != "toplist" {
		return fmt.Errorf("When using top-range option sorting must be set to toplist")
	}

	return nil
}

func validateCategory(c string) error {
	if len(c) == 0 {
		return nil
	}

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
	if len(r) == 0 {
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
	if len(s) == 0 {
		return nil
	}

	valid := map[string]bool{
		"date_added": true,
		"relevance":  true,
		"random":     true,
		"views":      true,
		"favorites":  true,
		"toplist":    true,
		"hot":        true,
	}

	return validateMap(&valid, s, fmt.Sprintf("%s is not a valid sorting type", s))
}

func validateOrder(o string) error {
	if len(o) == 0 {
		return nil
	}

	valid := map[string]bool{
		"desc": true,
		"asc":  true,
	}

	return validateMap(&valid, o, fmt.Sprintf("%s is not a valid sorting order", o))
}

func validateTimeRange(r string) error {
	if len(r) == 0 {
		return nil
	}

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

func validateQuery(q string) error {
	parts := strings.Split(q, ",")

	for _, part := range parts {
		//check validity of an ID
		if strings.HasPrefix(part, "id:") {
			if _, err := strconv.Atoi(part[3:]); err != nil {
				return fmt.Errorf("%s is not a valid wallpaper ID", part[3:])
			}
		}

		// check filetype validity
		if strings.HasPrefix(part, "type:") {
			rest := part[5:]
			if rest != "png" && rest != "jpg" {
				return fmt.Errorf("%s is not a valid wallpaper extension. Use png or jpg", rest)
			}
		}
	}

	return nil
}
