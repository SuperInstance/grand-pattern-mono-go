package mono

import (
	"errors"
	"math/rand"
)

func makeRooms(n int, alpha float64) ([]*Room, error) {
	rooms := make([]*Room, n)
	for i := 0; i < n; i++ {
		r, err := NewRoom(roomID(i), 0, alpha)
		if err != nil {
			return nil, err
		}
		rooms[i] = r
	}
	return rooms, nil
}

func roomID(i int) string {
	return "r" + itoa(i)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	digits := []byte{}
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}

// Chain creates a linear chain of n rooms.
func Chain(n int, alpha float64) (*CellGraph, error) {
	g := NewCellGraph()
	rooms, err := makeRooms(n, alpha)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		g.AddRoom(r)
	}
	for i := 0; i < n-1; i++ {
		g.AddEdge(roomID(i), roomID(i+1))
	}
	return g, nil
}

// Ring creates a ring of n rooms.
func Ring(n int, alpha float64) (*CellGraph, error) {
	g, err := Chain(n, alpha)
	if err != nil {
		return nil, err
	}
	if n > 2 {
		g.AddEdge("r0", roomID(n-1))
	}
	return g, nil
}

// Star creates a star with room 0 connected to all others.
func Star(n int, alpha float64) (*CellGraph, error) {
	g := NewCellGraph()
	rooms, err := makeRooms(n, alpha)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		g.AddRoom(r)
	}
	for i := 1; i < n; i++ {
		g.AddEdge("r0", roomID(i))
	}
	return g, nil
}

// Mesh creates a fully connected mesh.
func Mesh(n int, alpha float64) (*CellGraph, error) {
	g := NewCellGraph()
	rooms, err := makeRooms(n, alpha)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		g.AddRoom(r)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			g.AddEdge(roomID(i), roomID(j))
		}
	}
	return g, nil
}

// SmallWorld creates a Watts-Strogatz small-world graph.
func SmallWorld(n, k int, p, alpha float64) (*CellGraph, error) {
	if k%2 != 0 {
		return nil, errors.New("k must be even")
	}
	g := NewCellGraph()
	rooms, err := makeRooms(n, alpha)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		g.AddRoom(r)
	}

	type edge struct{ a, b int }
	edgeSet := make(map[edge]bool)
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			a, b := i, (i+j)%n
			if a > b {
				a, b = b, a
			}
			edgeSet[edge{a, b}] = true
		}
	}
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			if rand.Float64() < p {
				oldA, oldB := i, (i+j)%n
				if oldA > oldB {
					oldA, oldB = oldB, oldA
				}
				delete(edgeSet, edge{oldA, oldB})
				newB := rand.Intn(n - 1)
				if newB >= i {
					newB++
				}
				a, b := i, newB
				if a > b {
					a, b = b, a
				}
				edgeSet[edge{a, b}] = true
			}
		}
	}
	for e := range edgeSet {
		g.AddEdge(roomID(e.a), roomID(e.b))
	}
	return g, nil
}

// ScaleFree creates a Barabási-Albert scale-free graph.
func ScaleFree(n, m int, alpha float64) (*CellGraph, error) {
	if m < 1 || m >= n {
		return nil, errors.New("m must be in [1, n-1]")
	}
	g := NewCellGraph()
	rooms, err := makeRooms(n, alpha)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		g.AddRoom(r)
	}

	repeated := []int{}
	for i := 0; i <= m; i++ {
		for j := i + 1; j <= m; j++ {
			g.AddEdge(roomID(i), roomID(j))
			repeated = append(repeated, i, j)
		}
	}
	for newN := m + 1; newN < n; newN++ {
		targets := map[int]bool{}
		for len(targets) < m {
			targets[repeated[rand.Intn(len(repeated))]] = true
		}
		for t := range targets {
			g.AddEdge(roomID(newN), roomID(t))
			repeated = append(repeated, newN, t)
		}
	}
	return g, nil
}
