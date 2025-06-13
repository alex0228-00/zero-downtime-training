package src

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	port string
	mngr AssetManager
}

func NewServer(port string, mngr AssetManager) *Server {
	return &Server{
		port: port,
		mngr: mngr,
	}
}

func (srv *Server) Start() error {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/asset", func(w http.ResponseWriter, r *http.Request) {
		var handler http.HandlerFunc
		switch r.Method {
		case http.MethodGet:
			handler = srv.handleGetAsset()
		case http.MethodPost:
			handler = srv.handleCreateAsset()
		case http.MethodPut:
			handler = srv.handleUpdateAsset()
		case http.MethodDelete:
			handler = srv.handleDeleteAsset()
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		handler(w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%s", srv.port), nil)
}

func (srv *Server) handleGetAsset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		asset, err := srv.mngr.ReadAsset(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading asset: %v", err), http.StatusInternalServerError)
			return
		}

		if asset == nil {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (srv *Server) handleCreateAsset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var asset Asset
		if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding asset: %v", err), http.StatusBadRequest)
			return
		}

		if err := srv.mngr.CreateAsset(&asset); err != nil {
			http.Error(w, fmt.Sprintf("Error creating asset: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (srv *Server) handleUpdateAsset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var asset Asset
		if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding asset: %v", err), http.StatusBadRequest)
			return
		}

		if err := srv.mngr.UpdateAssetSourceByID(asset.ID, asset.Source); err != nil {
			http.Error(w, fmt.Sprintf("Error updating asset: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (srv *Server) handleDeleteAsset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		if err := srv.mngr.DeleteAsset(id); err != nil {
			http.Error(w, fmt.Sprintf("Error deleting asset: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
