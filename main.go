package main

import "flag"

func main() {
	configFile := flag.String("config", "", "config file")
	flag.Parse()

	_ = configFile

	return
}
