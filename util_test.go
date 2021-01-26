package ne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	//given
	input := "test"
	//when
	output := String(input)
	//then
	assert.Equal(t, input, *output, "Input value is same as value of output's pointer")
}

func TestStringValue(t *testing.T) {
	//given
	inputString := "test"
	input := []*string{
		&inputString,
		nil,
	}
	expected := []string{
		inputString,
		"",
	}
	//when
	output := make([]string, len(input))
	for i := range input {
		output[i] = StringValue(input[i])
	}
	//then
	assert.Equal(t, expected, output, "Output matches expected output")
}

func TestInt(t *testing.T) {
	//given
	input := 20
	//when
	output := Int(input)
	//then
	assert.Equal(t, input, *output, "Input value is same as value of output's pointer")
}

func TestIntValue(t *testing.T) {
	//given
	inputInt := 101
	input := []*int{
		&inputInt,
		nil,
	}
	expected := []int{
		inputInt,
		0,
	}
	//when
	output := make([]int, len(input))
	for i := range input {
		output[i] = IntValue(input[i])
	}
	//then
	assert.Equal(t, expected, output, "Output matches expected output")
}

func TestBool(t *testing.T) {
	//given
	input := false
	//when
	output := Bool(input)
	//then
	assert.Equal(t, input, *output, "Input value is same as value of output's pointer")
}

func TestBoolValue(t *testing.T) {
	//given
	inputBool := false
	input := []*bool{
		&inputBool,
		nil,
	}
	expected := []bool{
		inputBool,
		false,
	}
	//when
	output := make([]bool, len(input))
	for i := range input {
		output[i] = BoolValue(input[i])
	}
	//then
	assert.Equal(t, expected, output, "Output matches expected output")
}
