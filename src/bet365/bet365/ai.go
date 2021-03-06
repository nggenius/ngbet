package bet365

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

var (
	W1 = mat.NewDense(40, 1, []float64{
		8.27173516e-02,
		1.10541552e-01,
		6.95866272e-02,
		8.81155670e-01,
		1.77493319e-01,
		8.67342353e-02,
		3.92892994e-02,
		6.81846344e-04,
		5.17368899e-04,
		5.39300032e-04,
		-9.78800189e-03,
		5.85713312e-02,
		6.26042843e-01,
		1.01009011e+00,
		2.83224165e-01,
		-2.41606489e-01,
		-4.03864793e-02,
		-6.41123438e-03,
		2.36779347e-01,
		-6.65536821e-01,
		-1.02196574e-01,
		-3.93382221e-01,
		-3.10864169e-02,
		2.98097786e-02,
		1.42930135e-01,
		-3.97605501e-04,
		-3.47402529e-04,
		3.01471772e-03,
		-2.87759416e-02,
		-7.42398351e-02,
		-1.83466300e-01,
		-3.67188662e-01,
		-6.29819632e-01,
		2.84140646e-01,
		5.35171770e-04,
		3.11689824e-03,
		-5.00448793e-02,
		7.13185489e-01,
		-3.67122948e-01,
		-2.89085835e-01,
	})

	b1 = 0.3038962

	W2 = mat.NewDense(40, 1, []float64{
		-1.7900032,
		3.1018026,
		-0.3646772,
		1.0133697,
		0.98701483,
		0.4879744,
		0.2108514,
		-0.089339,
		0.01485642,
		1.2743756,
		-0.37466604,
		-0.05297791,
		-0.59075314,
		0.02118768,
		0.6159304,
		0.6311794,
		-0.08429201,
		-0.10814717,
		-0.1529978,
		-1.2239702,
		-0.0711984,
		-0.5028359,
		-0.45827562,
		-2.0819302,
		-1.2196221,
		0.0698214,
		-0.06364495,
		0.25537458,
		-0.03107651,
		0.6696233,
		0.00533889,
		0.1231402,
		-0.03139678,
		0.43439427,
		0.04774024,
		0.1013598,
		-0.13673092,
		0.65206677,
		0.02756749,
		0.4696105,
	})

	b2 = 0.00270589
)

func forecast(data [40]float64) float64 {
	c := mat.NewDense(1, 40, data[:])
	out := &mat.Dense{}
	out.Mul(c, W1)

	fmt.Println(out)
	ar, ac := out.Dims()
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			v := out.At(i, j) + b1
			v = 1.0 / (1.0 + math.Exp(-v))
			out.Set(i, j, v)
		}
	}
	return out.At(0, 0) * 100
}

func forecastHalf(data [40]float64) float64 {
	c := mat.NewDense(1, 40, data[:])
	out := &mat.Dense{}
	out.Mul(c, W2)

	fmt.Println(out)
	ar, ac := out.Dims()
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			v := out.At(i, j) + b2
			v = 1.0 / (1.0 + math.Exp(-v))
			out.Set(i, j, v)
		}
	}
	return out.At(0, 0) * 100
}
