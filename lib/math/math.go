package math

import (
	"math"
	"math/rand"
	"time"
)

// Standard math package for most common mathematical operations

type i interface{}

var (
	Includes = map[string]interface{}{}
)

func init() {
	rand.Seed(time.Hour.Milliseconds())

	noErr := map[string]func(float64) float64{
		/* func sin(num float64) float64 */
		"sin": math.Sin,
		/* func asin(num float64) float64 */
		"asin": math.Asin,
		/* func cos(num float64) float64 */
		"cos": math.Cos,
		/* func acos(num float64) float64 */
		"acos": math.Acos,
		/* func tan(num float64) float64 */
		"tan": math.Tan,
		/* func atan(num float64) float64 */
		"atan": math.Atan,
		/* func floor(num float64) float64 */
		"floor": math.Floor,
		/* func ceil(num float64) int */
		"ceil": math.Ceil,
		/* func abs(num float64) int */
		"abs": math.Abs,
		/* func ln(num float64) int */
		"ln": math.Log,
		/* func log10(num float64) float64 */
		"log10": math.Log10,
		/* func sqrt(num float64) float64 */
		"sqrt": math.Sqrt,
	}

	Includes = map[string]interface{}{
		"max": max,
		"min": min,
		"deg": deg,
		"rad": rad,
		"random": random,
	}

	for name, f := range noErr {
		fun := f // Store function in closure to keep definition for current name
		Includes[name] = func(num float64) (val i, err error) {
			return fun(num), err
		}
	}
}

/* func max(a float64, b float64) float64 */
func max(a float64, b float64) (val i, err error) {
	return math.Max(a, b), err
}

/* func min(a float64, b float64) float64 */
func min(a float64, b float64) (val i, err error) {
	return math.Min(a, b), err
}

/*
	Converts degrees to radians
	func rad(deg float64) float64
*/
func rad(num float64) (val i, err error) {
	return (math.Pi * 2 * num) / 360, err
}

/*
	Converts radians to degrees
	func deg(rad float64) float64
*/
func deg(num float64) (val i, err error) {
	return (num * 360) / (2 * math.Pi), err
}

/*
	Gets random number bewteen 0 and 1
	func rand() float64
*/
func random() (val i, err error) {
	return rand.Float64(), err
}