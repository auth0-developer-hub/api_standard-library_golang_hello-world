package server

import (
	"hello-golang-api/entities"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/snowzach/queryp"
)

// MessageSave saves a message
//
// @ID MessageSave
// @Tags messages
// @Summary Save a message
// @Description Save a message
// @Param message body entities.Message true "Message"
// @Success 200 {object} entities.Message
// @Failure 400 string true "Invalid Argument"
// @Failure 500 string true "Internal Error"
// @Router /messages [post]
func (s *Server) MessageSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var msg = new(entities.Message)
		if err := DecodeJSON(r.Body, msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uuid, err := uuid.NewUUID(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			msg.Id = uuid.String()
		}

		msg.Date = time.Now()
		//call the store save function
		err := s.store.MessageSave(ctx, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendMessage(w, &msg)
	}
}

// MessageUpdate updates a msg
//
// @ID MessageUpdate
// @Tags messages
// @Summary updates a message
// @Description updates a message
// @Param message body entities.Message true "Message"
// @Success 200 {object} entities.Message
// @Failure 400 string true "Invalid Argument"
// @Failure 500 string true "Internal Error"
// @Router /messages/{id} [put]
func (s *Server) MessageUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var msg = new(entities.Message)
		if err := DecodeJSON(r.Body, msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		//parse and validate
		uuid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg.Id = uuid.String()
		msg.Date = time.Now()

		//call the store save function
		err = s.store.MessageUpdate(ctx, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendMessage(w, &msg)
	}
}

// MessageGetByID get a Message by ID
//
// @ID MessageGetByID
// @Tags messages
// @Summary Get a Message by ID
// @Description Get a Message by ID
// @Param id path string true "ID"
// @Success 200 {object} entities.Message
// @Failure 400 string true "Invalid Argument"
// @Failure 500 string true "Internal Error"
// @Failure 404 string true "Not Found"
// @Router /messages/{id} [get]
func (s *Server) MessageGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		msg, err := s.store.MessageGetByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		sendMessage(w, &msg)
	}

}

// MessageDeleteByID deletes a msg by ID
//
// @ID MessageDeleteByID
// @Tags messages
// @Summary Delete a msg by ID
// @Description Delete a msg by ID
// @Param id path string true "ID"
// @Success 204 "Success"
// @Failure 400 string true "Invalid Argument"
// @Failure 500 string true "Internal Error"
// @Failure 404 string true "Not Found"
// @Router /messages/{id} [delete]
func (s *Server) MessageDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.store.MessageDeleteByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// MessageFind finds messages based on query-params
//
// @ID MessageFind
// @Tags messages
// @Summary Find messages
// @Description Find messages
// @Param id query string false "id"
// @Param offset query int false "offset"
// @Param limit query int false "limit"
// @Param sort query string false "query"
// @Success 200 {array} entities.Message
// @Failure 400 string true "Invalid Argument"
// @Failure 500 string true "Internal Error"
// @Router /messages [get]
func (s *Server) MessageFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		qp, err := queryp.ParseRawQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		messages, count, err := s.store.MessagesList(ctx, qp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sendMessage(w, &Results{Count: count, Results: messages})
	}
}

type Results struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
