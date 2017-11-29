package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("file", "problems/default.csv", "The file to get quiz questions and answers from")
	time := flag.Int64("time", 5, "Set the time for the quiz")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Cannot open file %s", *filename)
	}

	reader := csv.NewReader(file)
	problems := parseProblems(reader)
	done := make(chan bool)
	scoreUpdate := make(chan uint32)
	go askQuestions(problems, scoreUpdate, done)
	var score uint32
	go timeout(*time, done)
	for {
		select {
		case <-done:
			fmt.Printf("Your final score is %v/%v.", score, len(problems))
			return
		case point := <-scoreUpdate:

			score += point
		}
	}
}

type problem struct {
	Question string
	Answer   string
}

func parseProblems(file *csv.Reader) []problem {
	lines, _ := file.ReadAll()
	problems := make([]problem, len(lines))
	for idx, line := range lines {
		problems[idx] = problem{
			Question: line[0],
			Answer:   line[1],
		}
	}
	return problems
}

func askQuestions(problems []problem, score chan uint32, done chan bool) {
	in := bufio.NewReader(os.Stdin)
	for idx, problem := range problems {
		fmt.Printf("Problem #%d:\n%s ", idx+1, problem.Question)
		ans, _ := in.ReadString('\n')
		if strings.TrimSpace(ans) == problem.Answer {
			score <- 1
		}
	}
	fmt.Println("Good job! You solved all the puzzles.")
	done <- true
}

func timeout(seconds int64, done chan bool) {
	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
	fmt.Println("\nTime elapsed.")
	done <- true
}
