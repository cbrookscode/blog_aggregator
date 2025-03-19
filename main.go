package main

import (
	"fmt"

	config "github.com/cbrookscode/blog_aggregator/internal/config"
)

func main() {
	test_struct, err := config.Read()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(test_struct)
	}
	return
}
