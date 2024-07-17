package timetable_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/PraWater/tthelper/internal/timetable"
)

func TestInput(t *testing.T) {
  file, err := os.Open("timetable_test.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
	got, err := timetable.ReadFile(file)
	want := [][]string{{"BIO F215", "L1"}, {"BIO F215", "T1"}, {"BIO F231", "L1"}}

	assertNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseSlot(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Output [][]int
	}{
		{
			"Single day with one time slot",
			"Th  1",
			[][]int{{3, 1}},
		},
		{
			"Multiple days with same time slot",
			"M W F  5",
			[][]int{{0, 5}, {2, 5}, {4, 5}},
		},
		{
			"Multiple days with different time slot",
			"M  5  T  4",
			[][]int{{0, 5}, {1, 4}},
		},
		{
			"Single day with multiple time slot",
			"T  4 5",
			[][]int{{1, 4}, {1, 5}},
		},
		{
			"Multiple days with multiple time slot",
			"T  4 5  S  7 8",
			[][]int{{1, 4}, {1, 5}, {5, 7}, {5, 8}},
		},
    {
      "Slot with newline character",
      `T 4 5 S
      7 8`,
			[][]int{{1, 4}, {1, 5}, {5, 7}, {5, 8}},
    },
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
      got, err := timetable.ParseSlot(test.Input)

      assertNoError(t, err)
			if !reflect.DeepEqual(got, test.Output) {
				t.Errorf("got %v, want %v", got, test.Output)
			}
		})
	}

  t.Run("Unexpected day character", func(t *testing.T) {
    slot := "P  7"
    _, err := timetable.ParseSlot(slot)

    if err == nil {
      t.Errorf("Expected error when giving %q as input\n", slot)
    }
  })
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
