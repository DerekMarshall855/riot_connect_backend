package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/item.json")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))
}

func getChampions(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/champion.json")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))
}

func getChampionByName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	//Logic used for testing, allows me to use any format for champion name. May put conversion on frontend but does not matter currently
	champion := vars["name"]

	champion = strings.ToLower(champion)
	temp := []byte(champion)
	temp[0] = byte(unicode.ToUpper(rune(temp[0])))
	champion = string(temp)
	fmt.Fprintf(w, "Champion Name: %s\n", vars["name"])
	//End of logic

	// Champ name must have capital first letter followed by lowercase
	res, err := http.Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/champion/%s.json", champion))

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))

}

func getAccountInfo(w http.ResponseWriter, r *http.Request) {
	api_key := "RGAPI-666d2caa-bf6c-41b7-9388-958990b9f333" // Stupid, only works for 24/h, use .env variable later

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	//Logic used for testing, allows me to use any format for champion name. May put conversion on frontend but does not matter currently
	user := vars["name"]

	user = strings.ToLower(user)
	temp := []byte(user)
	temp[0] = byte(unicode.ToUpper(rune(temp[0])))
	user = string(temp)
	fmt.Fprintf(w, "user Name: %s\n", vars["name"])
	//End of logic

	// Currently just gets account info for NA1, later should be able to select region
	res, err := http.Get(fmt.Sprintf("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", user, api_key))

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))
}

// Use https://github.com/gorilla/mux to handle variable api requests such as:
// https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/champion/Aatrox.json

// `https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/champion.json`
// `https://ddragon.leagueoflegends.com/cdn/11.15.1/data/en_US/item.json`
// `https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/Doublelift?api_key=RGAPI-cb48e359-5533-4305-bc7c-895c7ba24ee7`
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/items", getItems)
	r.HandleFunc("/champions", getChampions)
	r.HandleFunc("/champion/{name}", getChampionByName)
	r.HandleFunc("/account/{name}", getAccountInfo)
	http.Handle("/", r)

	log.Fatal(
		http.ListenAndServe(":8080",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(r),
		),
	)
}
