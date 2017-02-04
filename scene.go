package main

import (
	"context"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	renderer *sdl.Renderer
	bg       *sdl.Texture
	bgt      int32
	bird     *bird
}

func newScene(r *sdl.Renderer, speed int32, gravity float64) (s *scene, err error) {
	s = &scene{renderer: r}

	s.bg, err = img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	s.bird = &bird{
		x:       10,
		y:       windowHeight / 2,
		w:       50,
		h:       43,
		gravity: gravity,
	}
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		frame, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird_frame_%d image: %v", i, err)
		}
		s.bird.frames = append(s.bird.frames, frame)
	}

	return s, nil
}

func (s *scene) restart() {
	s.bird.y = windowHeight / 2
	s.bird.dead = false
	s.bird.speed = 0
}

func (s *scene) run(ctx context.Context, fps float64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !s.bird.dead {
				s.update()
				s.draw()
			}

			sdl.Delay(uint32(1000 / fps))
		}
	}
}

func (s *scene) update() {
	s.bird.update()
	s.bgt = (s.bgt + 1) % 2000
}

func (s *scene) draw() error {
	s.renderer.Clear()

	bgRect := &sdl.Rect{X: s.bgt, Y: 0, W: windowWidth, H: windowHeight}
	err := s.renderer.Copy(s.bg, bgRect, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	s.bird.draw(s.renderer)

	if s.bird.dead {
		deadRect := &sdl.Rect{X: 100, Y: windowHeight / 4, W: windowWidth - 200, H: windowHeight / 2}
		drawText(s.renderer, "YOU DIED!", deadRect, sdl.Color{R: 255})
	}

	s.renderer.Present()
	return nil
}
