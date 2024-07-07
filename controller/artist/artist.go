package artist

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"log"
	"fmt"

	"uas_musik/database"
	"github.com/gorilla/mux"
	"uas_musik/model/artist"
)

func GetArtist(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM artists")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var artists []artist.Artist
    for rows.Next() {
        var c artist.Artist
        if err := rows.Scan(&c.ArtistId,&c.Name,&c.Genre); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        artists = append(artists, c)
    }
	
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// GetArtistByID handles GET requests to fetch a single album by its ID
func GetArtistByID(w http.ResponseWriter, r *http.Request) {
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

	var artist artist.Artist
	query := "SELECT * FROM artists WHERE id = ?"
	err = database.DB.QueryRow(query, id).Scan(&artist.ArtistId, &artist.Name, &artist.Genre)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Album not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

func PostArtist(w http.ResponseWriter, r *http.Request) {
	var pc artist.Artist
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new course
	query := `
		INSERT INTO artists (name, genre) 
		VALUES (?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.Name, pc.Genre)
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
		"message": "Artist added successfully",
		"id":      id,
	})
}
func PutArtist(w http.ResponseWriter, r *http.Request) {
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
	var pc artist.Artist
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Tutup body request setelah digunakan

	// Check if artist exists
	if !artistExists(id) {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Prepare the SQL statement for updating the category admin
	query := `
		UPDATE artists 
		SET name=?, genre=?
		WHERE id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.Name, pc.Genre, id)
	if err != nil {
		log.Println("Failed to update artist:", err)
		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to retrieve affected rows:", err)
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
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
		"message": fmt.Sprintf("Artist with ID %d updated successfully", id),
	})
}
func artistExists(id int) bool {
	// Persiapkan kueri SQL untuk memeriksa keberadaan artis berdasarkan ID
	query := "SELECT COUNT(*) FROM artists WHERE id = ?"

	// Eksekusi kueri dan periksa hasilnya
	var count int
	err := database.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		// Penanganan kesalahan jika terjadi kesalahan saat menjalankan kueri
		log.Println("Error checking artist existence:", err)
		return false // Asumsikan artis tidak ada jika terjadi kesalahan
	}

	return count > 0 // Kembalikan true jika jumlah artis dengan ID yang diberikan lebih dari 0
}

func DeleteArtist(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM artists
		WHERE id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete artist: "+err.Error(), http.StatusInternalServerError)
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
		"message": "Artist deleted successfully",
	})
}