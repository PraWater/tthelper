package excel_test

import (
	"reflect"
	"testing"

	"github.com/PraWater/tthelper/internal/excel"
)

func TestReadTT(t *testing.T) {
	cgot, sgot := excel.ReadTT("excel_test.xlsx")
	cwant := [][]string{{"BIO F215", "BIOPHYSICS", "3"}, {"BIO F231", "BIOLOGY PROJECT LAB", "3"}}
	swant := [][]string{{"BIO F215", "L1", "M W F  5"}, {"BIO F215", "T1", "Th  1"}, {"BIO F231", "L1", "S  1"}, {"BIO F231", "P1", "T  4 5  S  7 8"}}

	if !reflect.DeepEqual(cgot, cwant) {
		t.Errorf("Got: %s, want %s", cgot, cwant)
	}
	if !reflect.DeepEqual(sgot, swant) {
		t.Errorf("Got: %s, want %s", sgot, swant)
	}
}
