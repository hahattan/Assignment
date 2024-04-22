package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var symbols = []string{"+", "-", "*", "/"}

type Student struct {
	id         string
	questionCh <-chan string
	answerCh   chan<- Answer
}

type Teacher struct {
	questionCh    chan<- string
	answerCh      <-chan Answer
	correctAnswer int
	question      []string
}

type Answer struct {
	Value   int
	Student string
}

func NewStudent(id string, questionCh <-chan string, answerCh chan<- Answer) *Student {
	return &Student{
		id:         id,
		questionCh: questionCh,
		answerCh:   answerCh,
	}
}

func NewTeacher(questionCh chan<- string, answerCh <-chan Answer) *Teacher {
	return &Teacher{
		questionCh: questionCh,
		answerCh:   answerCh,
		question:   make([]string, 3),
	}
}

func (t *Teacher) Ask(questionReady chan struct{}) {
	time.Sleep(3 * time.Second)
	fmt.Println("\nTeacher: Guys, are you ready?")
	a, b, symbolIdx := 1+rand.IntN(99), 1+rand.IntN(99), rand.IntN(4)
	fmt.Printf("Teacher: %d %s %d = ?\n", a, symbols[symbolIdx], b)

	t.questionCh <- strconv.Itoa(a)
	t.questionCh <- symbols[symbolIdx]
	t.questionCh <- strconv.Itoa(b)

	t.question[0] = strconv.Itoa(a)
	t.question[1] = symbols[symbolIdx]
	t.question[2] = strconv.Itoa(b)
	switch t.question[1] {
	case "+":
		t.correctAnswer = a + b
	case "-":
		t.correctAnswer = a - b
	case "*":
		t.correctAnswer = a * b
	case "/":
		t.correctAnswer = a / b
	}
	// signal student to answer
	questionReady <- struct{}{}
}

func (t *Teacher) Check(questionReady chan struct{}) string {
	var a Answer
	var count int
	var winner string
	for {
		a = <-t.answerCh
		if a.Value == t.correctAnswer {
			winner = a.Student
			break
		}
		fmt.Printf("Teacher: %s, you are wrong.\n", a.Student)

		count++
		if count == 5 {
			break
		}

		t.questionCh <- t.question[0]
		t.questionCh <- t.question[1]
		t.questionCh <- t.question[2]
		// signal student to answer
		questionReady <- struct{}{}
	}

	if winner == "" {
		fmt.Printf("Teacher: Boooo~ Answer is %d.\n", t.correctAnswer)
	} else {
		fmt.Printf("Teacher: %s, you are right!\n", winner)
	}

	close(questionReady)
	return winner
}

func (s *Student) Answer(questionReady chan struct{}) {
	select {
	case <-questionReady:
		if len(s.questionCh) == 0 {
			return
		}
		time.Sleep(time.Second * (1 + time.Duration(rand.IntN(2))))
		var nums []int
		var symbol string
		for i := 0; i < 3; i++ {
			q := <-s.questionCh
			num, err := strconv.Atoi(q)
			if err != nil {
				symbol = q
				continue
			}
			nums = append(nums, num)
		}
		if len(nums) == 0 {
			return
		}

		var answer int
		if time.Now().Unix()%2 == 0 {
			switch symbol {
			case "+":
				answer = nums[0] + nums[1]
			case "-":
				answer = nums[0] - nums[1]
			case "*":
				answer = nums[0] * nums[1]
			case "/":
				answer = nums[0] / nums[1]
			}
		} else {
			// simulate wrong answer
			answer = rand.IntN(100)
		}

		fmt.Printf("Student %v: %d %s %d = %d\n", s.id, nums[0], symbol, nums[1], answer)
		s.answerCh <- Answer{Value: answer, Student: s.id}
	}
}

func (s *Student) Praise(id string) {
	fmt.Printf("Student %v: %v, you win.\n", s.id, id)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	questionCh := make(chan string, 3)
	answerCh := make(chan Answer, 1)
	teacher := NewTeacher(questionCh, answerCh)
	students := make([]*Student, 5)
	for i := 0; i < 5; i++ {
		students[i] = NewStudent(string('A'+i), questionCh, answerCh)
	}

	for {
		questionReady := make(chan struct{}, 1)
		select {
		case <-ctx.Done():
			close(answerCh)
			close(questionCh)
			close(questionReady)
			return
		default:
			// prepare student to answer
			for _, student := range students {
				go func(s *Student) {
					s.Answer(questionReady)
				}(student)
			}
			// teacher ask question
			teacher.Ask(questionReady)
			// teacher check answer
			winner := teacher.Check(questionReady)
			// student praise the winner
			if winner != "" {
				for _, s := range students {
					if s.id != winner {
						s.Praise(winner)
					}
				}
			}
		}
	}
}
