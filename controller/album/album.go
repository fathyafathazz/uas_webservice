package album

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"uas_musik/database"
	"uas_musik/model/albums"

	"github.com/gorilla/mux"
)

// GetAlbum handles GET requests to fetch all albums
func GetAlbum(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM albums")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var albums []album.Album
	for rows.Next() {
		var c album.Album
		if err := rows.Scan(&c.AlbumId, &c.Title, &c.Artist_Id, &c.Price, &c.Year); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		albums = append(albums, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

// GetAlbumByID handles GET requests to fetch a single album by its ID
func GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var album album.Album
	query := "SELECT * FROM albums WHERE id = ?"
	err = database.DB.QueryRow(query, id).Scan(&album.AlbumId, &album.Title, &album.Artist_Id, &album.Price, &album.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Album not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

// PostAlbum handles POST requests to create a new album
func PostAlbum(w http.ResponseWriter, r *http.Request) {
	var pc album.Album
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO albums (title, artist_id, price, year) VALUES (?, ?, ?, ?)`
	res, err := database.DB.Exec(query, pc.Title, pc.Artist_Id, pc.Price, pc.Year)
	if err != nil {
		http.Error(w, "Failed to insert album: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Album added successfully",
		"id":      id,
	})
}

// PutAlbum handles PUT requests to update an existing album
func PutAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var pc album.Album
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `UPDATE albums SET title=?, artist_id=?, price=?, year=? WHERE id=?`
	result, err := database.DB.Exec(query, pc.Title, pc.Artist_Id, pc.Price, pc.Year, id)
	if err != nil {
		http.Error(w, "Failed to update album: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Album updated successfully",
	})
}

// DeleteAlbum handles DELETE requests to delete an album
func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM albums WHERE id = ?`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete album: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Album deleted successfully",
	})
}