package timetable

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadFile(file *os.File) (sections [][]string, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		for i := 2; i < len(fields); i++ {
			sections = append(sections, []string{strings.Join(fields[:2], " "), fields[i]})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return
}

func ParseSlot(slot string) (times [][]int, err error) {
	sSplit := strings.Split(slot, " ")
  splitSlots := []string{}
  for _, sp := range sSplit {
    splitSlots = append(splitSlots, strings.Split(sp, "\n")...)
  }
	previousNum := false
	currentDays := []int{}
	for _, splits := range splitSlots {
		s := strings.Trim(splits, " ")
		if len(s) == 0 {
			continue
		}
		if time, e := strconv.Atoi(s); e == nil {
      previousNum = true
      for _, day := range currentDays {
        times = append(times, []int{day, time})
      }
		} else {
      if previousNum {
        currentDays = nil
      }
      previousNum = false
      day := 0
      switch s {
      case "M":
        day = 0
      case "T":
        day = 1
      case "W":
        day = 2
      case "Th":
        day = 3
      case "F":
        day = 4
      case "S":
        day = 5
      default:
        err = errors.New("Unexpected string for day")
        return
      }
      currentDays = append(currentDays, day)
    }
	}
	return
}
