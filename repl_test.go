package main

import (
	"testing"
	"reflect"
)

func TestCleanInput(t *testing.T) {
	resultOne := cleanInput("hello world")
	resultTwo := cleanInput("Charmander Bulbasaur PIKACHU")
	expectedOne := []string{"hello", "world"}
	expectedTwo := []string{"charmander", "bulbasaur", "pikachu"}
	if !reflect.DeepEqual(resultOne, expectedOne) {
    	t.Errorf("expected %v, got %v", expectedOne, resultOne)
	}
	if !reflect.DeepEqual(resultTwo, expectedTwo) {
		t.Errorf("expected %v, got %v", expectedTwo, resultTwo)
	}
}