package goptcha

import (
	"crypto/rand"
	"image"
	"image/color"
	"math/big"

	"github.com/fogleman/gg"
)

type ClockColor struct {
	Background color.Color
	Border     color.Color

	Hour   color.Color
	Minute color.Color

	MinorTick color.Color
	MajorTick color.Color
}

func (c *ClockColor) fillDefaults() {
	if c.Background == nil {
		c.Background = clockDefaultColor.Background
	}
	if c.Border == nil {
		c.Border = clockDefaultColor.Border
	}
	if c.Hour == nil {
		c.Hour = clockDefaultColor.Hour
	}
	if c.Minute == nil {
		c.Minute = clockDefaultColor.Minute
	}
	if c.MinorTick == nil {
		c.MinorTick = clockDefaultColor.MinorTick
	}
	if c.MajorTick == nil {
		c.MajorTick = clockDefaultColor.MajorTick
	}
}

type ClockTime struct {
	Hour   int
	Minute int
}

type ClockConfig struct {
	Time  *ClockTime
	Color *ClockColor
}

var clockDefaultColor = ClockColor{
	Background: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	Border:     color.RGBA{R: 0, G: 0, B: 0, A: 255},

	Hour:   color.RGBA{R: 0, G: 0, B: 255, A: 255},
	Minute: color.RGBA{R: 255, G: 0, B: 255, A: 255},

	MinorTick: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	MajorTick: color.RGBA{R: 255, G: 0, B: 0, A: 255},
}

type Clock struct {
	Time  ClockTime
	Image image.Image
}

func NewClock(size int) (Clock, error) {
	return NewClockWithConfig(size, ClockConfig{
		Time:  nil,
		Color: nil,
	})
}

func NewClockWithConfig(size int, config ClockConfig) (Clock, error) {
	if config.Time == nil {
		bigint, err := rand.Int(rand.Reader, big.NewInt(12))
		if err != nil {
			return Clock{}, err
		}
		h := bigint.Int64()

		bigint, err = rand.Int(rand.Reader, big.NewInt(60))
		if err != nil {
			return Clock{}, err
		}
		m := bigint.Int64()

		config.Time = &ClockTime{Hour: int(h), Minute: int(m)}
	}
	if config.Color == nil {
		config.Color = &clockDefaultColor
	} else {
		config.Color.fillDefaults()
	}

	img, err := generateClock(size, config)
	if err != nil {
		return Clock{}, err
	}

	return Clock{
		Time:  *config.Time,
		Image: img,
	}, nil
}

func generateClock(s int, config ClockConfig) (image.Image, error) {
	ctx := gg.NewContext(s, s)
	cx := float64(s / 2)
	cy := cx
	r := cx
	t := config.Time
	c := config.Color
	h, m := t.Hour, t.Minute

	ctx.SetColor(c.Background)
	ctx.DrawCircle(cx, cy, r)
	ctx.FillPreserve()
	ctx.SetStrokeStyle(gg.NewSolidPattern(c.Border))
	ctx.Stroke()

	r--

	l := r * 0.1
	ctx.SetStrokeStyle(gg.NewSolidPattern(c.MajorTick))
	for i := 0; i < 3; i++ {
		ctx.Push()
		ctx.RotateAbout(gg.Radians(float64(i*30)), cx, cy)
		ctx.DrawLine(cx, cy-r, cx, cy-r+l)
		ctx.DrawLine(cx, cy+r-l, cx, cy+r)
		ctx.DrawLine(cx-r, cy, cx-r+l, cy)
		ctx.DrawLine(cx+r-l, cy, cx+r, cy)
		ctx.Stroke()
		ctx.Pop()
	}

	l = r * 0.05
	ctx.SetStrokeStyle(gg.NewSolidPattern(c.MinorTick))
	for i := 1; i < 15; i++ {
		if i%5 == 0 {
			continue
		}
		ctx.Push()
		ctx.RotateAbout(gg.Radians(float64(i*6)), cx, cy)
		ctx.DrawLine(cx, cy-r, cx, cy-r+l)
		ctx.DrawLine(cx, cy+r-l, cx, cy+r)
		ctx.DrawLine(cx-r, cy, cx-r+l, cy)
		ctx.DrawLine(cx+r-l, cy, cx+r, cy)
		ctx.Stroke()
		ctx.Pop()
	}

	l = r * 0.8
	ctx.SetStrokeStyle(gg.NewSolidPattern(c.Minute))
	ctx.Push()
	ctx.RotateAbout(gg.Radians(float64(m*6)), cx, cy)
	ctx.DrawLine(cx, cy-l, cx, cy)
	ctx.Stroke()
	ctx.Pop()

	l = r * 0.5
	ctx.SetStrokeStyle(gg.NewSolidPattern(c.Hour))
	ctx.Push()
	ctx.RotateAbout(gg.Radians(float64(h*30+m/2)), cx, cy)
	ctx.DrawLine(cx, cy-l, cx, cy)
	ctx.Stroke()
	ctx.Pop()

	return ctx.Image(), nil
}
