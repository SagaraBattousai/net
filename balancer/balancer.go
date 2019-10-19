
package balancer

import (
  "container/heap"
)

func (b *Balancer) dispatch(req Request) {
  w := heap.Pop(&b.pool).(*Worker)

  w.requests <- req

  w.pending++

  heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
  w.pending--
  heap.Remove(&b.pool, w.index)
  heap.Push(&b.pool, w)
}

func (p pool) Len() int {
  return len(p)
}

func (p pool) resize() pool {
  n := len(p) * 2
  if n == 0 {
    n = 1
  }
  largerpool := make(pool, n)
  copy(largerpool, p)

  return largerpool
}


func (p pool) Less(i, j int) bool {
  return p[i].pending < p[j].pending
}

func (p pool) Swap(i, j int) {
  tmp := p[i]
  p[i] = p[j]
  p[j] = tmp
}

func (p *pool) Push(x interface{}) {
  if(cap(*p) == len(*p)) {
    *p = p.resize()
  }

  pv := *p

  nextIndex := len(pv)
  pv = pv[:nextIndex+1]
  w := x.(*Worker)
  w.index = nextIndex
  pv[nextIndex] = w

  *p = pv
}

func (p *pool) Pop() interface{} {
  pv := *p
  endIndex := len(pv) - 1
  w := pv[endIndex]
  pv = pv[:endIndex]

  *p = pv

  return w
}
