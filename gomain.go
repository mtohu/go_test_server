package main

import (
	. "fmt"
	"gotest/gopool"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)
var p *gopool.Pool
func main()  {
	p = gopool.NewPool(5, 10).
		Start().EnableWaitForAll(true)
	defer p.WaitForAll()
	defer p.StopAll()
	//defer p.StopAll()
	http.HandleFunc("/hello", hello)
	http.Handle("/rest/",http.HandlerFunc(say))
	//go func() {
	http.ListenAndServe("localhost:6069", nil)
	//}()

}
func say(w http.ResponseWriter, req *http.Request) {
	Print("=========333=======")
	//http.Handle("/sd",http.HandlerFunc(bbbs))
}
func bbbs(w http.ResponseWriter, req *http.Request)  {
	w.Write([]byte("2222Hello11111"))
}
func enqueue(q chan int) {

	time.Sleep(time.Second * 3)

	q <- 10

	close(q)

}


func hello(w http.ResponseWriter, req *http.Request) {
	var wg sync.WaitGroup

	var c =make(chan int,5)
	var bol =make(chan bool)
	for i := 0; i < 5; i++ {
		count := i
		wg.Add(1)
		p.AddJob(func() {
		  //fmt.Println("uuuu",count)
		  go cdo(count,c,&wg,w,req)
		})

	}
	//time.Sleep(2 * time.Second)
	w.Write([]byte("Hello11111"))
	/*for {
		select {
		    case x, ok := <- c:
				if ok {
					fmt.Println("----",x)
				}else{
					fmt.Println("----error")
					return
				}
		    default:
				fmt.Println("waiting")

		}
	}*/
	Println("0000000")
	//等待并取出channelc中的值，直到channel关闭，会阻塞
	var numbers []int

	go func(){
		for term := range c {

			Println("ggggggg------",term)
			numbers = append(numbers, term)
		}
		Println("bbbbbbbbb------")
		bol <- true
	}()
	wg.Wait()
	//WaitTimeout(&wg, 5 * time.Second)
	close(c)
	<-bol  //等待 c 通过结束
	close(bol)
	Printf("len=%d cap=%d slice=%v\n",len(numbers),cap(numbers),numbers)

}
/*****设置等待超时*****/
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func cdo(n int,q chan int,wg *sync.WaitGroup,w http.ResponseWriter, req *http.Request)  {
	defer wg.Done()
	time.Sleep(10 * time.Millisecond)
	//fmt.Printf("%d\r\n", n)
	q <- n


}
