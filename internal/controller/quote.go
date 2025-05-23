package controller

import (
	"context"
	"fmt"
	"net/http"
	"quotes/internal/entity"
	customErr "quotes/pkg/error"
	"quotes/pkg/utils"
	"strconv"
)

type IQuoteService interface {
	Add(ctx context.Context, data *entity.Quote) error
	GetAll(ctx context.Context, author string) ([]*entity.Quote, error)
	GetRandom(ctx context.Context) (*entity.Quote, error)
	DeleteByID(ctx context.Context, id uint32) error
}

type QuotesController struct {
	baseContrl   *BaseController
	quoteService IQuoteService
}

type QuotesControllerDeps struct {
	Router *http.ServeMux
	*BaseController
	IQuoteService
}

func NewQuotesController(deps *QuotesControllerDeps) *QuotesController {
	c := &QuotesController{baseContrl: deps.BaseController, quoteService: deps.IQuoteService}
	c.initRouter(deps.Router)
	return c
}

func (c *QuotesController) initRouter(r *http.ServeMux) {
	group := "quotes"

	r.Handle(fmt.Sprintf("POST /%s", group), c.Add())
	r.Handle(fmt.Sprintf("GET /%s", group), c.GetAll())
	r.Handle(fmt.Sprintf("GET /%s/random", group), c.GetRandom())
	r.Handle(fmt.Sprintf("DELETE /%s/{id}", group), c.DeleteByID())
}

func (c *QuotesController) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.DecodeBody[entity.Quote](r.Body)
		if err != nil {
			c.baseContrl.HandleError(w, fmt.Errorf("failed to read the request body, err: %v", err))
			return
		}

		if data.Author == "" || data.Text == "" {
			c.baseContrl.HandleError(w, customErr.ErrNoFields)
			return
		}

		if err := c.quoteService.Add(r.Context(), data); err != nil {
			c.baseContrl.HandleError(w, err)
			return
		}

		w.WriteHeader(201)
	}
}

func (c *QuotesController) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")

		resp, err := c.quoteService.GetAll(r.Context(), author)
		if err != nil {
			c.baseContrl.HandleError(w, err)
			return
		}

		c.baseContrl.SendJsonResp(w, 200, resp)
	}
}

func (c *QuotesController) GetRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resp, err := c.quoteService.GetRandom(r.Context())
		if err != nil {
			c.baseContrl.HandleError(w, err)
			return
		}

		c.baseContrl.SendJsonResp(w, 200, resp)
	}
}

func (c *QuotesController) DeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			c.baseContrl.HandleError(w, customErr.ErrInvalidTypeID)
			return
		}

		if err := c.quoteService.DeleteByID(r.Context(), uint32(id)); err != nil {
			c.baseContrl.HandleError(w, err)
			return
		}

		w.WriteHeader(204)
	}
}
