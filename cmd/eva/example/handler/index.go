
package handler

import (
	"context"
	"example/proto/example"
)

type HandlerService struct {

}

func (m *HandlerService) Ping(context.Context, *example.Nil) (*example.Nil, error) {
	return &example.Nil{}, nil
}

