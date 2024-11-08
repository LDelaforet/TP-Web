package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// Oui c'est completement con mais jsp comment faire autrement
type Taille struct {
	Taille string
}

type Produit struct {
	Id                 int
	ImageName          string
	Nom                string
	Description        string
	Prix               string
	Reduction          string
	TaillesDisponibles []Taille
}

type PassedData struct {
	ProductList     []Produit
	SelectedProduct Produit
	Error           string
}

var SelectedProduct Produit

func main() {
	ProductList := []Produit{}
	ProductList = ProductsFiller(ProductList)
	PData := PassedData{
		ProductList: ProductList,
	}

	var SelectedProduct Produit

	// Charger les templates HTML
	temp, errTemp := template.ParseGlob("templates/*.html")
	if errTemp != nil {
		fmt.Printf("Error: %v\n", errTemp)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "productList", PData)
	})

	http.HandleFunc("/product/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		prodIdStr := path[len("/product/"):]

		prodID, _ := strconv.Atoi(prodIdStr)

		for _, product := range PData.ProductList {
			if product.Id == prodID {
				SelectedProduct = product
				break
			}
		}
		temp.ExecuteTemplate(w, "productData", SelectedProduct)
	})

	http.HandleFunc("/productMgmt", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "productMgmt", PData)
	})

	http.HandleFunc("/productMgmt/newProduct", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()

			if r.FormValue("nom") == "" || r.FormValue("description") == "" || r.FormValue("prix") == "" || (r.FormValue("size1") == "" && r.FormValue("size2") == "" && r.FormValue("size3") == "" && r.FormValue("size4") == "") {
				PData.Error = "Veuillez remplir tous les champs"
				temp.ExecuteTemplate(w, "productMgmt", PData)
			} else {
				PData.Error = ""
			}

			Tailles := []Taille{}
			if r.FormValue("size1") != "" {
				Tailles = append(Tailles, Taille{Taille: "S"})
			}
			if r.FormValue("size2") != "" {
				Tailles = append(Tailles, Taille{Taille: "M"})
			}
			if r.FormValue("size3") != "" {
				Tailles = append(Tailles, Taille{Taille: "L"})
			}
			if r.FormValue("size4") != "" {
				Tailles = append(Tailles, Taille{Taille: "XL"})
			}
			fmt.Printf("Tailles: %v\n", Tailles)
			newProduct := Produit{
				Id:                 len(PData.ProductList),
				ImageName:          "19A.webp",
				Nom:                strings.ToUpper(r.FormValue("nom")),
				Description:        r.FormValue("description"),
				Prix:               r.FormValue("prix") + "€",
				TaillesDisponibles: Tailles,
				Reduction:          "0",
			}
			PData.ProductList = append(PData.ProductList, newProduct)
			productPage := "/product/" + strconv.Itoa(newProduct.Id)
			http.Redirect(w, r, productPage, http.StatusSeeOther)
		}
	})

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/img/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./assets/styles/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./assets/fonts/"))))

	// 0.0.0.0 >>> 127.0.0.1, ici on adore lancer sur ttes les interfaces
	http.ListenAndServe("0.0.0.0:8080", nil)
}

// Psq ca me saoulait de voir tt ce bordel au dessus de main
func ProductsFiller(ProductList []Produit) []Produit {
	ProductList = append(ProductList, Produit{
		Id:          0,
		ImageName:   "19A.webp",
		Nom:         "PALACE PULL A CAPUCHE UNISEX CHASSEUR",
		Description: "Description du produit",
		Prix:        "146€",
		Reduction:   "0",
		TaillesDisponibles: []Taille{
			Taille{Taille: "S"},
			Taille{Taille: "M"},
			Taille{Taille: "XL"},
		},
	})
	ProductList = append(ProductList, Produit{
		Id:        1,
		ImageName: "21A.webp",
		Nom:       "PALACE PULL A CAPUCHON MARINE",
		Prix:      "138€",
		Reduction: "0",
		TaillesDisponibles: []Taille{
			Taille{Taille: "S"},
			Taille{Taille: "L"},
			Taille{Taille: "XL"},
		},
	})
	ProductList = append(ProductList, Produit{
		Id:          2,
		ImageName:   "22A.webp",
		Nom:         "PALACE PULL CREW PASSEPOSE NOIR",
		Description: "Description du produit",
		Prix:        "128€",
		Reduction:   "48€",
		TaillesDisponibles: []Taille{
			Taille{Taille: "M"},
			Taille{Taille: "L"},
			Taille{Taille: "XL"},
		},
	})
	ProductList = append(ProductList, Produit{
		Id:          3,
		ImageName:   "16A.webp",
		Nom:         "PALACE WASHED TERRY 1/4 PLACKET HOOD MOJITO",
		Description: "Description du produit",
		Prix:        "168€",
		Reduction:   "0",
		TaillesDisponibles: []Taille{
			Taille{Taille: "S"},
			Taille{Taille: "M"},
			Taille{Taille: "L"},
			Taille{Taille: "XL"},
		},
	})
	ProductList = append(ProductList, Produit{
		Id:          4,
		ImageName:   "34B.webp",
		Nom:         "PALACE PANTALON BOSSY JEAN STONE",
		Description: "Description du produit",
		Prix:        "125€",
		Reduction:   "0",
		TaillesDisponibles: []Taille{
			Taille{Taille: "S"},
			Taille{Taille: "XL"},
		},
	})
	ProductList = append(ProductList, Produit{
		Id:          5,
		ImageName:   "33B.webp",
		Nom:         "PALACE PANTALON CARGO GORE-TEX R-TEK NOIR",
		Description: "Description du produit",
		Prix:        "110€",
		Reduction:   "0",
		TaillesDisponibles: []Taille{
			Taille{Taille: "S"},
		},
	})
	return ProductList
}
