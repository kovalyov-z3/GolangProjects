package computation

import (
	"fmt"
	"math"
)

func Solve(a float64, b float64, c float64) string {
	D := b*b - 4*(a*c)
	var x1 float64
	var x2 float64
	var answer string
	if D >= 0.0 {
		x1 = (b*(-1.0) + math.Sqrt(D)) / (2 * a)
		x2 = (b*(-1.0) - math.Sqrt(D)) / (2 * a)
		answer = "Первый корень = " + fmt.Sprint(x1) + " Второй корень = " + fmt.Sprint(x2)

	}
	if D == 0.0 {
		x1 = (b*(-1.0) + math.Sqrt(D)) / (2 * a)
		answer = "Единственный корень = " + fmt.Sprint(x1)
	}
	if D < 0.0 {
		answer = "Нет действительных корней"
	}
	return answer
}
