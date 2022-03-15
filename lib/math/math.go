package math

import (
	"math"
	"math/rand"
	"time"
)

// Standard math package for most common mathematical operations

type i interface{}

func init() {
	rand.Seed(time.Hour.Milliseconds())
}

/* func sin(num float64) float64 */
func Sin(num float64) (val i, err error) {
	return math.Sin(num), err
}

/* func cos(num float64) float64 */
func Cos(num float64) (val i, err error) {
	return math.Cos(num), err
}

/* func asin(num float64) float64 */
func Asin(num float64) (val i, err error) {
	return math.Asin(num), err
}

/* func acos(num float64) float64 */
func Acos(num float64) (val i, err error) {
	return math.Acos(num), err
}

/* func tan(num float64) float64 */
func Tan(num float64) (val i, err error) {
	return math.Tan(num), err
}

/* func atan(num float64) float64 */
func Atan(num float64) (val i, err error) {
	return math.Atan(num), err
}

/* func floor(num float64) float64 */
func Floor(num float64) (val i, err error) {
	return math.Floor(num), err
}

/* func ceil(num float64) float64 */
func Ceil(num float64) (val i, err error) {
	return math.Ceil(num), err
}

/* func abs(num float64) float64 */
func Abs(num float64) (val i, err error) {
	return math.Abs(num), err
}

/* func ln(num float64) float64 */
func Ln(num float64) (val i, err error) {
	return math.Log(num), err
}

/* func log10(num float64) float64 */
func Log10(num float64) (val i, err error) {
	return math.Log10(num), err
}

/* func sqrt(num float64) float64 */
func Sqrt(num float64) (val i, err error) {
	return math.Sqrt(num), err
}

/* func max(a float64, b float64) float64 */
func Max(a float64, b float64) (val i, err error) {
	return math.Max(a, b), err
}

/* func min(a float64, b float64) float64 */
func Min(a float64, b float64) (val i, err error) {
	return math.Min(a, b), err
}

/*
	Converts degrees to radians
	func rad(deg float64) float64
*/
func Rad(num float64) (val i, err error) {
	return (math.Pi * 2 * num) / 360, err
}

/*
	Converts radians to degrees
	func deg(rad float64) float64
*/
func Deg(num float64) (val i, err error) {
	return (num * 360) / (2 * math.Pi), err
}

/*
	Gets random number bewteen 0 and 1
	func random() float64
*/
func Random() (val i, err error) {
	return rand.Float64(), err
}
