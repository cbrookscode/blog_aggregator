package main

import (
	"fmt"

	config "github.com/cbrookscode/blog_aggregator/internal/config"
)

func main() {
	test_struct, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(test_struct)

	err = test_struct.SetUser("Hilda Brown")
	if err != nil {
		fmt.Println(err)
		return
	}

	test_struct, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(test_struct)
}
