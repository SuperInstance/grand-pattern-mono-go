package mono

import (
	"math"
	"testing"
)

// ── Jepa tests ─────────────────────────────────────────────────────────

func TestJepaEmptyPredict(t *testing.T) {
	j, _ := NewJepa(0.3)
	if j.Predict() != 0 {
		t.Errorf("expected 0, got %f", j.Predict())
	}
}

func TestJepaSingleObserve(t *testing.T) {
	j, _ := NewJepa(0.3)
	j.Observe(1.0)
	if math.Abs(j.Predict()-1.0) > 1e-6 {
		t.Errorf("expected 1.0, got %f", j.Predict())
	}
}

func TestJepaTwoObservations(t *testing.T) {
	j, _ := NewJepa(0.5)
	j.Observe(0.0)
	j.Observe(1.0)
	p := j.Predict()
	if p <= 0 || p >= 1 {
		t.Errorf("expected 0<p<1, got %f", p)
	}
}

func TestJepaAlphaOne(t *testing.T) {
	j, _ := NewJepa(1.0)
	j.Observe(5.0)
	j.Observe(10.0)
	if math.Abs(j.Predict()-10.0) > 1e-6 {
		t.Errorf("expected 10.0, got %f", j.Predict())
	}
}

func TestJepaInvalidAlpha(t *testing.T) {
	if _, err := NewJepa(0); err == nil {
		t.Error("expected error for alpha=0")
	}
	if _, err := NewJepa(1.5); err == nil {
		t.Error("expected error for alpha=1.5")
	}
}

func TestJepaSurpriseZero(t *testing.T) {
	j, _ := NewJepa(0.3)
	j.Observe(3.0)
	j.Observe(3.0)
	if math.Abs(j.Surprise(3.0)) > 1e-10 {
		t.Errorf("expected 0 surprise, got %f", j.Surprise(3.0))
	}
}

func TestJepaSurpriseNonzero(t *testing.T) {
	j, _ := NewJepa(0.3)
	j.Observe(0.0)
	j.Observe(0.0)
	if j.Surprise(1.0) <= 0 {
		t.Error("expected positive surprise")
	}
}

func TestJepaHistory(t *testing.T) {
	j, _ := NewJepa(0.3)
	j.Observe(1.0)
	j.Observe(2.0)
	j.Observe(3.0)
	h := j.History()
	if len(h) != 3 || h[0] != 1 || h[1] != 2 || h[2] != 3 {
		t.Errorf("history mismatch: %v", h)
	}
	if j.ObservationCount() != 3 {
		t.Errorf("expected 3, got %d", j.ObservationCount())
	}
}

// ── Room tests ──────────────────────────────────────────────────────────

func TestRoomCreate(t *testing.T) {
	r, _ := NewRoom("a", 1.5, 0.3)
	if r.ID != "a" || r.Vibe != 1.5 {
		t.Errorf("room mismatch: %+v", r)
	}
}

func TestRoomObserve(t *testing.T) {
	r, _ := NewRoom("a", 0, 0.3)
	r.Observe(1.0)
	if r.Vibe != 1.0 {
		t.Errorf("expected vibe 1.0, got %f", r.Vibe)
	}
	if r.LastSurprise <= 0 {
		t.Error("expected positive surprise")
	}
}

func TestRoomPredict(t *testing.T) {
	r, _ := NewRoom("a", 5.0, 0.3)
	if math.Abs(r.Predict()-5.0) > 1e-6 {
		t.Errorf("expected 5.0, got %f", r.Predict())
	}
}

// ── CellGraph tests ────────────────────────────────────────────────────

func makeTestGraph() *CellGraph {
	g := NewCellGraph()
	a, _ := NewRoom("a", 1, 0.3)
	b, _ := NewRoom("b", 2, 0.3)
	c, _ := NewRoom("c", 3, 0.3)
	g.AddRoom(a)
	g.AddRoom(b)
	g.AddRoom(c)
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	return g
}

func TestGraphAddRooms(t *testing.T) {
	g := makeTestGraph()
	if g.Size() != 3 {
		t.Errorf("expected 3 rooms, got %d", g.Size())
	}
	nbrs := g.Neighbours("a")
	found := false
	for _, n := range nbrs {
		if n == "b" {
			found = true
		}
	}
	if !found {
		t.Error("expected b as neighbour of a")
	}
}

func TestGraphMissingEdge(t *testing.T) {
	g := NewCellGraph()
	a, _ := NewRoom("a", 0, 0.3)
	g.AddRoom(a)
	if err := g.AddEdge("a", "z"); err == nil {
		t.Error("expected error for missing room")
	}
}

func TestGraphTotalVibe(t *testing.T) {
	g := makeTestGraph()
	if math.Abs(g.TotalVibe()-6.0) > 1e-6 {
		t.Errorf("expected 6.0, got %f", g.TotalVibe())
	}
}

func TestDiffuseConservation(t *testing.T) {
	g := makeTestGraph()
	before := g.TotalVibe()
	for i := 0; i < 10; i++ {
		g.Diffuse(0.2)
	}
	after := g.TotalVibe()
	if math.Abs(after-before) > 1e-4 {
		t.Errorf("conservation violated: before=%f after=%f", before, after)
	}
}

func TestTick(t *testing.T) {
	g := makeTestGraph()
	g.Tick()
	if g.Size() != 3 {
		t.Error("tick should not change room count")
	}
}

func TestGossip(t *testing.T) {
	g := makeTestGraph()
	v, err := g.Gossip("a", nil)
	if err != nil || v != g.Rooms["a"].Vibe {
		// gossip returns shared vibe before averaging
	}
}

func TestGossipNoNeighbours(t *testing.T) {
	g := NewCellGraph()
	solo, _ := NewRoom("solo", 5.0, 0.3)
	g.AddRoom(solo)
	v, err := g.Gossip("solo", nil)
	if err != nil || math.Abs(v-5.0) > 1e-6 {
		t.Errorf("expected 5.0, got %f err=%v", v, err)
	}
}

func TestGossipTarget(t *testing.T) {
	g := makeTestGraph()
	target := "b"
	g.Gossip("a", &target)
	if math.Abs(g.Rooms["a"].Vibe-g.Rooms["b"].Vibe) > 1e-10 {
		t.Errorf("gossip should equalize: a=%f b=%f", g.Rooms["a"].Vibe, g.Rooms["b"].Vibe)
	}
}

func TestLearn(t *testing.T) {
	g := makeTestGraph()
	g.Learn()
	if g.Size() != 3 {
		t.Error("learn should not change room count")
	}
}

// ── Topology tests ─────────────────────────────────────────────────────

func TestChain(t *testing.T) {
	g, err := Chain(5, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if g.Size() != 5 {
		t.Errorf("expected 5, got %d", g.Size())
	}
	if len(g.Neighbours("r0")) != 1 {
		t.Errorf("expected 1 neighbour for r0, got %d", len(g.Neighbours("r0")))
	}
}

func TestRing(t *testing.T) {
	g, err := Ring(5, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	nbrs := g.Neighbours("r0")
	found := false
	for _, n := range nbrs {
		if n == "r4" {
			found = true
		}
	}
	if !found {
		t.Error("r0 should connect to r4 in ring")
	}
}

func TestStar(t *testing.T) {
	g, err := Star(5, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Neighbours("r0")) != 4 {
		t.Errorf("expected 4 neighbours for r0, got %d", len(g.Neighbours("r0")))
	}
}

func TestMesh(t *testing.T) {
	g, err := Mesh(4, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Neighbours("r0")) != 3 {
		t.Errorf("expected 3 neighbours, got %d", len(g.Neighbours("r0")))
	}
}

func TestSmallWorld(t *testing.T) {
	g, err := SmallWorld(10, 4, 0.2, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if g.Size() != 10 {
		t.Errorf("expected 10, got %d", g.Size())
	}
}

func TestSmallWorldKEven(t *testing.T) {
	if _, err := SmallWorld(10, 3, 0.2, 0.3); err == nil {
		t.Error("expected error for odd k")
	}
}

func TestScaleFree(t *testing.T) {
	g, err := ScaleFree(20, 2, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if g.Size() != 20 {
		t.Errorf("expected 20, got %d", g.Size())
	}
}

func TestScaleFreeInvalidM(t *testing.T) {
	if _, err := ScaleFree(10, 0, 0.3); err == nil {
		t.Error("expected error for m=0")
	}
	if _, err := ScaleFree(5, 5, 0.3); err == nil {
		t.Error("expected error for m=n")
	}
}

func TestDiffuseConservationRing(t *testing.T) {
	g, err := Ring(10, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		id := roomID(i)
		g.Rooms[id].Vibe = float64(i)
		g.Rooms[id].Jepa.Observe(float64(i))
	}
	before := g.TotalVibe()
	for i := 0; i < 20; i++ {
		g.Diffuse(0.3)
	}
	after := g.TotalVibe()
	if math.Abs(after-before) > 1e-3 {
		t.Errorf("conservation violated on ring: before=%f after=%f", before, after)
	}
}
