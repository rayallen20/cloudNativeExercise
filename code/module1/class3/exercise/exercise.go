package main

import "fmt"

func main() {
	origin := []string{"I", "am", "stupid", "and", "week"}
	fmt.Println(origin)
	changeElementsByRange(origin)
	fmt.Println(origin)
}

func changeElementsByRange(slice []string) {
	for key, value := range slice {
		if value == "stupid" {
			slice[key] = "smart"
		}

		if value == "weak" {
			slice[key] = "strong"
		}
	}
}
