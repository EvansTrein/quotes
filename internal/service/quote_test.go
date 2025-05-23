package service

import (
	"context"
	"errors"
	"quotes/internal/entity"
	"quotes/internal/repository"
	customErr "quotes/pkg/error"
	"quotes/pkg/logs"
	"testing"
)

var (
	testData = []*entity.Quote{
		{ID: 1, Author: "Author1", Text: "Quote1"},
		{ID: 2, Author: "Author1", Text: "Quote2"},
		{ID: 3, Author: "Author2", Text: "Quote3"},
		{ID: 4, Author: "Author3", Text: "Quote4"},
	}
)

func TestQuoteService_Add(t *testing.T) {
	log := logs.NewDiscardLogger()
	repository.InitQuoteRepo()

	service := &QuoteService{log: log}

	type args struct {
		ctx  context.Context
		data *entity.Quote
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "successful add",
			args: args{
				ctx: context.Background(),
				data: &entity.Quote{
					Author: "Test Author",
					Text:   "Test Quote",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Add(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("QuoteService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuoteService_GetAll(t *testing.T) {
	log := logs.NewDiscardLogger()
	repository.InitQuoteRepo()

	repo := repository.GetQuoteRepo()
	for _, q := range testData {
		repository.SetCountIncrement()
		repo.Store(q.ID, q)
	}

	service := &QuoteService{log: log}

	type args struct {
		ctx    context.Context
		author string
	}
	tests := []struct {
		wantErr   error
		args      args
		name      string
		wantCount int
	}{
		{
			name: "get all quotes",
			args: args{
				ctx:    context.Background(),
				author: "",
			},
			wantCount: 4,
			wantErr:   nil,
		},
		{
			name: "get quotes by existing author",
			args: args{
				ctx:    context.Background(),
				author: "Author1",
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name: "get quotes by non-existing author",
			args: args{
				ctx:    context.Background(),
				author: "UnknownAuthor",
			},
			wantCount: 0,
			wantErr:   customErr.ErrAuthorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetAll(tt.args.ctx, tt.args.author)

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("QuoteService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("QuoteService.GetAll() count = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestQuoteService_GetRandom(t *testing.T) {
	log := logs.NewDiscardLogger()
	repository.InitQuoteRepo()

	repo := repository.GetQuoteRepo()
	for _, q := range testData {
		repository.SetCountIncrement()
		repo.Store(q.ID, q)
	}

	service := &QuoteService{log: log}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		args    args
		wantErr error
		name    string
	}{
		{
			name:    "successful get random quote",
			args:    args{ctx: context.Background()},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := service.GetRandom(tt.args.ctx)

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("QuoteService.GetRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Error("expected non-nil quote, got nil")
				return
			}

			if got.ID == 0 {
				t.Error("quote ID should not be 0")
				return
			}

			if got.Author == "" {
				t.Error("quote Author should not be empty")
				return
			}

			if got.Text == "" {
				t.Error("quote Text should not be empty")
				return
			}

		})
	}
}

func TestQuoteService_DeleteByID(t *testing.T) {
	log := logs.NewDiscardLogger()
	repository.InitQuoteRepo()

	repo := repository.GetQuoteRepo()
	for _, q := range testData {
		repository.SetCountIncrement()
		repo.Store(q.ID, q)
	}

	service := &QuoteService{log: log}

	type args struct {
		ctx context.Context
		id  uint32
	}
	tests := []struct {
		wantErr error
		args    args
		name    string
	}{
		{
			name:    "successful delete",
			args:    args{ctx: context.Background(), id: 2},
			wantErr: nil,
		},
		{
			name:    "delete non-existent quote",
			args:    args{ctx: context.Background(), id: 999},
			wantErr: customErr.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteByID(tt.args.ctx, tt.args.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("QuoteService.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				_, exists := repo.Load(tt.args.id)
				if exists {
					t.Errorf("expected quote with ID %d to be deleted, but it still exists", tt.args.id)
				}
			}
		})
	}
}
