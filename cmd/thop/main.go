package main

import "os"

func main() {
	if err := root(autowireThop()).Execute(); err != nil {
		os.Exit(1)
	}
}
