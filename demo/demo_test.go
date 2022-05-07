package demo

import "testing"

func TestAdd(t *testing.T) {
	actual := Add(1, 2)
	expected := 3
	if actual != expected {
		t.Fail()
	}
}
