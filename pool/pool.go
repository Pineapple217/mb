package pool

import (
	"github.com/gomarkdown/markdown/html"
)

type Pool struct {
	// buffered channel for connection pooling
	c chan *html.Renderer

	// factory to create new connection
	f func() (*html.Renderer, error)
}

func NewPool(s int, f func() (*html.Renderer, error)) *Pool {
	return &Pool{
		c: make(chan *html.Renderer, s),
		f: f,
	}
}

// get one idle conn from pool, if pool empty, create a new one
func (p *Pool) Get() (*html.Renderer, error) {
	select {
	case c := <-p.c:
		// log.Printf("Got connection from pool: %p", c)
		return c, nil

	default:
		// log.Printf("Creating new connection")
		c, _ := p.f()
		// if err != nil {
		// 	log.Printf("Failed to create new connection: %v", err)
		// 	return *html.Renderer{}, err
		// }

		return c, nil
	}
}

// put conn into pool, if the pool full, close the conn instead
func (p *Pool) Put(r *html.Renderer) {
	select {
	case p.c <- r:
		// log.Printf("Connection idle, joined pool: %p", r)

	default:
		// log.Printf("Pool full: closing current connection: %p", r)
		// r.Close()
	}
}
