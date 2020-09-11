package ent

import "github.com/faiface/pixel"

type Existence interface {
	GetPos() pixel.Vec
}

type Velocity interface {
	GetVel() pixel.Vec
}

type Speed interface {
	GetSpeed() float64
}

type Volume interface {
	GetSize() float64
}

type Averness interface {
	GetSight() float64
}

type Vitality interface {
	IsAlive() bool
}

type State interface {
	IsActive() bool
}

type Stamina interface {
	GetHealth() int
}

type Passable interface {
	Passable(team uint16) bool
}

type Strength interface {
	GetDamage() int
}

type Competence interface {
	GetTeam() uint16
}

type Agility interface {
	GetSteer() float64
}

type Limit interface {
	GetMax() int
}

type Rot interface {
	GetRot() float64
}

type Finite interface {
	IsDead() bool
}

type Origin interface {
	GetID() uint16
}

type SelfControl interface {
	IsControlled() bool
}
