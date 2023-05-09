package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("go run main.go test1|test2|test3|test4|test5")
		os.Exit(1)
		return
	}
	args := map[string]func(){
		"test1": test1,
		"test2": test2,
		"test3": test3,
		"test4": test4,
		"test5": test5,
	}
	f, ok := args[os.Args[1]]
	if !ok {
		fmt.Println("go run main.go test1|test2|test3|test4|test5")
		os.Exit(1)
		return
	}
	f()
}

func test1() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	event := <-sigs
	fmt.Println(event)
	os.Exit(0)
}

func test2() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	event := <-sigs
	fmt.Println(event)
	os.Exit(1)
}

func test3() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	event := <-sigs
	fmt.Println(event)
	fmt.Println("fall asleep for 10 seconds")
	time.Sleep(10 * time.Second)
	fmt.Println("the sleep is over")
	os.Exit(0)
}

func test4() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	event := <-sigs
	fmt.Println(event)
	fmt.Println("fall asleep for 10 seconds")
	time.Sleep(10 * time.Second)
	fmt.Println("the sleep is over")
	os.Exit(1)
}

func test5() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	event := <-sigs
	fmt.Println(event)
	fmt.Println("fall asleep for 40 seconds")
	time.Sleep(40 * time.Second)
	fmt.Println("not killed")
	os.Exit(0)
}
