package mono

import (
	"errors"
	"math/rand"
)

// CellGraph is a graph of rooms connected by edges.
// Conservation holds by construction.
type CellGraph struct {
	Rooms map[string]*Room
	Edges map[string]map[string]bool
}

// NewCellGraph creates an empty CellGraph.
func NewCellGraph() *CellGraph {
	return &CellGraph{
		Rooms: make(map[string]*Room),
		Edges: make(map[string]map[string]bool),
	}
}

// AddRoom adds a room to the graph.
func (g *CellGraph) AddRoom(room *Room) {
	g.Rooms[room.ID] = room
	if _, ok := g.Edges[room.ID]; !ok {
		g.Edges[room.ID] = make(map[string]bool)
	}
}

// AddEdge connects two rooms.
func (g *CellGraph) AddEdge(a, b string) error {
	if _, ok := g.Rooms[a]; !ok {
		return errors.New("room " + a + " not in graph")
	}
	if _, ok := g.Rooms[b]; !ok {
		return errors.New("room " + b + " not in graph")
	}
	g.Edges[a][b] = true
	g.Edges[b][a] = true
	return nil
}

// Neighbours returns the neighbour IDs of a room.
func (g *CellGraph) Neighbours(id string) []string {
	nbrs := g.Edges[id]
	out := make([]string, 0, len(nbrs))
	for k := range nbrs {
		out = append(out, k)
	}
	return out
}

// Tick runs one simulation tick.
func (g *CellGraph) Tick() {
	for _, room := range g.Rooms {
		predicted := room.Predict()
		perturbation := (rand.Float64() - 0.5) * 0.02
		room.Observe(predicted + perturbation)
	}
}

// Diffuse redistributes vibe between connected rooms (conservation by construction).
func (g *CellGraph) Diffuse(rate float64) {
	deltas := make(map[string]float64)
	for id := range g.Rooms {
		deltas[id] = 0
	}
	for rid, room := range g.Rooms {
		nbrs := g.Neighbours(rid)
		if len(nbrs) == 0 {
			continue
		}
		share := rate * room.Vibe / float64(len(nbrs))
		deltas[rid] -= rate * room.Vibe
		for _, n := range nbrs {
			deltas[n] += share
		}
	}
	for rid, delta := range deltas {
		room := g.Rooms[rid]
		room.Observe(room.Vibe + delta)
	}
}

// Gossip shares vibe between rooms.
func (g *CellGraph) Gossip(roomID string, targetID *string) (float64, error) {
	if _, ok := g.Rooms[roomID]; !ok {
		return 0, errors.New("room " + roomID + " not in graph")
	}
	nbrs := g.Neighbours(roomID)
	if len(nbrs) == 0 {
		return g.Rooms[roomID].Vibe, nil
	}
	target := ""
	if targetID != nil {
		for _, n := range nbrs {
			if n == *targetID {
				target = n
				break
			}
		}
	}
	if target == "" {
		target = nbrs[rand.Intn(len(nbrs))]
	}
	shared := g.Rooms[roomID].Vibe
	avg := (shared + g.Rooms[target].Vibe) / 2
	g.Rooms[roomID].Observe(avg)
	g.Rooms[target].Observe(avg)
	return shared, nil
}

// Learn runs one learning step.
func (g *CellGraph) Learn() {
	for _, room := range g.Rooms {
		predicted := room.Predict()
		room.Observe(0.9*room.Vibe + 0.1*predicted)
	}
}

// TotalVibe returns the sum of all room vibes.
func (g *CellGraph) TotalVibe() float64 {
	total := 0.0
	for _, room := range g.Rooms {
		total += room.Vibe
	}
	return total
}

// Size returns the number of rooms.
func (g *CellGraph) Size() int {
	return len(g.Rooms)
}
