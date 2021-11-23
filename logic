package main

import "fmt"

// Question 4
func main() {
	strings := []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}
	fmt.Println(groupByAnagrams(strings))
}

func groupByAnagrams(strings []string) [][]string {
	anagramMap := make(map[[2]int][]string)
	for _, str := range strings {
		// convert str to rune so it can be converted to int
		runes := []rune(str)
		// get sum and product of the array of runes
		sumAndProduct := calculateSumAndProduct(runes)
		// group strings together if their sum and product
		// are equal to each other
		anagramMap[sumAndProduct] = append(anagramMap[sumAndProduct], str)
	}

	// remap the strings from map into array to return the expected output
	anagramSlice := [][]string{}
	for _, value := range anagramMap {
		anagramSlice = append(anagramSlice, value)
	}

	return anagramSlice
}

func calculateSumAndProduct(runes []rune) [2]int {
	sum := 0
	product := 1

	for _, runeValue := range runes {
		sum += int(runeValue)
		product *= int(runeValue)
	}

	return [2]int{sum, product}
}
