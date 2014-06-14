package main

import (
    "fmt"
    "sync"
    "math/rand"
    "time"
    "github.com/garyburd/redigo/redis"
    )

func pub(c redis.Conn){
    for i := 0; i < 2; i++ {
        time.Sleep(1e9)
        c.Do("PUBLISH", "example", rand.Intn(10))
    }
}

func sub(psc redis.PubSubConn, wg sync.WaitGroup){
    go func(){
        for {
            switch v := psc.Receive().(type) {
            case redis.Message:
                fmt.Printf("%s: %s\n", v.Channel, v.Data)
            case redis.Subscription:
                fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
                if v.Count == 0 {
                    fmt.Printf("Returning from sub routine")
                    wg.Done()
                    return
                }
            case error:
                fmt.Printf("error: %v\n", v)
                wg.Done()
                return
            }
        }
    }()
}

func main(){
    c1, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("Connection to redis failed")
    }
    defer c1.Close()
    c2, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("Connection to redis failed")
    }
    defer c2.Close()
    psc := redis.PubSubConn{Conn: c1}

    var wg sync.WaitGroup
    wg.Add(1)
    psc.Subscribe("example")
    sub(psc, wg) // go routine

    pub(c2)
    psc.Unsubscribe("example")
    wg.Wait()
    fmt.Printf("Exciting from pubsub module")
}
