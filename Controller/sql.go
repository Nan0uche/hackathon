package controller

import (
	"database/sql"
	"encoding/base64"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func GetLastDeliveryIDForUser(username string) string {
	db := InitDatabase()
	defer db.Close()

	// Exécuter la requête SQL pour récupérer l'ID du dernier enregistrement de la table Deliverys pour l'utilisateur donné
	query := `
        SELECT ID
        FROM Deliverys 
        WHERE Expéditeur = ? OR Destinataire = ?
        ORDER BY Livraison DESC
        LIMIT 1
    `
	var deliveryID string
	err := db.QueryRow(query, username, username).Scan(&deliveryID)
	if err != nil {
		log.Println("Error getting last delivery ID for user:", username, err)
		return ""
	}

	return deliveryID
}

func GetAllReceivePackageIDs(user string) []string {
	db := InitDatabase()
	defer db.Close()
	// Exécutez une requête SQL pour récupérer les ID des colis expédiés par l'utilisateur
	rows, err := db.Query("SELECT ID FROM Deliverys WHERE Destinataire = ?", user)
	if err != nil {
		return nil
	}
	defer rows.Close()
	// Parcourez les résultats et stockez les ID des colis dans une slice
	var packageIDs []string
	for rows.Next() {
		var packageID string
		if err := rows.Scan(&packageID); err != nil {
			return nil
		}
		packageIDs = append(packageIDs, packageID)
	}
	// Gestion des erreurs de parcours des lignes
	if err := rows.Err(); err != nil {
		return nil
	}
	return packageIDs
}
func GetAllShippedPackageIDs(user string) []string {
	db := InitDatabase()
	defer db.Close()
	// Exécutez une requête SQL pour récupérer les ID des colis expédiés par l'utilisateur
	rows, err := db.Query("SELECT ID FROM Deliverys WHERE Expéditeur = ?", user)
	if err != nil {
		return nil
	}
	defer rows.Close()
	// Parcourez les résultats et stockez les ID des colis dans une slice
	var packageIDs []string
	for rows.Next() {
		var packageID string
		if err := rows.Scan(&packageID); err != nil {
			return nil
		}
		packageIDs = append(packageIDs, packageID)
	}
	// Gestion des erreurs de parcours des lignes
	if err := rows.Err(); err != nil {
		return nil
	}
	return packageIDs
}
func IsAdminUser(username string) bool {
	// Établir une connexion à la base de données
	db := InitDatabase()
	defer db.Close()
	// Préparer la requête SQL pour récupérer le grade de l'utilisateur
	query := "SELECT Grade FROM Users WHERE Username = ?"
	var grade string
	err := db.QueryRow(query, username).Scan(&grade)
	if err != nil {
		log.Println("Erreur lors de la récupération du grade de l'utilisateur:", err)
		return false
	}
	// Vérifier si le grade de l'utilisateur est "Admin"
	return grade == "Admin"
}
func CheckIDExists(deliveryID string) bool {
	db := InitDatabase()
	defer db.Close()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Deliverys WHERE ID = ?", deliveryID).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'ID:", err)
		return false
	}
	return count > 0
}
func GetDeliveryWithID(deliveryID string) (expediteur, destinataire, livraison, status string, err error) {
	db := InitDatabase()
	defer db.Close()
	err = db.QueryRow("SELECT Expéditeur, Destinataire, Livraison, Status FROM Deliverys WHERE ID = ?", deliveryID).Scan(&expediteur, &destinataire, &livraison, &status)
	if err != nil {
		log.Println("Error getting delivery with ID:", deliveryID, err)
		return "", "", "", "", err
	}
	return expediteur, destinataire, livraison, status, nil
}

// GetColisWithDeliveryID récupère les informations d'un colis associé à une livraison avec l'ID de livraison donné.
func GetColisWithDeliveryID(deliveryID string) (titre, message string, imageBase64 string, err error) {
	db := InitDatabase()
	defer db.Close()
	var image []byte
	err = db.QueryRow("SELECT Titre, Message, Image FROM Colis WHERE ID = (SELECT Colis FROM Deliverys WHERE ID = ?)", deliveryID).Scan(&titre, &message, &image)
	if err != nil {
		log.Println("Error getting colis with delivery ID:", deliveryID, err)
		return "", "", "", err
	}
	imageBase64 = base64.StdEncoding.EncodeToString(image)
	return titre, message, imageBase64, nil
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func generateUniqueID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	for i := 4; i < len(b); i += 5 {
		b = append(b[:i], append([]byte{'-'}, b[i:]...)...)
	}
	return string(b)
}
func generateUniqueColisID(db *sql.DB) (string, error) {
	for {
		id := generateUniqueID(12)
		// Vérifier si l'identifiant est déjà présent dans la base de données
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM Colis WHERE ID = ?", id).Scan(&count)
		if err != nil {
			return "", err
		}
		if count == 0 {
			// L'identifiant est unique, donc le retourner
			return id, nil
		}
	}
}
func InitDatabase() *sql.DB {
	dbPath := "./Model/database.db"
	// Vérifier si le fichier de base de données existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Créer le répertoire /Model s'il n'existe pas
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatal(err)
		}
	}
	db, err := sql.Open("sqlite3", "./Model/database.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func InitTables(db *sql.DB) {
	sqlUser := `CREATE TABLE IF NOT EXISTS Users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT,
			Email TEXT,
			Password TEXT,
			Grade TEXT,
			Location TEXT)`
	_, err := db.Exec(sqlUser)
	if err != nil {
		log.Fatal(err)
	}
	sqlDelivery := `CREATE TABLE IF NOT EXISTS Deliverys (
			ID TEXT PRIMARY KEY,
			Expéditeur TEXT,
			Destinataire TEXT,
			Colis TEXT,
			Livraison TEXT,
			Status TEXT)`
	_, err = db.Exec(sqlDelivery)
	if err != nil {
		log.Fatal(err)
	}
	sqlColis := `CREATE TABLE IF NOT EXISTS Colis (
			ID TEXT PRIMARY KEY,
			Titre TEXT,
			Message TEXT,
			Image TEXT)`
	_, err = db.Exec(sqlColis)
	if err != nil {
		log.Fatal(err)
	}
}
func InsertUser(username, email, password, location string) error {
	var grade string
	err := InitDatabase().QueryRow("SELECT Grade FROM Users WHERE Grade != '' LIMIT 1").Scan(&grade)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	// Si la colonne Grade est vide, insérer l'utilisateur avec le grade "Member"
	if grade == "" {
		_, err := InitDatabase().Exec("INSERT INTO Users (Username, Email, Password, Grade, Location) VALUES (?, ?, ?, 'Member', ?)", username, email, password, location)
		if err != nil {
			return err
		}
	} else {
		// Sinon, ne pas toucher à la colonne Grade et insérer l'utilisateur avec la valeur fournie
		_, err := InitDatabase().Exec("INSERT INTO Users (Username, Email, Password, Grade, Location) VALUES (?, ?, ?, ?, ?)", username, email, password, grade, location)
		if err != nil {
			return err
		}
	}
	return nil
}
func InsertDelivery(expediteur, destinataire, colis, livraison string) error {
	// Insérer la livraison dans la base de données
	deliveryID, err := generateUniqueColisID(InitDatabase())
	if err != nil {
		log.Println("Erreur lors de la génération de l'identifiant du colis:", err)
		return err
	}
	_, err = InitDatabase().Exec("INSERT INTO Deliverys (ID, Expéditeur, Destinataire, Colis, Livraison, Status) VALUES (?, ?, ?, ?, ?, ?)", deliveryID, expediteur, destinataire, colis, livraison, "Enregistré")
	if err != nil {
		log.Println("Erreur lors de l'insertion de la livraison:", err)
		return err
	}
	return err
}
func InsertColis(titre, message string, image []byte) (string, error) {
	// Générer un identifiant unique pour le colis
	colisID, err := generateUniqueColisID(InitDatabase())
	if err != nil {
		log.Println("Erreur lors de la génération de l'identifiant du colis:", err)
		return "", err
	}
	// Insérer le colis dans la base de données avec l'identifiant unique
	_, err = InitDatabase().Exec("INSERT INTO Colis (ID, Titre, Message, Image) VALUES (?, ?, ?, ?)", colisID, titre, message, image)
	if err != nil {
		log.Println("Erreur lors de l'insertion du colis:", err)
		return "", err
	}
	return colisID, nil
}

// UserExists vérifie si un utilisateur avec le nom d'utilisateur donné existe déjà dans la base de données.
func UserExists(username string) bool {
	// Établir une connexion à la base de données
	db := InitDatabase()
	defer db.Close()
	// Préparer la requête SQL pour vérifier l'existence de l'utilisateur
	query := "SELECT COUNT(*) FROM Users WHERE Username = ?"
	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'utilisateur:", err)
		return false
	}
	// Si le nombre d'utilisateurs trouvés est supérieur à zéro, alors l'utilisateur existe
	return count > 0
}
func MailExists(email string) bool {
	// Établir une connexion à la base de données
	db := InitDatabase()
	defer db.Close()
	// Préparer la requête SQL pour vérifier l'existence de l'utilisateur
	query := "SELECT COUNT(*) FROM Users WHERE Email = ?"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'email:", err)
		return false
	}
	// Si le nombre d'utilisateurs trouvés est supérieur à zéro, alors l'utilisateur existe
	return count > 0
}
func GetPasswordByEmail(email string) (string, error) {
	// Ouverture de la connexion à la base de données
	db := InitDatabase()
	defer db.Close()
	// Préparez la requête SQL pour sélectionner le mot de passe haché en fonction de l'e-mail
	query := "SELECT password FROM users WHERE email = ?"
	row := db.QueryRow(query, email)
	// Variables pour stocker les résultats de la requête
	var hashedPassword string
	// Exécuter la requête et scanner le résultat dans les variables
	err := row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with email:", email)
			return "", err
		}
		log.Println("Error getting password for email:", email, err)
		return "", err
	}
	return hashedPassword, nil
}
func GetUsernameByEmail(email string) (string, error) {
	db := InitDatabase()
	defer db.Close()
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with email:", email)
			return "", err
		}
		log.Println("Error getting username for email:", email, err)
		return "", err
	}
	return username, nil
}
func GetRoleByEmail(email string) string {
	db := InitDatabase()
	defer db.Close()
	var grade string
	err := db.QueryRow("SELECT grade FROM users WHERE email = ?", email).Scan(&grade)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with email:", email)
			return ""
		}
		log.Println("Error getting grade for email:", email, err)
		return ""
	}
	return grade
}
func GetLocationByEmail(email string) (string, error) {
	db := InitDatabase()
	defer db.Close()
	var location string
	err := db.QueryRow("SELECT location FROM users WHERE email = ?", email).Scan(&location)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with email:", email)
			return "", err
		}
		log.Println("Error getting location for email:", email, err)
		return "", err
	}
	return location, nil
}
func UpdateUsername(email, ancientName, newUsername string) error {
	db := InitDatabase()
	defer db.Close()
	_, err := db.Exec("UPDATE users SET username = ? WHERE email = ?", newUsername, email)
	if err != nil {
		log.Println("Error updating username for email:", email, err)
		return err
	}
	_, err = db.Exec("UPDATE Deliverys SET Expéditeur = ? WHERE Expéditeur =?", newUsername, ancientName)
	if err != nil {
		log.Println("Error updating Deliverys for username:", ancientName, err)
		return err
	}
	_, err = db.Exec("UPDATE Deliverys SET Destinataire = ? WHERE Destinataire =?", newUsername, ancientName)
	if err != nil {
		log.Println("Error updating Deliverys for username:", ancientName, err)
		return err
	}
	return nil
}
func UpdateLocation(email, newLocation string) error {
	db := InitDatabase()
	defer db.Close()
	_, err := db.Exec("UPDATE users SET location = ? WHERE email = ?", newLocation, email)
	if err != nil {
		log.Println("Error updating location for email:", email, err)
		return err
	}
	return nil
}
func UpdatePassword(email, hashedPassword string) error {
	db := InitDatabase()
	defer db.Close()
	_, err := db.Exec("UPDATE users SET password = ? WHERE email = ?", hashedPassword, email)
	if err != nil {
		log.Println("Error updating password for email:", email, err)
		return err
	}
	return nil
}
func UpdateEmail(email, newMail string) error {
	db := InitDatabase()
	defer db.Close()
	_, err := db.Exec("UPDATE users SET email = ? WHERE email = ?", newMail, email)
	if err != nil {
		log.Println("Error updating password for email:", email, err)
		return err
	}
	return nil
}
func DeleteUser(email string) error {
	db := InitDatabase()
	defer db.Close()
	_, err := db.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		log.Println("Error deleting user:", email, err)
		return err
	}
	return nil
}
func Getallusername() []string {
	db := InitDatabase()
	defer db.Close()
	rows, err := db.Query("SELECT Username FROM Users")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil
		}
		usernames = append(usernames, username)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return usernames
}
func GetUserLocation(username string) string {
	db := InitDatabase()
	defer db.Close()
	var location string
	err := db.QueryRow("SELECT location FROM users WHERE username = ?", username).Scan(&location)
	if err != nil {
		if err == sql.ErrNoRows {

			return ""
		}
		log.Println("Error getting location for username:", username, err)
		return ""
	}
	return location
}
