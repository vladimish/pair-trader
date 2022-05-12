package correlations

import (
	"testing"
)

func TestFindPearsonCorrelation(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6}
	y := []float64{2, 4, 7, 9, 12, 14}

	r := FindPearsonCorrelation(x, y)
	if r < 0.9980 || r > 0.9985 {
		t.Errorf("Wrong answer. Got: %f, Exp: ~0.988", r)
	}
	//fmt.Println(r)
}
