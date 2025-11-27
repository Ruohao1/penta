package report

import (
	"fmt"
	"io"
	"sync"

	"github.com/Ruohao1/penta/internal/scan"
)

type orderedPrinter struct {
	mu          sync.Mutex
	nextToPrint int
	pending     map[int]scan.Result
	w           io.Writer
}

func NewOrderedPrinter(w io.Writer) *orderedPrinter {
	return &orderedPrinter{
		pending: make(map[int]scan.Result),
		w:       w,
	}
}

// Add is safe to call from concurrent goroutines.
func (op *orderedPrinter) Add(idx int, res scan.Result) {
	op.mu.Lock()
	defer op.mu.Unlock()

	op.pending[idx] = res
	op.drainLocked()
}

// Flush is optional; just in case you want to force draining at end.
func (op *orderedPrinter) Flush() {
	op.mu.Lock()
	defer op.mu.Unlock()
	op.drainLocked()
}

func (op *orderedPrinter) drainLocked() {
	for {
		res, ok := op.pending[op.nextToPrint]
		if !ok {
			return
		}
		delete(op.pending, op.nextToPrint)

		fmt.Fprintf(
			op.w,
			"host=%s status=%s method=%s meta=%v\n",
			res.Addr.String(),
			res.Status,
			res.Method,
			res.Meta,
		)

		op.nextToPrint++
	}
}
