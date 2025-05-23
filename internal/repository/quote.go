package repository

import (
	"sync"
	"sync/atomic"
)

var db *QuoteRepo

type QuoteRepo struct {
	*sync.Map
	idNumder uint32
	count    uint32
}

func InitQuoteRepo() {
	db = &QuoteRepo{Map: &sync.Map{}}
}

func GetQuoteRepo() *QuoteRepo {
	return db
}

func GetNextID() uint32 {
	return atomic.AddUint32(&db.idNumder, 1)
}

func GetCount() uint32 {
	return db.count
}

func SetCountIncrement() {
	atomic.AddUint32(&db.count, 1)
}

func SetCountDecrement() {
	atomic.AddUint32(&db.count, ^uint32(0))
}
