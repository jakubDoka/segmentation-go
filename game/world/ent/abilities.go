package ent

import "github.com/faiface/pixel"

type Updatable interface {
	Update()
}

type Mortal interface {
	Kill()
}

type Notifiable interface {
	Notify(id int)
}

type Drawable interface {
	Draw(pixel.Target)
}

type Movable interface {
	SetPos(pos pixel.Vec)
}

type Rotatable interface {
	SetRot(rot float64)
}

type Hitable interface {
	Demage(int)
}

type Controlable interface {
	Set(bool)
}
