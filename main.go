package main

import (
	controller "HACKATHON/Controller"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type HomeData struct {
	IsLoggedIn   bool
	Username     string
	ErrorMessage string
	IsAdmin      bool
}

var tmplHome = template.Must(template.ParseFiles("./View/html/home.html"))
var tmpl404 = template.Must(template.ParseFiles("./View/html/404.html"))
var tmplLogin = template.Must(template.ParseFiles("./View/html/login.html"))
var tmplRegister = template.Must(template.ParseFiles("./View/html/register.html"))
var tmplProfil = template.Must(template.ParseFiles("./View/html/profil.html"))
var tmplTracking = template.Must(template.ParseFiles("./View/html/tracking.html"))
var tmplExped = template.Must(template.ParseFiles("./View/html/expedition.html"))
var tmplSuivi = template.Must(template.ParseFiles("./View/html/suivi.html"))
var tmplFaq = template.Must(template.ParseFiles("./View/html/faq.html"))
var tmplContactus = template.Must(template.ParseFiles("./View/html/contactus.html"))
var tmplDashboard = template.Must(template.ParseFiles("./View/html/dashboard.html"))
var tmplSuccesfuly = template.Must(template.ParseFiles("./View/html/succesfuly.html"))

func main() {
	db := controller.InitDatabase()
	defer db.Close()
	controller.InitTables(db)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./View/html"))))
	fmt.Printf("\n")
	fmt.Println("http://localhost:8080/")
	fmt.Printf("\n")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/404", pageNotFoundHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/update_profile", updateProfileHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/tracking", trackingHandler)
	http.HandleFunc("/expedition", expeditionHandler)
	http.HandleFunc("/envoyer_colis", sendColisHandler)
	http.HandleFunc("/faq", faqHandler)
	http.HandleFunc("/contactus", contactusHandler)
	http.HandleFunc("/suivi", suiviHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/succesfuly", sendColisHandler)
	http.HandleFunc("/get-location", getLocationHandler)
	http.ListenAndServe("localhost:8080", nil)
}
func getLocationHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le nom du destinataire de la requête GET
	destinataire := r.URL.Query().Get("destinataire")
	if destinataire == "" {
		http.Error(w, "No destinataire provided", http.StatusBadRequest)
		return
	}
	// Récupérer la localisation du destinataire
	location := controller.GetUserLocation(destinataire)
	if location == "" {
		http.Error(w, "Failed to retrieve location for destinataire", http.StatusInternalServerError)
		return
	}
	// Créer une structure pour la réponse JSON
	response := struct {
		Success  bool   `json:"success"`
		Location string `json:"location"`
	}{
		Success:  true,
		Location: location,
	}
	// Encoder la réponse en JSON et l'envoyer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		log.Println("Failed to encode JSON response:", err)
		return
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var Username string
	var isAdmin bool
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	isLoggedIn := false
	cookie, err := r.Cookie("username")
	if err == nil {
		isLoggedIn = true
		email := cookie.Value
		// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
		Username, err = controller.GetUsernameByEmail(email)
		if err != nil {
			// Gérer l'erreur, afficher un message d'erreur, selon votre logique
			http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
			return
		}
		isAdmin = controller.IsAdminUser(Username)
	}
	// Déterminer si l'utilisateur est connecté
	// Rediriger vers la page 404 si l'URL n'est pas l'une des pages autorisées
	if r.URL.Path != "/" && r.URL.Path != "/login" && r.URL.Path != "/register" && r.URL.Path != "/profil" && r.URL.Path != "/update_profile" && r.URL.Path != "/expeditionFr" && r.URL.Path != "/expeditionInt" && r.URL.Path != "/colis_envoye" && r.URL.Path != "/suivi" {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	// Passer les données au modèle et exécuter le template
	data := HomeData{
		IsLoggedIn: isLoggedIn,
		Username:   Username,
		IsAdmin:    isAdmin,
	}
	err = tmplHome.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	email := cookie.Value
	role := controller.GetRoleByEmail(email)
	if role != "Admin" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	Username, err := controller.GetUsernameByEmail(email)
	if err != nil {
		// Gérer l'erreur, afficher un message d'erreur, selon votre logique
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}
	data := HomeData{
		Username: Username,
	}
	err = tmplDashboard.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl404.Execute(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func faqHandler(w http.ResponseWriter, r *http.Request) {
	var Username string
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	isLoggedIn := false
	var isAdmin bool
	cookie, err := r.Cookie("username")
	if err == nil {
		isLoggedIn = true
		email := cookie.Value
		// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
		Username, err = controller.GetUsernameByEmail(email)
		if err != nil {
			// Gérer l'erreur, afficher un message d'erreur, selon votre logique
			http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
			return
		}
		isAdmin = controller.IsAdminUser(Username)
	}
	data := HomeData{
		IsLoggedIn: isLoggedIn,
		Username:   Username,
		IsAdmin:    isAdmin,
	}
	err = tmplFaq.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func suiviHandler(w http.ResponseWriter, r *http.Request) {
	deliveryID := getTrackingID(r.URL.RawQuery)
	if deliveryID != "" {
		if controller.CheckIDExists(deliveryID) {
			expediteur, destinataire, livraison, status, err := controller.GetDeliveryWithID(deliveryID)
			if err != nil {
				log.Fatal("Erreur lors de la récupération des informations de livraison:", err)
			}
			titre, message, image, err := controller.GetColisWithDeliveryID(deliveryID)
			if err != nil {
				log.Fatal("Erreur lors de la récupération des informations du colis associé à la livraison:", err)
			}
			data := struct {
				ID           string
				Expediteur   string
				Destinataire string
				Livraison    string
				Status       string
				Titre        string
				Message      string
				Image        string
			}{
				ID:           deliveryID,
				Expediteur:   expediteur,
				Destinataire: destinataire,
				Livraison:    livraison,
				Status:       status,
				Titre:        titre,
				Message:      message,
				Image:        image,
			}
			err = tmplSuivi.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
func contactusHandler(w http.ResponseWriter, r *http.Request) {
	var Username string
	var isAdmin bool
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	isLoggedIn := false
	cookie, err := r.Cookie("username")
	if err == nil {
		isLoggedIn = true
		email := cookie.Value
		// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
		Username, err = controller.GetUsernameByEmail(email)
		if err != nil {
			// Gérer l'erreur, afficher un message d'erreur, selon votre logique
			http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
			return
		}
		isAdmin = controller.IsAdminUser(Username)
	}
	data := HomeData{
		IsLoggedIn: isLoggedIn,
		Username:   Username,
		IsAdmin:    isAdmin,
	}
	err = tmplContactus.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := tmplRegister.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// Récupérer les valeurs du formulaire
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	// Hacher le mot de passe avec bcrypt
	hashedPassword := HashPassword(password)
	location := r.FormValue("location")
	if controller.UserExists(username) {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Username already taken",
		}
		err := tmplRegister.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if controller.MailExists(email) {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Email already taken",
		}
		err := tmplRegister.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err := controller.InsertUser(username, email, hashedPassword, location)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("Failed to insert user into database:", err)
		return
	}

	expiration := time.Now().Add(2 * time.Hour) // Durée de validité du cookie (1 jour dans cet exemple)
	cookie := http.Cookie{Name: "username", Value: email, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := tmplLogin.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Vérifier si l'e-mail existe dans la base de données
	if !controller.MailExists(email) {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Email non enregistré",
		}
		err := tmplLogin.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// Récupérer le mot de passe haché associé à l'e-mail
	hashedPassword, err := controller.GetPasswordByEmail(email)
	if err != nil {
		http.Error(w, "Failed to retrieve password", http.StatusInternalServerError)
		log.Println("Failed to retrieve password:", err)
		return
	}
	// Vérifier si le mot de passe saisi correspond au mot de passe haché
	if !CheckPasswordHash(password, hashedPassword) {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Mot de passe incorrect",
		}
		err := tmplLogin.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	expiration := time.Now().Add(2 * time.Hour) // Durée de validité du cookie (1 jour dans cet exemple)
	cookie := http.Cookie{Name: "username", Value: email, Expires: expiration}
	http.SetCookie(w, &cookie)
	// Connexion réussie, vous pouvez maintenant rediriger l'utilisateur vers une page sécurisée
	role := controller.GetRoleByEmail(email)
	if role == "Admin" {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
func profilHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	cookie, err := r.Cookie("username")
	if err != nil {
		// Gérer l'erreur, rediriger ou afficher un message d'erreur, selon votre logique
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	email := cookie.Value
	// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
	username, err := controller.GetUsernameByEmail(email)
	if err != nil {
		// Gérer l'erreur, afficher un message d'erreur, selon votre logique
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}
	location, err := controller.GetLocationByEmail(email)
	if err != nil {
		// Gérer l'erreur, afficher un message d'erreur, selon votre logique
		http.Error(w, "Failed to retrieve location", http.StatusInternalServerError)
		return
	}
	// Créer un modèle avec les valeurs du pseudo et de l'email
	data := struct {
		Username string
		Email    string
		Location string
	}{
		Username: username,
		Email:    email,
		Location: location,
	}
	// Rendre la page HTML du profil en passant le modèle contenant les valeurs
	err = tmplProfil.Execute(w, data)
	if err != nil {
		// Gérer l'erreur, afficher un message d'erreur, selon votre logique
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si le bouton de suppression a été cliqué
	if r.FormValue("deleteAccount") == "true" {
		// Récupérer l'email de l'utilisateur à partir du cookie
		cookie, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		email := cookie.Value
		// Supprimer l'utilisateur de la base de données
		err = controller.DeleteUser(email)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			log.Println("Failed to delete user:", err)
			return
		}
		// Effacer le cookie de l'utilisateur
		expiration := time.Now().AddDate(0, 0, -1) // Réglage de la date d'expiration du cookie dans le passé pour le supprimer
		cookie = &http.Cookie{Name: "username", Value: "", Expires: expiration}
		http.SetCookie(w, cookie)
		// Rediriger l'utilisateur vers la page d'accueil ou une autre page appropriée
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Vérifier si l'utilisateur est connecté en vérifiant le cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		// Cookie non trouvé, rediriger vers la page de connexion
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Récupérer les nouvelles informations fournies par l'utilisateur
	newEmail := r.FormValue("newEmail")
	newUsername := r.FormValue("newUsername")
	newLocation := r.FormValue("newLocation")
	newPassword := r.FormValue("newPassword")
	// Récupérer l'e-mail de l'utilisateur à partir du cookie
	email := cookie.Value
	currentUsername, err := controller.GetUsernameByEmail(email)
	if err != nil {
		http.Error(w, "Failed to retrieve current username", http.StatusInternalServerError)
		log.Println("Failed to retrieve current username:", err)
		return
	}
	if controller.MailExists(newEmail) && email != newEmail {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}
	if controller.UserExists(newUsername) && currentUsername != newUsername {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}
	if currentUsername != newUsername {
		// Mettre à jour le pseudo de l'utilisateur dans la base de données
		err = controller.UpdateUsername(email, currentUsername, newUsername)
		if err != nil {
			http.Error(w, "Failed to update username", http.StatusInternalServerError)
			log.Println("Failed to update username:", err)
			return
		}
	}
	// Mettre à jour le mot de passe de l'utilisateur s'il a fourni un nouveau mot de passe
	if newPassword != "" {
		hashedPassword := HashPassword(newPassword)
		err = controller.UpdatePassword(email, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			log.Println("Failed to update password:", err)
			return
		}
	}
	// Mettre à jour le mail de l'utilisateur dans la base de données
	if email != newEmail {
		err = controller.UpdateEmail(email, newEmail)
		if err != nil {
			http.Error(w, "Failed to update username", http.StatusInternalServerError)
			log.Println("Failed to update username:", err)
			return
		}
	}
	err = controller.UpdateLocation(email, newLocation)
	if err != nil {
		http.Error(w, "Failed to update location", http.StatusInternalServerError)
		log.Println("Failed to update location:", err)
		return
	}
	expiration := time.Now().Add(2 * time.Hour)
	cookie.Value = newEmail
	cookie.Expires = expiration
	http.SetCookie(w, cookie)
	// Rediriger l'utilisateur vers la page de profil après la mise à jour réussie
	http.Redirect(w, r, "/profil", http.StatusFound)
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{Name: "username", Value: "", Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
func trackingHandler(w http.ResponseWriter, r *http.Request) {
	var Username string
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	isLoggedIn := false
	var isAdmin bool
	cookie, err := r.Cookie("username")
	if err == nil {
		isLoggedIn = true
		email := cookie.Value
		// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
		Username, err = controller.GetUsernameByEmail(email)
		if err != nil {
			// Gérer l'erreur, afficher un message d'erreur, selon votre logique
			http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
			return
		}
		isAdmin = controller.IsAdminUser(Username)
	} else {
		// Gérer l'erreur, rediriger ou afficher un message d'erreur, selon votre logique
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	IDExped := controller.GetAllShippedPackageIDs(Username)
	IDRecu := controller.GetAllReceivePackageIDs(Username)
	// Passer les données au modèle et exécuter le template
	data := struct {
		IsLoggedIn bool
		Username   string
		IsAdmin    bool
		IDExped    []string
		IDRecu     []string
	}{
		IsLoggedIn: isLoggedIn,
		Username:   Username,
		IsAdmin:    isAdmin,
		IDExped:    IDExped,
		IDRecu:     IDRecu,
	}
	err = tmplTracking.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func expeditionHandler(w http.ResponseWriter, r *http.Request) {
	var Username string
	isLoggedIn := false
	var allusers []string
	var isAdmin bool
	errorMessage := r.URL.Query().Get("error")
	// Récupérer les informations de l'utilisateur depuis le cookie ou toute autre source
	cookie, err := r.Cookie("username")
	if err == nil {
		email := cookie.Value
		isLoggedIn = true
		allusers = controller.Getallusername()
		// Récupérer le pseudo associé à l'email de l'utilisateur depuis la base de données
		Username, err = controller.GetUsernameByEmail(email)
		if err != nil {
			// Gérer l'erreur, afficher un message d'erreur, selon votre logique
			http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
			log.Println("Failed to retrieve username:", err)
			return
		}
		isAdmin = controller.IsAdminUser(Username)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Passer les données au modèle et exécuter le template
	data := struct {
		Allusers     []string
		Username     string
		IsLoggedIn   bool
		IsAdmin      bool
		ErrorMessage string
	}{
		Allusers:     allusers,
		Username:     Username,
		IsLoggedIn:   isLoggedIn,
		IsAdmin:      isAdmin,
		ErrorMessage: errorMessage,
	}
	if r.Method != http.MethodPost {
		err := tmplExped.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Failed to execute template:", err)
			return
		}
		return
	}
}
func sendColisHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté en vérifiant le cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		// Si l'utilisateur n'est pas connecté, rediriger vers la page de connexion
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Récupérer l'expéditeur depuis le cookie
	email := cookie.Value
	expediteur, _ := controller.GetUsernameByEmail(email)
	// Récupérer les informations du formulaire
	destinataire := r.FormValue("destinataire")
	titre := r.FormValue("titre")
	contenu := r.FormValue("contenu")
	var image []byte
	// Vérifier si le destinataire existe dans la base de données
	if !controller.UserExists(destinataire) {
		http.Redirect(w, r, "/expedition?error=Utilisateur inconnu", http.StatusFound)
		return
	}
	// Récupérer le fichier du formulaire s'il existe
	file, _, err := r.FormFile("fichier")
	if err == nil {
		defer file.Close()
		// Lire le contenu du fichier dans un tableau de bytes
		image, err = ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			log.Println("Failed to read file:", err)
			return
		}
	}
	// Insérer le colis dans la base de données
	colisID, err := controller.InsertColis(titre, contenu, image)
	if err != nil {
		http.Error(w, "Failed to create colis", http.StatusInternalServerError)
		log.Println("Failed to create colis:", err)
		return
	}
	// Insérer la livraison dans la base de données
	err = controller.InsertDelivery(expediteur, destinataire, colisID, getDateTimeFR())
	if err != nil {
		http.Error(w, "Failed to create delivery", http.StatusInternalServerError)
		log.Println("Failed to create delivery:", err)
		return
	}

	ID := controller.GetLastDeliveryIDForUser(destinataire)

	data := struct {
		ID string
	}{
		ID: ID,
	}

	err = tmplSuccesfuly.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func getDateTimeFR() string {
	// Obtenir la date et l'heure actuelles
	now := time.Now()
	// Définir la locale française
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return "Erreur lors du chargement de la localisation"
	}
	// Convertir la date et l'heure actuelles dans la locale française
	nowFR := now.In(loc)
	nowFR = nowFR.Add(24 * time.Hour)
	// Formater la date et l'heure au format français
	dateTimeFR := nowFR.Format("02 janvier 2006 15:04")
	return dateTimeFR
}
func getTrackingID(query string) string {
	params := strings.Split(query, "&")
	for _, param := range params {
		kv := strings.Split(param, "=")
		if kv[0] == "id" {
			return kv[1]
		}
	}
	return ""
}
