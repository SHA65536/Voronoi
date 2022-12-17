package voronoi

import (
	"math"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Animation struct {
	Points     []*Point
	ShowPoints bool
}

func MakeAnimation(numPoints int, showPoints bool) (*Animation, error) {
	var res = &Animation{
		Points:     make([]*Point, numPoints),
		ShowPoints: showPoints,
	}
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	r := *rand.New(rand.NewSource(time.Now().UnixNano()))
	w, h := termbox.Size()
	for i := range res.Points {
		// Spreading initial dots
		res.Points[i] = &Point{
			Y:     math.Mod(r.Float64(), float64(h)),
			X:     math.Mod(r.Float64(), float64(w)),
			DY:    math.Mod(r.Float64(), 4),
			DX:    math.Mod(r.Float64(), 4),
			Color: (termbox.Attribute(i) % 8) + 1,
		}
		// Movement can be backwards
		if r.Int()%2 == 0 {
			res.Points[i].DX *= -1
		}
		if r.Int()%2 == 0 {
			res.Points[i].DY *= -1
		}
	}
	return res, nil
}

// Start starts the animation
func (a *Animation) Start() {
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	a.Draw()
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
				break loop
			}
		default:
			a.Draw()
			time.Sleep(60 * time.Millisecond)
		}
	}
}

// Draw draws the animation on screen and updates the dots
func (a *Animation) Draw() {
	w, h := termbox.Size()
	// Updating
	for i := range a.Points {
		a.Points[i].Update(h, w)
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// Drawing
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var minDist float64 = math.MaxFloat64
			var curColor termbox.Attribute
			// Finding closest point
			for i := range a.Points {
				d := a.Points[i].Distance(y, x)
				if d < minDist {
					minDist = d
					curColor = a.Points[i].Color
				}
			}
			// Drawing it's colour
			if a.ShowPoints && minDist <= 0.5 {
				termbox.SetCell(x, y, 'x', termbox.ColorDefault, curColor)
			} else {
				termbox.SetCell(x, y, ' ', termbox.ColorDefault, curColor)
			}
		}
	}
	termbox.Flush()
}

type Point struct {
	Y, X   float64           // Actual position
	DY, DX float64           // Movement Diff
	Color  termbox.Attribute // Color of the point
}

// Update updates the position of the point
// and wraps it around the width and height
func (p *Point) Update(h, w int) {
	// Updating position
	p.Y += p.DY
	p.X += p.DX
	// Wrapping
	p.Y = math.Mod(p.Y, float64(h))
	if p.Y < 0 {
		p.Y += float64(h)
	}
	p.X = math.Mod(p.X, float64(w))
	if p.X < 0 {
		p.X += float64(w)
	}
}

// Distance calculates the distance of (y,x) from the point
func (p *Point) Distance(y, x int) float64 {
	return math.Sqrt(math.Pow(p.X-float64(x), 2) + math.Pow(p.Y-float64(y), 2))
}
