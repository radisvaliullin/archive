package flatter

import "testing"

func TestFlatterInt(t *testing.T) {

	// test values
	var intSl []int
	var interfSl []interface{}
	var interfSlWithStr = []interface{}{"one"}

	// test cases
	testData := []struct {
		id      int
		testArr interface{}
		expRes  []int
		expErr  error
	}{
		// positive cases
		{
			id: 1,
			testArr: []interface{}{
				[]interface{}{
					1, 2,
					[]int{3},
				},
				4,
			},
			expRes: []int{1, 2, 3, 4},
		},
		{
			id:      2,
			testArr: intSl,
			expRes:  nil,
		},
		{
			id:      3,
			testArr: []int{},
			expRes:  []int{},
		},
		{
			id:      4,
			testArr: interfSl,
			expRes:  []int{},
		},
		{
			id: 5,
			testArr: []interface{}{
				intSl, intSl, interfSl,
			},
			expRes: []int{},
		},

		// negative cases
		{
			id:      10,
			testArr: nil,
			expErr:  ErrWrongType,
		},
		{
			id:      11,
			testArr: interfSlWithStr,
			expErr:  ErrWrongType,
		},
	}

	// pass test cases
	for _, td := range testData {

		a, err := Int(td.testArr)
		if err != nil && err != td.expErr {
			t.Fatalf("flatter error: in arr - %v; err - %v", td.testArr, err)
		} else if !isEqualIntSlice(a, td.expRes) {
			t.Fatalf("flatter wrong result: in arr - %v; out arr - %v; exp res - %v", td.testArr, a, td.expRes)
		}
		t.Logf("flatter: id - %v: in arr - %v, out arr - %v, expRes - %v, expErr - %v", td.id, td.testArr, a, td.expRes, err)
	}
}

//
func isEqualIntSlice(a, b []int) bool {

	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
