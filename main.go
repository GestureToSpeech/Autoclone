package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

func main() {
	cfg, err := toml.LoadFile("config.tml")
	catch(err)

	list := cfg.Get("repos").([]interface{})

	// Iterate through the list and print the strings
	for _, item := range list {
		fmt.Println(item.(string))
	}
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
