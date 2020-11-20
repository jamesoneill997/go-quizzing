package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
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

func printScore(score int, totalQuestions int) string {
	return fmt.Sprintf("Quiz completed! \n You scored %d out of a possible %d \n", score, totalQuestions)
}

func quiz(initQuiz map[string]string, timer *time.Timer) {

	questions := make([]string, 0, len(initQuiz))

	totalQuestions := len(initQuiz)

	score := 0

	for k := range initQuiz {
		questions = append(questions, k)
	}

	for i := 0; i < totalQuestions; i++ {
		question := questions[i]
		answer := strings.ToLower(initQuiz[question])
		fmt.Printf("Question %d/%d", i+1, totalQuestions)
		fmt.Println()
		fmt.Println(question)

		answerCh := make(chan string, 1)

		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			guess := scanner.Text()
			answerCh <- guess
		}()

		select {

		case <-timer.C:
			fmt.Printf("You ran out of time. You scored %d out of %d \n", score, totalQuestions)
		case guess := <-answerCh:
			if guess == answer {
				score++

				fmt.Println(strings.Repeat("-", 15), "\n ")
				fmt.Println("Correct! \n ")
				fmt.Println(strings.Repeat("-", 15), "\n ")
			} else {
				fmt.Println(strings.Repeat("-", 15), "\n ")
				fmt.Printf("Incorrect, the correct answer was %s \n", answer)
				fmt.Println(strings.Repeat("-", 15), "\n ")

			}
		}

	}

	printScore(score, totalQuestions)

}

func main() {
	//handle file input
	file := "./qa.csv"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	initQuiz := initialise(file)
	timer := time.NewTimer(time.Duration(3) * time.Second)

	quiz(initQuiz, timer)
}
