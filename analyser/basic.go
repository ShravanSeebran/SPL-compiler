package analyser

import (
	"fmt"
	"strings"
)

func ValidateTranslateToBasic(program []string) (instrs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	TranslateToBasic(program)
	return program, nil
}

func TranslateToBasic(program []string) {
	labelMap := getLabel(program)
	addLineNumbers(labelMap, program)
}

func getLabel(program []string) map[string]int {
	output := make(map[string]int)
	for i, line := range program {
		if strings.Contains(line, "REM") {
			label := strings.Split(line, " ")[1]
			output[label] = (i + 1) * 10
		}
	}
	return output
}

func addLineNumbers(labelMap map[string]int, program []string) {
	for i, line := range program {
		if strings.Contains(line, "GOTO") || strings.Contains(line, "THEN") {
			tokens := strings.Split(line, " ")
			label := tokens[len(tokens)-1]
			ln, ok := labelMap[label]
			if !ok {
				panic(fmt.Sprintf("label %s not found", label))
			}
			instruction := strings.Join(tokens[:len(tokens)-1], " ")

			program[i] = fmt.Sprintf("%s %d", instruction, ln)
		}
		program[i] = fmt.Sprintf("%-3d %s", (i+1)*10, program[i])
	}
}
