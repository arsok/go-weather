package main

import (
	"bufio"
	"fmt"
	"go-weather/internal"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("Enter a location to get weather information, or type 'exit' to quit:")
		command, inputErr := reader.ReadString('\n')
		if inputErr != nil {
			fmt.Println("Please try again: error reading input")
		} else {
			command = strings.TrimSpace(command[:len(command)-1])
			if command == "exit" {
				fmt.Println("Come again, rain or shine!")
				return
			}

			weather, weatherErr := internal.GetWeather(command)
			if weatherErr != nil {
				fmt.Println("Please try again:", weatherErr)
				continue
			}

			fmt.Println(weather)
		}
	}
}
