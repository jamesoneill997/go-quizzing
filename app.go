package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func initialise(filename string) map[string]string {

	qna := make(map[string]string)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file")
		log.Fatal()
	}
	dataStr := string(data)

	reader := csv.NewReader(strings.NewReader(dataStr))

	records, e := reader.ReadAll()

	if e != nil {
		fmt.Println("Error reading csv")
		log.Fatal()
	}

	for _, pair := range records {
		qna[pair[0]] = pair[1]
	}

	return qna
}

func main() {
	initQuiz := initialise(os.Args[1])
	questions := make([]string, 0, len(initQuiz))
	score := 0

	for k := range initQuiz {
		questions = append(questions, k)
	}

	for i := 0; i < len(questions); i++ {
		question := questions[i]
		answer := initQuiz[question]
		fmt.Printf("Question %d/%d", i+1, len(questions))
		fmt.Println()
		fmt.Println(question)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		guess := scanner.Text()

		if guess == answer {
			score++

			fmt.Println(strings.Repeat("-", 15), "\n ")
			fmt.Println("Correct! \n")
			fmt.Println(strings.Repeat("-", 15), "\n ")
		} else {
			fmt.Println(strings.Repeat("-", 15), "\n ")
			fmt.Printf("Incorrect, the correct answer was %s \n", answer)
			fmt.Println(strings.Repeat("-", 15), "\n ")

		}

	}

	fmt.Printf("Quiz completed! \n You scored %d out of a possible %d \n", score, len(questions))

}
