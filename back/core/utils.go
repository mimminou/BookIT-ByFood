package core

import "os"

//Utility functions for debugging

// Writes anything passed to it to query.txt
func LogQuery(query string) {
	file, _ := os.OpenFile("query.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(query)
	defer file.Close()
}
