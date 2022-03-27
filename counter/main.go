package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	i := 0
	for {
		if i == 5 {
			return
		}

		fmt.Fprintln(os.Stdout, "stdout: ", i)
		fmt.Fprintln(os.Stderr, "stderr: ", i*2)
		time.Sleep(1 * time.Second)

		i++
	}
}
