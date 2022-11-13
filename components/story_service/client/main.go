package main

const (
	address = "localhost:8080"

	// Adjust the size for which a large file will be broken
	// down into multiple parts during a stream request
	chunkSize = 5 * MB
)

const MB = 1 << 20

func main() {

}
