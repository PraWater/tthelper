package input_test

import (
	"testing"
  "reflect"

	"github.com/PraWater/tthelper/internal/input"
)

func TestInput(t *testing.T) {
  got, err := input.ReadFile("input_test.txt")
	want := [][]string{{"BIO F215", "L1"}, {"BIO F215", "T1"}, {"BIO F231", "L1"}}

  assertNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
