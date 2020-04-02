package database

import "time"

type Votes struct {
	Id uint64 `gorm:"type:bigint;unique_index;primary_key"`
	VoteTime time.Time
}

func (Votes) TableName() string {
	return "votes"
}

func HasVoted(owner uint64, ch chan bool) {
	var node Votes
	Db.Where(Votes{Id: owner}).First(&node)

	ch <- time.Now().Sub(node.VoteTime) < 24 * time.Hour
}
