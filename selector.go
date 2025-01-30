package main

import (
	"fmt"
	"os"
	"strconv"
)

func selector[slice any](s []slice) int {
	var choice int
	for {
		var choiceStr string
		_, err := fmt.Fscanln(os.Stdin, &choiceStr)
		if err != nil {
			fmt.Print("\033[1A\033[K")
			continue
		}
		choice, err = strconv.Atoi(choiceStr)
		if err != nil || choice >= len(s) {
			fmt.Print("\033[1A\033[K")
		} else {
			break
		}
	}
	fmt.Print("\033[1A\033[K")
	fmt.Printf("[%d]\n", choice)
	return choice
}
