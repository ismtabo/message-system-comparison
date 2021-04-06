package measurement

import (
	"time"
)

type Measurements struct {
	n              uint
	sumX           float64
	sumX2          float64
	start          time.Time
	end            time.Time
	firstReceiving time.Time
	lastReceiving  time.Time
}

func (m *Measurements) AddMeasurement(t time.Duration) {
	x := t.Seconds()
	m.n++
	m.sumX += x
	m.sumX2 += x * x
}

func (m *Measurements) Mean() float64 {
	return m.sumX / float64(m.n)
}

func (m *Measurements) Var() float64 {
	fN := float64(m.n)
	return (m.sumX2 - (m.sumX*m.sumX)/fN) / (fN - 1)
}

func (m *Measurements) N() uint {
	return m.n
}

func (m *Measurements) Start() {
	m.start = time.Now()
}

func (m *Measurements) Stop() {
	m.end = time.Now()
}

func (m *Measurements) Elapsed() time.Duration {
	return m.end.Sub(m.start)
}

func (m *Measurements) SetFirstReceiving(t time.Time) {
	m.firstReceiving = t
}

func (m *Measurements) RegisterLastReceiving() {
	m.lastReceiving = time.Now()
}

func (m *Measurements) ElapsedMessagingTime() time.Duration {
	return m.lastReceiving.Sub(m.firstReceiving)
}
