package mono

// Room is a cell in the mono-vibe architecture.
type Room struct {
	ID           string
	Vibe         float64
	Jepa         *Jepa
	LastSurprise float64
}

// NewRoom creates a room with the given id, initial vibe, and JEPA alpha.
func NewRoom(id string, vibe float64, alpha float64) (*Room, error) {
	j, err := NewJepa(alpha)
	if err != nil {
		return nil, err
	}
	r := &Room{ID: id, Vibe: vibe, Jepa: j}
	r.Jepa.Observe(vibe)
	return r, nil
}

// Observe records a new vibe, updates JEPA, and records surprise.
func (r *Room) Observe(newVibe float64) {
	r.LastSurprise = r.Jepa.Surprise(newVibe)
	r.Vibe = newVibe
	r.Jepa.Observe(newVibe)
}

// Predict returns the JEPA prediction.
func (r *Room) Predict() float64 {
	return r.Jepa.Predict()
}
