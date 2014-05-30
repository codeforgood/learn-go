package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
	"runtime"
)

func randNumGen() <-chan int{
	out := make(chan int, 100)
	go func(){
		for i := 0; i < 100; i++ {
			fmt.Println("Generating random number")
			out <- rand.Intn(10)
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int{
	out := make(chan int)
	go func(){
		for num := range in{
			fmt.Println("Squaring from channel")
			out <- num * num
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))
	c := randNumGen()
	out1 := sq(c)
	out2 := sq(c)
	out3 := sq(c)
	out4 := sq(c)

	for p := range merge(out1, out2, out3, out4) {
		fmt.Println("Square: ", p)
	}
}
