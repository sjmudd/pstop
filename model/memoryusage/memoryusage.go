// Package memoryusage manages collecting data from performance_schema which holds
// information about memory usage
package memoryusage

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // keep golint happy

	"github.com/sjmudd/ps-top/baseobject"
	"github.com/sjmudd/ps-top/config"
)

// MemoryUsage represents a table of rows
type MemoryUsage struct {
	baseobject.BaseObject      // embedded
	last                  Rows // last loaded values
	Results               Rows // results (maybe with subtraction)
	Totals                Row  // totals of results
	db                    *sql.DB
}

// NewMemoryUsage returns a pointer to a MemoryUsage struct
func NewMemoryUsage(cfg *config.Config, db *sql.DB) *MemoryUsage {
	mu := &MemoryUsage{
		db: db,
	}
	mu.SetConfig(cfg)

	return mu
}

// Collect data from the db, no merging needed
// DEPRECATED
func (mu *MemoryUsage) Collect() {
	mu.AddRows(collect(mu.db))
}

// add new rows to the dataset
func (mu *MemoryUsage) AddRows(rows Rows) {
	mu.last = rows
	mu.LastCollected = time.Now()

	mu.calculate()
}

// ResetStatistics resets the statistics to current values
func (mu *MemoryUsage) ResetStatistics() {

	mu.calculate()
}

// Rows returns the rows we have which are interesting
func (mu MemoryUsage) Rows() []Row {
	rows := make([]Row, 0, len(mu.Results))

	for i := range mu.Results {
		rows = append(rows, mu.Results[i])
	}

	return rows
}

// HaveRelativeStats returns if the values returned are relative to a previous collection
func (mu MemoryUsage) HaveRelativeStats() bool {
	return false
}

func (mu *MemoryUsage) calculate() {
	mu.Results = make(Rows, len(mu.last))
	copy(mu.Results, mu.last)
	mu.Totals = totals(mu.Results)
}
