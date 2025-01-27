package handler

import (
	"encoding/json"
	"net/http"

	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetNotes ...
// swagger:route GET /api/profile/{id}/notes Notes getNotes
//
// # Retrieves the list of notes for the selected user
//
// Parameters:
//   - +name: id
//     in: path
//     description: The userID of the selected user
//     type: string
//     required: true
//
// Responses:
// 200: []Note
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func GetNotes(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || len(id) <= 0 {
		config.Log("Unable to retrieve details without an id. ", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		config.Log("Unable to retrieve details with provided id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	resp, err := db.RetrieveNotes(user, parsedUUID)
	if err != nil {
		config.Log("Unable to retrieve notes", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// AddNewNote ...
// swagger:route POST /api/profile/{id}/notes Notes addNewNote
//
// # Add a new note to the database
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected note
//     type: string
//     required: true
//   - +name: Note
//     in: body
//     description: The note object to add into the db
//     type: Note
//     required: true
//
// Responses:
// 200: Note
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func AddNewNote(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to update notes with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		config.Log("Error decoding data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := db.AddNewNote(user, userID, note)
	if err != nil {
		config.Log("Unable to add new note", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// UpdateNote ...
// swagger:route PUT /api/profile/{id}/notes Notes updateNote
//
// # Updates an existing note in the database
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected note
//     type: string
//     required: true
//   - +name: Note
//     in: body
//     description: The note object to update into the db
//     type: Note
//     required: true
//
// Responses:
// 200: Note
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func UpdateNote(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if len(userID) <= 0 {
		config.Log("Unable to update notes with empty id", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		config.Log("Error decoding data", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := db.UpdateNote(user, userID, note)
	if err != nil {
		config.Log("Unable to update notes", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// RemoveNote ...
// swagger:route DELETE /api/profile/{id}/notes Notes removeNote
//
// # Removes the note from the db
//
// Parameters:
//   - +name: id
//     in: path
//     description: The id of the selected note
//     type: string
//     required: true
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func RemoveNote(rw http.ResponseWriter, r *http.Request, user string) {
	vars := mux.Vars(r)
	userID := vars["id"]
	noteID := vars["noteID"]

	if len(userID) <= 0 {
		config.Log("Unable to update notes with empty userID", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	if len(noteID) <= 0 {
		config.Log("Unable to update notes with empty noteID", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(nil)
		return
	}

	var note model.Note
	note.ID = noteID

	err := db.RemoveNote(user, note.ID)
	if err != nil {
		config.Log("Unable to remove notes", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(note.ID)
}
