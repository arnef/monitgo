package utils

import (
	"math"
	"regexp"
	"strconv"
)

// ParseMegabyte converts given value to Megabyte value max 4 decimials
func ParseMegabyte(val string) (float64, error) {
	re := regexp.MustCompile(`(\d*\.?\d+)((G|M|T|k)?i?(B))`)
	result := re.FindStringSubmatch(val)

	value, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return -1, err
	}
	switch result[2] {
	case "GiB":
		value = value * 1073.74
	case "GB":
		value = value * 1000
	case "MiB":
		value = value * 1.04858
	case "kB":
		value = value * 0.001
	case "B":
		value = value * 1e-6
	}
	return math.Round(value*1000) / 1000, nil
}

// MustParseMegabyte parse val and panic on failure
func MustParseMegabyte(val string) float64 {
	res, err := ParseMegabyte(val)
	if err != nil {
		panic(err)
	}
	return res
}

func Round(value float64) float64 {
	return math.Round(value*100) / 100
}
