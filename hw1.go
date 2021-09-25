package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2*z)
		fmt.Println("Now z is", z)
	}
	return z
}

func main() {
	fmt.Println("The answer is:", Sqrt(2))
}