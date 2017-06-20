package ticket

import (
	"fmt"
)

// GoTickets is a interface of throttle system
type GoTickets interface {
	Take()
	Return()
	Active() bool
	Total() int
	Remainder() int
}

type baseGoTickets struct {
	total    int
	ticketCh chan byte
	active   bool
}

// NewGoTicket returns a GoTickets handler
func NewGoTicket(total int) (GoTickets, error) {
	gt := baseGoTickets{}
	if !gt.init(total) {
		return nil, fmt.Errorf("The goroutine ticket pool can NOT be initialized! (total=%d)", total)
	}
	return &gt, nil
}

func (gt *baseGoTickets) Take() {
	<-gt.ticketCh
}

func (gt *baseGoTickets) Return() {
	gt.ticketCh <- 1
}

func (gt *baseGoTickets) Active() bool {
	return gt.active
}

func (gt *baseGoTickets) Total() int {
	return gt.total
}

func (gt *baseGoTickets) Remainder() int {
	return len(gt.ticketCh)
}

func (gt *baseGoTickets) init(total int) bool {
	if gt.active {
		return false
	}
	if total == 0 {
		return false
	}

	ch := make(chan byte, total)
	for i := 0; i < total; i++ {
		ch <- 1
	}

	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}
