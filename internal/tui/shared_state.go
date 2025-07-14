package tui

import (
	"github.com/bgics/pmjay-go/model"
	"github.com/bgics/pmjay-go/store"
)

type SharedState struct {
	SelectedRecord model.FormData
	LastPageIndex  PageIndex
	Store          *store.Store
	Error          error
}

func NewSharedState() *SharedState {
	s := &SharedState{}
	s.Store = store.NewStore()

	return s
}
