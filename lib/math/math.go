package math

import "math"

// Standard math package for most common mathematical operations

type i interface{}

var (
	Includes = map[string]interface{}{}
)

func init() {
	noErr := map[string]func(float64) float64{
		"sin":   math.Sin,
		"asin":  math.Asin,
		"cos":   math.Cos,
		"acos":  math.Acos,
		"tan":   math.Tan,
		"atan":  math.Atan,
		"floor": math.Floor,
		"ceil":  math.Ceil,
		"abs":   math.Abs,
		"ln":    math.Log,
		"log10": math.Log10,
		"sqrt":  math.Sqrt,
	}

	Includes = map[string]interface{}{
		// Max of a and b
		"max": max,
		// Min of a and b
		"min": min,
		// Converts radians to degrees
		"deg": deg,
		// Converts degrees to radians
		"rad": rad,
	}

	for name, f := range noErr {
		fun := f // Store function in closure to keep definition for current name
		Includes[name] = func(num float64) (val i, err error) {
			return fun(num), err
		}
	}
}

func max(a float64, b float64) (val i, err error) {
	return math.Max(a, b), err
}

func min(a float64, b float64) (val i, err error) {
	return math.Min(a, b), err
}

func rad(num float64) (val i, err error) {
	return (math.Pi * 2 * num) / 360, err
}

func deg(num float64) (val i, err error) {
	return (num * 360) / (2 * math.Pi), err
}
