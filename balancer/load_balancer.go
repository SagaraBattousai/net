
package balancer

import (
  "container/heap"
  "fmt"
  "runtime"
  "net/http"
)

type Request struct {
  fn func(http.ResponseWriter, *http.Request)
  c chan int
}

type Worker struct {
  requests chan Request
  pending int
  index int
}

type pool []*Worker

type Balancer struct {
  pool pool
  done chan *Worker
}

func requester(work chan<- Request) {
  c := make(chan int)
  for {
    work <- Request{tmp, c}
    result := <-c
    fmt.Println(result)
  }
}
func (w *Worker) work(done chan *Worker) {
  for {
    req := <-w.requests
    req.c <- req.fn()
    done <- w

  }
}
func (b *Balancer) balance(work chan Request) {
  for {
    select {
    case req := <-work:
      b.dispatch(req)
    case w:= <-b.done:
      b.completed(w)
    }
  }
}

func run() {
  reqChannel := make(chan Request)
  workOverChannel := make(chan *Worker)

  go requester(reqChannel)

  numWorkers := runtime.GOMAXPROCS(0)

  pool := make(pool, 0, numWorkers)

  for i := 0; i < numWorkers; i++ {
    w := createWorker()
    go w.work(workOverChannel)
    heap.Push(&pool, w)
  }

  //heap.Init(pool)

  bal := Balancer{pool, workOverChannel}

  go bal.balance(reqChannel)

}

func createWorker() *Worker {
  reqChannel := make(chan Request)
  return &Worker{reqChannel, 0, -1}
}

/*
func main() {

  cmd := exec.Command("python", "-c", "while True:\n\tx=input()\n\tprint('python', x)")
  stdin, _ := cmd.StdinPipe()
  stdout, _ := cmd.StdoutPipe()
  cmd.Start()
  b := make([]byte, 50)
  for {
    stdin.Write([]byte("hello from Go\n"))
    stdout.Read(b)
    fmt.Printf("%s\n", b)
  }

}
*/



































