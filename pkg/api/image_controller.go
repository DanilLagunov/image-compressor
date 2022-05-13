package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) getImage(w http.ResponseWriter, r *http.Request) {
	quality := r.URL.Query().Get("q")

	vars := mux.Vars(r)

	id := vars["id"]

	img, err := h.storage.GetImage(id, quality)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(img)

	return
}

func (h *Handler) uploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the photo key - First Approach
	file, header, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := make(chan error, 1)
	req := UploadRequest{
		img:      file,
		header:   header,
		response: response,
	}
	h.queue <- req
	fmt.Println("1")

	if err := <-req.response; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("1")
	w.WriteHeader(http.StatusOK)
	return
}
