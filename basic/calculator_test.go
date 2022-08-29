package basic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		panic("Result is not 3")
	}
}

func TestSum2(t *testing.T) {
	result := Sum(2, 2)
	if result != 4 {
		panic("Result is not 4")
	}
}

func TestSumFail(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		t.Fail()
	}
	fmt.Println("TestSum done")
}

func TestSumFailNow(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		t.FailNow()
	}
	fmt.Println("TestSum done")
}

func TestSumError(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		t.Error("Result is not 3")
	}
	fmt.Println("TestSum done")
}

func TestSumFatal(t *testing.T) {
	result := Sum(1, 2)
	if result != 3 {
		t.Fatal("Result is not 3")
	}
	fmt.Println("TestSum done")
}

func TestSumAssert(t *testing.T) {
	result := Sum(1, 2)
	assert.Equal(t, 3, result)
	fmt.Println("TestSum done")
}

func TestSumRequire(t *testing.T) {
	result := Sum(1, 2)
	require.Equal(t, 3, result)
	fmt.Println("TestSum done")
}

func TestSumSub(t *testing.T) {
	x := 2

	t.Run("1plus2", func(t *testing.T) {
		result := Sum(1, x)
		assert.Equal(t, 3, result)
	})

	t.Run("2plus2", func(t *testing.T) {
		result := Sum(2, x)
		assert.Equal(t, 4, result)
	})
}

func TestSumTable(t *testing.T) {
	testCases := []struct {
		name     string
		paramA   int
		paramB   int
		expected int
	}{
		{
			name:     "1plus2",
			paramA:   1,
			paramB:   2,
			expected: 3,
		}, {
			name:     "2plus4",
			paramA:   2,
			paramB:   2,
			expected: 4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Sum(testCase.paramA, testCase.paramB)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
