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

func printScore(score int, totalQuestions int) {
	fmt.Printf("Quiz completed! \n You scored %d out of a possible %d \n", score, totalQuestions)
}

func timer(ch chan int, wg *sync.WaitGroup) {
	timer := time.NewTimer(time.Second * 3)
	<-timer.C
	ch <- -1
	close(ch)
	wg.Done()
}

func quiz(initQuiz map[string]string, ch chan int, wg *sync.WaitGroup) {

	questions := make([]string, 0, len(initQuiz))

	totalQuestions := len(initQuiz)

	score := 0

	for k := range initQuiz {
		questions = append(questions, k)
	}

	fmt.Println("Press enter to start the quiz.")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	for i := 0; i < totalQuestions; i++ {
		question := questions[i]
		answer := strings.ToLower(initQuiz[question])
		fmt.Printf("Question %d/%d", i+1, totalQuestions)
		fmt.Println()
		fmt.Println(question)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		guess := scanner.Text()

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

	printScore(score, totalQuestions)
	ch <- 0
	close(ch)
	wg.Done()

}

func main() {
	//handle file input
	file := "./qa.csv"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	initQuiz := initialise(file)

	timerChan := make(chan int)
	quizChan := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)
	go timer(timerChan, &wg)
	go quiz(initQuiz, quizChan, &wg)

	select {
	case status := <-timerChan:
		fmt.Println(status)
	case status := <-quizChan:
		fmt.Println(status)

		wg.Wait()
	}
}
