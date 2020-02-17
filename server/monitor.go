package server

import (
	"log"
	"time"
)

const (
	numPollers     = 2                // number of Poller goroutines to launch
	pollInterval   = 60 * time.Second // how often to poll each URL
	statusInterval = 10 * time.Second // how often to log status to stdout
	errTimeout     = 10 * time.Second // back-off timeout on error
)

func NewStatusMonitor() *StatusMonitor {
	results := make(chan *Result)
	monitor := &StatusMonitor{
		Results:      nil,
		ResultChan:   results,
		OKResults:    make(map[string]int),
		ErrResults:   make(map[string]int),
	}
	monitor.StartMonitoring()
	return monitor
}

func (m *StatusMonitor) StartMonitoring() {
	ticker := time.NewTicker(statusInterval)
	start := time.Now()
	go func() {
		totalJobsProcessed := 0

		for {
			select {

			case <-ticker.C:
				secSince := time.Since(start)
				rate := totalJobsProcessed / int(secSince.Seconds())
				log.Printf("STAT: jobs: %v, avg: %v\n", totalJobsProcessed, rate)

			case r := <- m.ResultChan:
				totalJobsProcessed++
				//log.Println("RESULT:", r.ConvertedHTML)
				log.Println("RESULT:", r.Id)

				/*
				if r.Errors != nil {
					for _, error := range r.Errors {
						log.Println("ERROR:", error)
					}
					log.Println("RESULT:", totalPageCrawled, r.Page.Board, r.Page.Num, len(r.Page.Posts), "ERROR")
					m.RLock()
					value, ok := m.ErrResults[r.Page.Board]
					m.RUnlock()
					if !ok {
						m.Lock()
						m.ErrResults[r.Page.Board] = 0
						m.Unlock()
					}
					m.Lock()
					m.ErrResults[r.Page.Board] = value + 1
					m.Unlock()
				} else {
					log.Println("RESULT:", totalPageCrawled, r.Page.Board, r.Page.Num, len(r.Page.Posts), "OK")
					m.RLock()
					value, ok := m.OKResults[r.Page.Board]
					m.RUnlock()
					if !ok {
						m.Lock()
						m.OKResults[r.Page.Board] = 0
						m.Unlock()
					}
					m.Lock()
					m.OKResults[r.Page.Board] = value + 1
					m.Unlock()
				}
				 */
			}
		}
	}()
}

func (m *StatusMonitor) IsCompleted(name string, expected int) bool {
	if m.OKResults[name] >= expected {
		return true
	}
	return false
}
