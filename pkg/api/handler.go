package api

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/DanilLagunov/image-compressor/pkg/storage"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	storage storage.Storage
	queue   chan UploadRequest
}

type UploadRequest struct {
	img      multipart.File
	header   *multipart.FileHeader
	response chan error
}

func New(ctx context.Context, s storage.Storage) *Handler {
	queue := make(chan UploadRequest, 5)
	h := Handler{
		storage: s,
		queue:   queue,
	}
	go h.QueueController(ctx)
	h.Router = h.initRoutes()
	return &h
}

func (h Handler) initRoutes() *mux.Router {
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/upload", h.uploadImage).Methods(http.MethodPost)
	h.Router.HandleFunc("/images/{id}", h.getImage).Methods(http.MethodGet)
	return h.Router
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	h.Router.ServeHTTP(w, req)

	fmt.Printf("request time is %v \n", time.Since(start))
}

func (h Handler) QueueController(ctx context.Context) {
	for {
		select {
		case imgReq := <-h.queue:
			fmt.Println("here")
			err := h.storage.SaveImage(imgReq.img, imgReq.header)
			imgReq.response <- err
		case <-ctx.Done():

		}
	}
}
