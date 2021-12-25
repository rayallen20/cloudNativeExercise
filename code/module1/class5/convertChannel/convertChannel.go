package main

func main() {
	c := make(chan int)
	go prod(c)
	go consume(c)
}

func prod(ch chan<- int) {
	for {
		ch <- 1
	}
}

func consume(ch <-chan int) {
	for {
		<-ch
	}
}
