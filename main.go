package quiz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	filename := flag.String("file", "problems/default.csv", "The file to get quiz questions and answers from")

	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Cannot open file %s", *filename)
	}

	reader := csv.NewReader(file)
	go questions(reader)

}

func questions(reader *csv.Reader) {

}
