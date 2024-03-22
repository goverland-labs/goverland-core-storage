package stats

type Dao struct {
	Total         int64
	TotalVerified int64
}

type Proposals struct {
	Total int64
}

type Totals struct {
	Dao       Dao
	Proposals Proposals
}
