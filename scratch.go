package main

import "fmt"

func sum(start int, end int) <-chan int{
    out := make(chan int)
    //goroutine
    go func(){
        sum := 0
        for i := start; i < end; i++ {
            sum += i
        }
        out <- sum
        close(out)
    }()
    return out
}

func main() {
    c1 := sum(1, 50) // call to sum returns immediately with a channel
    c2 := sum(50, 100)
    sum_first_half, sum_second_half := <-c1, <-c2
    fmt.Println(sum_first_half + sum_second_half)
}
