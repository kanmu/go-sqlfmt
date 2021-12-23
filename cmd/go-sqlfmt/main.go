package main

import "flag"

func main() {
	flag.Usage = usage
	flag.Parse()

	sqlfmtMain()
}
