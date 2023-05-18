package events

import "gorm.io/gorm"

type RegisteredEvent struct {
	gorm.Model

	Type   string // dao/proposal etc
	TypeID string // type identifier: dao.id, proposal.id, etc
	Event  string // core.proposal.created, etc
}
