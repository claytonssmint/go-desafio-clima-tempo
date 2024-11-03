package utils

import "testing"

func TestIsvalidCep(t *testing.T) {
	valid := IsvalidCep("12345678")
	invalid := IsvalidCep("1234567")
	if !valid || invalid {
		t.Errorf("IsvalidCep failed")
	}
}
