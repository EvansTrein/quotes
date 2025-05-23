package service

import (
	"context"
	"log/slog"
	"math/rand"
	"quotes/internal/entity"
	"quotes/internal/repository"
	customErr "quotes/pkg/error"
)

type QuoteService struct {
	log *slog.Logger
}

type QuoteServiceDeps struct {
	*slog.Logger
}

func NewQuoteService(deps *QuoteServiceDeps) *QuoteService {
	return &QuoteService{log: deps.Logger}
}

func (q *QuoteService) Add(ctx context.Context, data *entity.Quote) error {
	op := "quote service: adding"
	log := q.log.With(slog.String("operation", op))
	log.Debug("QuoteService call func Add", "data", data)

	db := repository.GetQuoteRepo()

	newID := repository.GetNextID()
	data.ID = newID

	db.Store(newID, data)
	repository.SetCountIncrement()

	log.Info("successfully created")
	return nil
}

func (q *QuoteService) GetAll(ctx context.Context, author string) ([]*entity.Quote, error) {
	op := "quote service: getting"
	log := q.log.With(slog.String("operation", op))
	log.Debug("QuoteService call func GetAll", "author", author)

	data := make([]*entity.Quote, 0)
	db := repository.GetQuoteRepo()

	db.Range(func(key, value any) bool {
		quote, ok := value.(*entity.Quote)
		if !ok {
			q.log.Error("unexpected storage format", "key", key, "value", value)
			return false
		}

		if author == "" || quote.Author == author {
			data = append(data, quote)
		}

		return true
	})

	if author != "" && len(data) == 0 {
		log.Warn("no author record found", "author", author)
		return nil, customErr.ErrAuthorNotFound
	}

	log.Info("successfully received")
	return data, nil
}

func (q *QuoteService) GetRandom(ctx context.Context) (*entity.Quote, error) {
	op := "quote service: getting random"
	log := q.log.With(slog.String("operation", op))
	log.Debug("QuoteService call func GetRandom")

	db := repository.GetQuoteRepo()
	count := repository.GetCount()

	if count == 0 {
		log.Warn("no quotes available")
		return nil, customErr.ErrNoQuotesAvailable
	}

	var quotes []*entity.Quote
	db.Range(func(key, value any) bool {
		quote, ok := value.(*entity.Quote)
		if !ok {
			q.log.Error("unexpected storage format", "key", key, "value", value)
			return false
		}
		quotes = append(quotes, quote)
		return true
	})

	result := quotes[rand.Intn(len(quotes))]

	log.Info("successfully received random quote", "id", result.ID)
	return result, nil
}

func (q *QuoteService) DeleteByID(ctx context.Context, id uint32) error {
	op := "quote service: deleting"
	log := q.log.With(slog.String("operation", op))
	log.Debug("QuoteService call func DeleteByID", "id", id)

	db := repository.GetQuoteRepo()

	if _, exists := db.Load(id); !exists {
		log.Warn("quote not found", "id", id)
		return customErr.ErrRecordNotFound
	}

	db.Delete(id)
	repository.SetCountDecrement()

	log.Info("successfully deleted")
	return nil
}
