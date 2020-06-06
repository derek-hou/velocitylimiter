package main

import (
	"bufio"
	"fmt"
	"os"
)

/**
 *
 */
func main() {
	// parse through the text file
	file, _ := os.Open("input.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	//var data = Customers{}

	//_ = json.Unmarshal([]byte(file), &data)

	// for i := 0; i < len(data.Customers); i++ {
	// 	fmt.Println("Id: ", data.Customers[i].Id)
	// 	fmt.Println("Customer_id: ", data.Customers[i].Customer_id)
	// 	fmt.Println("Time: ", data.Customers[i].Time)
	// }

	// check based on business rule
	// output response to server
}
