package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	err_val, err := cli()
	if err != nil {
		fmt.Println(err)
		os.Exit(err_val)
	}
}
