package input

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(filename string) (sections [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
