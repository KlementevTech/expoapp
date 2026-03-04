package model

import (
	"time"

	"github.com/google/uuid"
)

const firstVersion = 1

// PartID Должен быть UUID V7.
type PartID = uuid.UUID

type Part struct {
	ID        PartID
	Name      string
	Version   int
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewPart(id PartID, name string) *Part {
	return &Part{
		ID:        id,
		Name:      name,
		Version:   firstVersion,
		CreatedAt: timeNow(),
	}
}

func (p *Part) Update(name string) {
	p.UpdatedAt = new(timeNow())
	p.Version++
	p.Name = name
}

func (p *Part) Delete() {
	p.DeletedAt = new(timeNow())
}

func timeNow() time.Time {
	return time.Now().UTC()
}
