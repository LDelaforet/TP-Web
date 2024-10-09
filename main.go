package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

type Etudiant struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string
}

type Promotion struct {
	Nom         string
	Filiere     string
	Niveau      string
	NbEtudiants int
	Etudiants   []Etudiant
}

type PassedVars struct {
	Compteur     int
	CompteurPair string
}

type UserForm struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
	Errors        []string // Les erreur de saisie seront stockées ici
}

// J'aurais pu le mettre dans main mais je veut pas trop le charger
func promoInit() []Promotion {
	PromoList := []Promotion{}
	StudentList := []Etudiant{}

	StudentList = append(StudentList, Etudiant{
		Nom:    "Mederreg",
		Prenom: "Kheir-Eddine",
		Age:    22,
		Sexe:   "M",
	})

	StudentList = append(StudentList, Etudiant{
		Nom:    "Rodrigues",
		Prenom: "Cyril",
		Age:    23,
		Sexe:   "M",
	})

	Promo1 := Promotion{
		Nom:         "MENTOR",
		Filiere:     "Informatique",
		Niveau:      "Mentors",
		NbEtudiants: 2,
		Etudiants:   StudentList,
	}

	PromoList = append(PromoList, Promo1)

	return PromoList
}

func main() {
	PromoList := promoInit()
	currentPromo := PromoList[0]
	pVars := PassedVars{
		Compteur: 0,
	}
	Utilisateur := UserForm{}
	fmt.Println(PromoList)

	// Charger les templates HTML
	temp, errTemp := template.ParseGlob("webPages/*.html")
	if errTemp != nil {
		fmt.Printf("Error: %v\n", errTemp)
		return
	}

	// Bouton pour aller a /promo
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index", currentPromo)
	})

	// Route pour afficher les informations de la promotion
	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "promoDisplay", currentPromo)
	})

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		// Alors oui c'est totalement débile de faire comme ca mais j'y ait pensé en faisant le if et j'ai trouvé ca marrant
		pVars.CompteurPair = strings.Repeat("im", pVars.Compteur%2) + "pair"

		temp.ExecuteTemplate(w, "viewCount", pVars)
		pVars.Compteur++
	})
	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "userForm", Utilisateur)
	})
	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			Utilisateur.Errors = nil
			r.ParseForm()
			Utilisateur.Nom = r.FormValue("nom")
			Utilisateur.Prenom = r.FormValue("prenom")
			Utilisateur.DateNaissance = r.FormValue("dateNaissance")
			Utilisateur.Sexe = r.FormValue("sexe")
			if len(Utilisateur.Nom) < 1 || len(Utilisateur.Nom) > 32 {
				Utilisateur.Errors = append(Utilisateur.Errors, "Le nom doit être compris entre 1 et 32 caractères")
			}
			if len(Utilisateur.Prenom) < 1 || len(Utilisateur.Prenom) > 32 {
				Utilisateur.Errors = append(Utilisateur.Errors, "Le prénom doit être compris entre 1 et 32 caractères")
			}
			if Utilisateur.Errors != nil {
				http.Redirect(w, r, "/user/error", http.StatusSeeOther)
			}
			http.Redirect(w, r, "/user/display", http.StatusSeeOther)
			fmt.Println(Utilisateur)
		}
	})
	http.HandleFunc("/user/error", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "userError", Utilisateur)
	})
	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		Utilisateur.Errors = nil
		if Utilisateur.Nom == "" {
			Utilisateur.Errors = append(Utilisateur.Errors, "Veuillez renseigner votre nom")
		} else if len(Utilisateur.Nom) < 1 || len(Utilisateur.Nom) > 32 {
			Utilisateur.Errors = append(Utilisateur.Errors, "Le nom doit être compris entre 1 et 32 caractères")
		}
		if Utilisateur.Prenom == "" {
			Utilisateur.Errors = append(Utilisateur.Errors, "Veuillez renseigner votre prénom")
		} else if len(Utilisateur.Prenom) < 1 || len(Utilisateur.Prenom) > 32 {
			Utilisateur.Errors = append(Utilisateur.Errors, "Le prénom doit être compris entre 1 et 32 caractères")
		}
		if Utilisateur.DateNaissance == "" {
			Utilisateur.Errors = append(Utilisateur.Errors, "Veuillez renseigner votre date de naissance")
		}
		if Utilisateur.Sexe == "" {
			Utilisateur.Errors = append(Utilisateur.Errors, "Veuillez renseigner votre sexe")
		}

		if Utilisateur.Errors != nil {
			http.Redirect(w, r, "/user/error", http.StatusSeeOther)
		}
		temp.ExecuteTemplate(w, "userDisplay", Utilisateur)
	})
	// Fichiers statiques
	http.Handle("/webPages/", http.StripPrefix("/webPages/", http.FileServer(http.Dir("./webPages"))))
	http.Handle("/Pictures/", http.StripPrefix("/Pictures/", http.FileServer(http.Dir("./Pictures"))))

	// Lancer le serveur
	http.ListenAndServe("0.0.0.0:8080", nil)
}

/*
Challenge 1 : Affichage des données
Vous devez implémenter une route « /promo » qui permet d'afficher les informations
liées à une classe ainsi que la liste des étudiants qui la composent. Les informations
décrivant une classe sont les suivantes :
• Un nom de classe (ex : B1 Informatique),
• La filière (ex : Informatique),
• Le niveau (ex : Bachelor 1),
• Le nombre d'étudiants,
• La liste des étudiants avec les informations suivantes : nom, prénom, âge, et sexe
pour chaque étudiant.
Vous devez afficher l'ensemble des informations citées ci-dessus en utilisant des actions
(variables et boucles) pour faire passer les données. Vous êtes libre pour la mise en page
et son style. Le sexe doit être représenté par une image (masculin ou féminin), utilisez une
condition dans le template pour gérer cela
*/
