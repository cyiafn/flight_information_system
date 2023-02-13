package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	UnsignedInt uint8
	Integer64   int64
	Integer     int
	Integer32   int32
	String      string
	Structure   struct {
		ABC   string
		Int64 int64
	}
	ArrayOfInt         []int
	ArrayOfString      []string
	ArrayOfEmptyString []string
	ArrayOfStruct      []struct {
		ABC   string
		Int64 int64
	}
}

type nestedStruct struct {
	ABC   string
	Int64 int64
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		Name      string
		TestValue *testStruct
	}{
		{
			Name:      "normal test",
			TestValue: newTestStruct(),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := Marshal(test.TestValue)
			assert.Nil(t, err)

			//logs.Info("%v", result)
			newStruct := &testStruct{}
			err = Unmarshal(result, newStruct)
			assert.Nil(t, err)

			assert.Equal(t, *test.TestValue, *newStruct)
		})
	}
}

func newTestStruct() *testStruct {
	return &testStruct{
		UnsignedInt: 5,
		Integer64:   51,
		Integer:     234,
		Integer32:   6431,
		String:      "hello world.",
		Structure: struct {
			ABC   string
			Int64 int64
		}{ABC: "weeew", Int64: 999999},
		ArrayOfInt:         []int{5, 2, 3},
		ArrayOfString:      []string{"hello", "world", "."},
		ArrayOfEmptyString: []string{},
		ArrayOfStruct: []struct {
			ABC   string
			Int64 int64
		}{
			{
				ABC:   "",
				Int64: 0,
			},
			{
				ABC:   "hadsfad",
				Int64: 123213,
			},
		},
	}
}
