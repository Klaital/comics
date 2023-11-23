package comicserver

import (
	"database/sql"
	"errors"
	"net/http"
)

func (srv *Server) GetAllComics(w http.ResponseWriter, r *http.Request) {
	// TODO: validate logged-in user instead of hardcoded value
	var userId uint64 = 1
	allComics, err := srv.storage.GetAllComics(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		srv.l.WithError(err).WithField("UserID", userId).Debug("Error fetching comics list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: convert
}
