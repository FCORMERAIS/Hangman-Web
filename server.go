package main

import (
	"html/template"
    "log"
    "net/http"
    "fmt"
    "io/ioutil"
    "math/rand"
    "time"
    "strings"
)

type Page struct { // la class Page est la classe permettant d'envoyer les variables que l'on souhaite dans notre fichier html et permet de l'afficher sur notre site avec un template 
	Letter    string
	Articles []string
	Word 	string
	Articles2 []string
	Vie int
	Image string
	End string
}

func main() {
	fs:= http.FileServer(http.Dir("tmpl/pos_hangman")) 
	http.Handle("/pos_hangman/", http.StripPrefix("/pos_hangman/" , fs)) // cela nous permet d'utiliser les images pour afficher la position du pendu en fonction du nombre de tentatives restantes 
	images_src := []string{"pos_hangman/pos_10.png","pos_hangman/pos_9.png","pos_hangman/pos_8.png","pos_hangman/pos_7.png","pos_hangman/pos_6.png","pos_hangman/pos_5.png","pos_hangman/pos_4.png","pos_hangman/pos_3.png","pos_hangman/pos_2.png","pos_hangman/pos_1.png","pos_hangman/pos_0.png"}
    // on stocks les endroits ou sont stockés les images du pendu dans uen liste 
	alphabet := []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}// liste contenant toutes les lettres que l'on peut utiliser
	word_tempo := chooseWord()// nous permet d'afficher la progression de l'utilisateur
	word := strings.Split(word_tempo, "")// il s'agit du mot a trouver mais qu'on a split pour qu'il soit contenu dans une liste (cela nous permet d'utiliser des focntions)
	attemps := 10 // initialisation du nombre de tentatives
	end := "" // initialisation du message qu'on affiche our savoir 
	letter := "" 
	letter_choose := take_letter(word) // nous permet d'enregistrer les lettres que l'utilisateur a choisie et il commence avec une lettre deja donné 
	for i:=0 ; i<len(alphabet);i++ {
		if letter_choose[0] == alphabet[i] {
			alphabet = remove(alphabet,i) // on retire la lettre que l'on a donné de l'alphabet
		}
	}
	// 
    http.HandleFunc("/Hangman", func(w http.ResponseWriter, r *http.Request) {
        // Création d'une page
        letter = strings.ToUpper(r.FormValue("letter")) // on stock la lettre que l'uilisateur a rentré 
        for i:=0 ; i<len(alphabet);i++ { // debut for
            if letter == alphabet[i] {//debut if 
                letter_choose = append(letter_choose,alphabet[i]) 
				alphabet = remove(alphabet,i)
				if letterChooseTest(letter, word) == false {
					if attemps > 0 {
						attemps-- // on retire une vie si la lettre n'est pas contenu dans le mot et que la lettre n'a pas déja été rentré
					}
				}
            } // fin if
        } // fin for 
		if win(word,letter_choose) && attemps > 0 { // on verifie si l'utilisateur a gagné et a plus de 0 tentatives restantes
			end = "You Win" // on affiche "You Win"
			alphabet = []string{}
			letter_choose = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"} // on change les valeurs
		}else if attemps == 0 { // l'utilisateur a perdu
			end = "You Loose"
			alphabet = []string{}
			letter_choose = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
		}
		word_tempo = printWord(letter_choose,word)
        p := Page{letter, alphabet,word_tempo,letter_choose,attemps,images_src[attemps],end}// Création d'une nouvelle  de template
        t := template.New("Label de ma template")// Déclaration des fichiers à parser
        t = template.Must(t.ParseFiles("tmpl/layout.html", "tmpl/content.html"))// Exécution de la fusion et injection dans le flux de sortie / La variable p sera réprésentée par le "." dans le layout
        err := t.ExecuteTemplate(w, "layout", p)
        if err != nil {
			error501()
            log.Fatalf("Template execution: %s", err)
        }
    })
	//
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //crée une page avec n'importe quel URL 
		tmpl, err := template.ParseFiles("./tmpl/index.html")
		tmpl.ExecuteTemplate(w, "index", nil)
		if err != nil {
			error501()
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}

func error404() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //crée une page avec n'importe quel URL si il y a une erreur 404
		tmpl, err := template.ParseFiles("./error404.html")
		tmpl.ExecuteTemplate(w, "error404", nil)
		if err != nil {
			error501()
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}

func error501() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //crée une page avec n'importe quel URL si il y a une erreur 501
		tmpl, err := template.ParseFiles("./error501.html")
		tmpl.ExecuteTemplate(w, "error501", nil)
		if err != nil {
			// error501()
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}
func remove(slice []string, s int) []string {
	/*
	fonction permettant de retirer n'importe quel élément a l'index "s" que l'on souhaite
	*/
    return append(slice[:s], slice[s+1:]...)
}

func printWord(letter_choose []string, word []string)string {
	/*
	fonction permettant d'afficher le mot en fonction des lettres que l'utilisateur a déjà trouver
	input : -letter_choose type []string il s'agit des lettres que l'utilisateur a déjà rentrer 
			-word : type string il s'agit du mot a deviner 
	copléxité O(n²)
	*/
	mot := ""
	count := 0 
	for i:= 0 ; i < len(word) ;i++{
		for k := 0; k < len(letter_choose); k++ { // debut de la boucle
			if string(word[i]) == string(letter_choose[k]) {
				mot = mot+string(word[i])
				count++
			}
		} //fin de la boucle
		if count == 0 {
			mot=mot+"_" // renvoie un underscore si le joueur n'a pas trouver la lettre 
		}
		count = 0
		mot = mot + " "
	} 
	return mot
}

func chooseWord() string {
	/*
	fonction permettant de prendre un mot aléatoire dans une banque de mots 
	return : un mot aléatoire de type sring
	complexité : O(n) ; n = nombre de caractère de tout les mots reunis 
	*/
	s, err := ioutil.ReadFile("words.txt") // ouverture du fichier word1 contenant tout les mots et qui seront stockés dans s
	if err != nil {
		error404()
		fmt.Println(err.Error()) // renvois de l'erreur lors de louverture du fichier si il y a un problème 
		fmt.Println(" Le fichier word1.txt a planter...")
	}else {
		list := strings.Split(string(s),"\n") 
		rand.Seed(time.Now().UnixNano()) // ceci permet de faire en sorte que l'aléatoire marche 
		return strings.ToUpper(list[rand.Int31n(83)]) // renvois un mot choisis aléatoirement dans notre liste de mots 
	}
	return " "
}

func take_letter(word2 []string) []string{
	/*
	fonction permettant de donner une lettre qui est présente dans le mot a deviner 
	input : word2 type string il s'agit du mot a deviner 
	return : List/string il s'agit de la liste des lettres choisie par l'utlisateur 
	compléxité : O(n) ; n = longueur du mot
	*/
	var tab []string
	rand.Seed(time.Now().UnixNano())
	tab = append(tab,string(word2[rand.Intn(len(word2)-1)])) // prend une lettre aléatoire dans le mot a deviner et l'ajoute dans les lettres déja choisie 
	for i := 0; i < len(word2); i++ {
		if string(word2[i]) == "-" { // si il y a un tiret l'ajoute aussi dans notre liste 
			tab = append(tab,"-") 
			return tab 
		}
	}
	return tab // renvois le tableau contenant 
}

func win(wordChoose []string, group_letter []string) bool {
	/*
	fonction qui permet de savoir si l'utilisateur a gagné en trouver toutes les lettres 
	input : -wordChoose type string il sagit du mot que l'utilisateur doit choisir 
			-group_letter type List/string il s'agit des lettres que l'utilisateur a choisi 
	return : bool 
	Complexité : O(n²)
	*/
	count := 0 
	for i := 0; i < len(wordChoose); i++ { // debut boucle i
		for k := 0; k < len(group_letter); k++ { // debut boucle k 
			if group_letter[k] == string(wordChoose[i]) { // on vérifie si une lettre du mot est dans notre liste de mot 
				count++ // si la lettre du mot est dans notre liste de lettre on ajoute 1 a un compteur 
			}// fin obucle k
		} // fin boucle i
	}
	if count == len(wordChoose) { // si notre compteur est égal a la longueur du mot a deviner alors nous avons trouver toutes les lettres
		return true
	}else { // sinon non
		return false
	}
}

func letterChooseTest(letter string, word []string) bool {
	/*
	fonction permettant de vérifier si la lettre chosi par l'utilisateur est contenu dans le mot a deviner
	input : -letter : type string il sagit de la lettre choisi par l'utilisateur
			-word : type string il sagit du mot a deviner
	return : Bool
	complexité : O(2n) ; n = len(word)
	*/
	for i := 0; i < len(word); i++ {
		if letter == string(word[i]) { // on vérifie si la lettre choisi par l'utilisateur est présente dans le mot a deviner si c'est le cas en renvois true 
			return true
		}
	}
	return false
}