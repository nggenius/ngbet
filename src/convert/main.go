package main

import (
	"fmt"
	"math"
)

var (
	Odds = [][15]float64{
		[15]float64{2.5, 2.55, 2.6, 2.65, 2.7, 2.75, 2.8, 2.85, 2.9, 2.95, 3, 3.05, 3.1, 3.15, 3.2},
		[15]float64{2, 2.03, 2.07, 2.1, 2.13, 2.17, 2.2, 2.23, 2.27, 2.3, 2.33, 2.37, 2.4, 2.43, 2.46},
		[15]float64{1.75, 1.78, 1.8, 1.83, 1.85, 1.88, 1.9, 1.93, 1.95, 1.98, 2, 2.03, 2.05, 2.08, 2.1},
		[15]float64{1.5, 1.52, 1.53, 1.55, 1.57, 1.58, 1.6, 1.62, 1.63, 1.65, 1.67, 1.68, 1.7, 1.72, 1.73},
		[15]float64{1.38, 1.39, 1.4, 1.41, 1.43, 1.44, 1.45, 1.46, 1.48, 1.49, 1.5, 1.51, 1.53, 1.54, 1.55},
		[15]float64{1.3, 1.31, 1.32, 1.33, 1.34, 1.35, 1.36, 1.37, 1.38, 1.39, 1.40, 1.41, 1.42, 1.43, 1.44},
		[15]float64{1.25, 1.26, 1.27, 1.28, 1.28, 1.29, 1.3, 1.31, 1.32, 1.33, 1.33, 1.34, 1.35, 1.36, 1.37},
		[15]float64{1.21, 1.22, 1.22, 1.24, 1.24, 1.25, 1.26, 1.26, 1.27, 1.28, 1.29, 1.29, 1.3, 1.31, 1.31},
		[15]float64{1.19, 1.19, 1.2, 1.21, 1.21, 1.22, 1.23, 1.23, 1.24, 1.24, 1.25, 1.26, 1.26, 1.27, 1.28},
	}

	Let        = [9]float64{0, -0.25, -0.5, -0.75, -1, -1.25, -1.5, -1.75, -2}
	LetStr     = [9]string{"平手", "平手/半球", "半球", "半球/一球", "一球", "一球/球半", "球半", "球半/两球", "两球"}
	Water      = [15]float64{0.75, 0.775, 0.8, 0.825, 0.85, 0.875, 0.9, 0.925, 0.95, 0.975, 1, 1.025, 1.05, 1.075, 1.1}
	water_min  = 0.75
	water_max  = 1.1
	water_step = 0.025
)

func equal(a, b float64) bool {
	return math.Dim(a, b) < 0.000001
}

func main() {
	fmt.Println("input odd:")
	var odd float64
	fmt.Scanf("%f\n", &odd)
	//fmt.Println("input", odd)
	//odd = 1.41
	for odd > 0 {

		for k1, l := range Odds {

			if odd < l[0]-0.01 || odd > l[14]+0.05 {
				continue
			}

			pos := -1
			for k, od := range l {
				if od < odd || equal(od, odd) {
					pos = k
					continue
				}
				break
			}

			if pos == -1 {
				continue
			}

			if pos == 14 {
				if odd > l[14] {
					fmt.Printf("让球:%s 大于%.3f\n", LetStr[k1], Water[14])
					continue
				}
			}

			if equal(odd, l[pos]) {
				fmt.Printf("让球:%s %.3f\n", LetStr[k1], Water[pos])
				continue
			}

			if odd > l[pos] && odd < l[pos+1] {
				fmt.Printf("让球:%s %.3f\n", LetStr[k1], Water[pos]+(Water[pos+1]-Water[pos])*((odd-l[pos])/(l[pos+1]-l[pos])))
				continue
			}
		}

		fmt.Println("input odd:")
		fmt.Scanf("%f\n", &odd)
	}

}
