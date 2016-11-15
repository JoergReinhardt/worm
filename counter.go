package main

import (
	"strings"
	"time"
)

type digits []string

func newDigits() digits {
	var dig []string

	dig = append(dig, "┏━┓")
	dig = append(dig, "┃┃┃")
	dig = append(dig, "┗━┛")
	dig = append(dig, "╺┓ ")
	dig = append(dig, " ┃ ")
	dig = append(dig, "╺┻╸")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┏━┛")
	dig = append(dig, "┗━╸")
	dig = append(dig, "┏━┓")
	dig = append(dig, "╺━┫")
	dig = append(dig, "┗━┛")
	dig = append(dig, "╻ ╻")
	dig = append(dig, "┗━┫")
	dig = append(dig, "  ╹")
	dig = append(dig, "┏━╸")
	dig = append(dig, "┗━┓")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┣━┓")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "  ┃")
	dig = append(dig, "  ╹")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┣━┫")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┗━┫")
	dig = append(dig, "┗━┛")

	return dig
}

type sector struct {
	pos  int
	str  string
	prev *sector
	next *sector
}
type ring struct {
	*sector
}

func (r *ring) stringRow() string { return r.str }
func (r *ring) stringDig() string {
	ret := r.prev.prev.str + "\n"
	ret = ret + r.prev.str + "\n"
	ret = ret + r.str + "\n"
	return ret
}
func (r *ring) advanceRow() {
	(*r).sector = (*r).next
}
func (r *ring) backupRow() {
	(*r).sector = (*r).prev
}
func (r *ring) incDigit() {
	for i := 0; i < 3; i++ {
		(*r).advanceRow()
	}
}
func (r *ring) decDigit() {
	for i := 0; i < 3; i++ {
		(*r).backupRow()
	}
}
func (r *ring) nextDigitDelayed() {
	for i := 0; i < 3; i++ {
		(*r).advanceRow()
		time.Sleep(100 * time.Millisecond)
	}
}
func (r *ring) prevDigitDelayed() {
	for i := 0; i < 3; i++ {
		(*r).backupRow()
		time.Sleep(100 * time.Millisecond)
	}
}

// builds a ring of rows
func newRing() *ring {
	d := newDigits()
	r := &ring{&sector{
		pos: 0,
		str: d[0],
	}}
	pred := r.sector

	for i := 1; i < 30; i++ {
		i := i
		// allocate new sector
		s := &sector{
			pos: i,
			str: d[i],
			// set predescessor as previous segment
			prev: pred,
		}
		// set predescessors next to point to new element
		(*pred).next = s

		if i == 29 {
			(*s).next = (*r).sector
			(*r.sector).prev = s
		}

		// let predescessor be the new element
		pred = s
	}
	// center ring on first element
	(*r).backupRow()

	// buffered de-/increase runs in the background, so that inc and dec
	// calls don't block (until the buffer is full) and
	// next/prevDigitAnimated calls run sequentially

	return r
}

type digit struct {
	*ring
	hasNext bool
	next    *digit
}

func (d *digit) stringRows() []string {
	var rows []string
	if (*d).hasNext {
		rows = strings.Split(d.next.String(), "\n")
		append := strings.Split((*d).stringDig(), "\n")
		for i := 0; i < 3; i++ {
			i := i
			rows[i] = rows[i] + " " + append[i] + "\n"
		}
	} else {
		rows = strings.Split(d.stringDig(), "\n")
	}
	return rows
}
func (d *digit) String() string {
	var rows []string
	var str string
	if (*d).hasNext {
		rows = strings.Split(d.next.String(), "\n")
		append := strings.Split((*d).stringDig(), "\n")
		for i := 0; i < 3; i++ {
			i := i
			rows[i] = rows[i] + " " + append[i] + "\n"
			str = str + rows[i]
		}
	} else {
		str = d.stringDig()
	}
	return str
}
func (d *digit) addDigit() {
	if !(*d).hasNext {
		(*d).next = newDigit()
		(*d).hasNext = true
	} else {
		(*d.next).addDigit()
	}
}
func (d *digit) increase() {
	(*d).incDigit()
	if (*d).pos <= 2 { // overflow
		if !(*d).hasNext {
			(*d).addDigit()
		}
		(*d).next.incDigit()
	}
}
func (d *digit) decrease() {
	(*d).decDigit()
	if (*d).pos >= 26 { // underflow
		(*d.next).decDigit()
	}
}
func newDigit() *digit {
	r := newRing()
	(*r).incDigit()
	return &digit{r, false, nil}
}
