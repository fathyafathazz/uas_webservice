package album

import (
	"encoding/json"
	"net/http"
	"strconv"

	"uas_musik/database"
	"github.com/gorilla/mux"
	"uas_musik/model/albums"
)

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
        if err := rows.Scan(&c.AlbumId,&c.Title,&c.Artist_Id, &c.Price, &c.Year); err != nil {
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

func PostAlbum(w http.ResponseWriter, r *http.Request) {
	var pc album.Album
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new course
	query := `
		INSERT INTO albums (title, artist_id, price, year) 
		VALUES (?, ?, ?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.Title, pc.Artist_Id, pc.Price, pc.Year)
	if err != nil {
		http.Error(w, "Failed to insert artist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Course added successfully",
		"id":      id,
	})
}

func PutAlbum(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
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

	// Decode JSON body
	var pc album.Album
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the category admin
	query := `
		UPDATE albums 
		SET title=?, artist_id=?, price=?, year=?
		WHERE id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.Title, pc.Artist_Id, pc.Price, pc.Year, id)
	if err != nil {
		http.Error(w, "Failed to update albums: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Course updated successfully",
	})
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
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

	// Prepare the SQL statement for deleting a category admin
	query := `
		DELETE FROM albums
		WHERE id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete album: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Album deleted successfully",
	})
}