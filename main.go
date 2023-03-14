package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	fileFlag         string
	timeLimit        int
	shouldBeShuffled bool
	timeInterrupt    = make(chan bool)
	newAnswer        = make(chan bool)
)

func init() {
	const (
		defaultFile      = "./problems.csv"
		fileUsage        = "The problems CSV file path"
		defaultTimeLimit = 30
		timeLimitUsage   = "Number of seconds to finish the quiz"
	)

	flag.StringVar(&fileFlag, "file", defaultFile, fileUsage)
	flag.StringVar(&fileFlag, "f", defaultFile, fileUsage+" (shorthand)")
	flag.IntVar(&timeLimit, "time-limit", defaultTimeLimit, timeLimitUsage+" (shorthand)")
	flag.IntVar(&timeLimit, "t", defaultTimeLimit, timeLimitUsage)
	flag.BoolVar(&shouldBeShuffled, "shuffle", false, timeLimitUsage)
	flag.BoolVar(&shouldBeShuffled, "s", false, timeLimitUsage)
}

func main() {
	flag.Parse()

	records := readCSVFromFile(fileFlag)

	fmt.Print("\nPress [enter] to start the quiz")
	fmt.Scanln()
	fmt.Println()

	go startTimer()

	var correctAnswersCount byte

	loopOverQuestions(records, &correctAnswersCount)

	fmt.Println("\n---------------------------------\n")
	fmt.Printf("You answered %d questions correctly out of %d\n", correctAnswersCount, len(records))
}

func readCSVFromFile(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func loopOverQuestions(records [][]string, correctAnswersCount *byte) {
	for index, record := range records {
		go askUser(index+1, record)

		select {
		case <-timeInterrupt:
			return
		case isCorrect := <-newAnswer:
			if isCorrect {
				*correctAnswersCount++
			}
		}
	}
}

func askUser(questionNumber int, record []string) {
	var userAnswer string
	question, correctAnswer := record[0], record[1]

	fmt.Printf("Q #%d: %s ?\n", questionNumber, question)
	fmt.Print("Your answer: ")
	fmt.Scan(&userAnswer)

	defer fmt.Print("\n\n")

	newAnswer <- userAnswer == correctAnswer
}

func startTimer() {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	fmt.Print("\n\nTime Up!!!\n")

	timeInterrupt <- true
}
