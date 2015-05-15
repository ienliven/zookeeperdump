package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func walk(path string, depth int, ch chan string) {
	defer close(ch)
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	var innerWalk func(string, int)
	innerWalk = func(path string, depth int) {
		if depth <= 0 {
			return
		}
		children, _, err := conn.Children(path)
		if err != nil {
			log.Print(err)
		}
		if len(children) == 0 {
			return
		}
		wg.Add(len(children))
		for _, c := range children {
			ch <- filepath.Join(path, c)
			go func(path, c string, depth int) {
				defer wg.Done()
				innerWalk(filepath.Join(path, c), depth-1)
			}(path, c, depth)
		}
	}
	innerWalk(path, depth)
	wg.Wait()
}

func main() {
	res := make(chan string)
	go walk("/", 20, res)
	for x := range res {
		fmt.Println(x)
	}
}
