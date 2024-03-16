package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Status struct merepresentasikan struktur dari JSON status.
type Status struct {
	Status struct {
		Water int `json:"water"` // Nilai air
		Wind  int `json:"wind"`  // Nilai angin
	} `json:"status"`
}

func main() {
	// Set up HTTP server untuk menyajikan halaman status dan endpoint JSON
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/status", statusJSONHandler)

	go updateStatusJSON() // Goroutine untuk memperbarui status JSON setiap 15 detik

	// Start HTTP server
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// statusHandler menangani permintaan HTTP untuk halaman status.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Memuat status dari file JSON
	status, err := loadStatusFromJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Menentukan status air dan angin
	waterStatus := determineWaterStatus(status.Status.Water)
	windStatus := determineWindStatus(status.Status.Wind)

	// Merender halaman HTML
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Status Page</title>
			<meta http-equiv="refresh" content="5"> <!-- Auto-reload setiap 5 detik -->
		</head>
		<body>
			<h1>Status Page</h1>
			<p>Water Status: %s</p>
			<p>Wind Status: %s</p>
			<p>Terakhir Update: %s</p>
		</body>
		</html>
	`, waterStatus, windStatus, time.Now().Format("2006-01-02 15:04:05"))

	// Menulis respons HTML
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, html)
}

// statusJSONHandler menangani permintaan HTTP untuk JSON status.
func statusJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Memuat status dari file JSON
	status, err := loadStatusFromJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal status ke JSON
	jsonData, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Menulis respons JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// updateStatusJSON secara berkala memperbarui file JSON status setiap 15 detik.
func updateStatusJSON() {
	for {
		// Menghasilkan nilai acak untuk air dan angin
		water := rand.Intn(100) + 1
		wind := rand.Intn(100) + 1

		// Membuat struktur status
		status := Status{}
		status.Status.Water = water
		status.Status.Wind = wind

		// Marshal status ke JSON
		jsonData, err := json.Marshal(status)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Menulis data JSON ke file
		err = writeJSONToFile(jsonData)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Tidur selama 15 detik
		time.Sleep(15 * time.Second)
	}
}

// loadStatusFromJSON memuat status dari file JSON.
func loadStatusFromJSON() (Status, error) {
	// Membuka file JSON
	file, err := os.Open("status.json")
	if err != nil {
		return Status{}, err
	}
	defer file.Close()

	// Mendekode file JSON menjadi struktur status
	var status Status
	err = json.NewDecoder(file).Decode(&status)
	if err != nil {
		return Status{}, err
	}

	return status, nil
}

// writeJSONToFile menulis data JSON ke file JSON status.
func writeJSONToFile(data []byte) error {
	// Membuka atau membuat file JSON
	file, err := os.OpenFile("status.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Menulis data JSON ke file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// determineWaterStatus menentukan status air berdasarkan nilainya.
func determineWaterStatus(water int) string {
	switch {
	case water < 5:
		return "Aman"
	case water >= 6 && water <= 8:
		return "Siaga"
	default:
		return "Bahaya"
	}
}

// determineWindStatus menentukan status angin berdasarkan nilainya.
func determineWindStatus(wind int) string {
	switch {
	case wind < 6:
		return "Aman"
	case wind >= 7 && wind <= 15:
		return "Siaga"
	default:
		return "Bahaya"
	}
}
