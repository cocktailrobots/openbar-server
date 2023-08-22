package dbutils

import "sync"

type DBBranchTracker struct {
	mu     *sync.Mutex
	branch string
}

func NewDBBranchTracker() *DBBranchTracker {
	return &DBBranchTracker{
		mu:     &sync.Mutex{},
		branch: "main",
	}
}

func (dbbt DBBranchTracker) SetBranch(branch string) {
	dbbt.mu.Lock()
	defer dbbt.mu.Unlock()

	dbbt.branch = branch
}

func (dbbt DBBranchTracker) GetBranch() string {
	dbbt.mu.Lock()
	defer dbbt.mu.Unlock()

	return dbbt.branch
}
