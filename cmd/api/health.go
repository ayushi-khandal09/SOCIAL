package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Health check
//	@Description	Checks the health of the server
//	@Tags			ops
//	@Produce		json
//	@Success		200		{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
