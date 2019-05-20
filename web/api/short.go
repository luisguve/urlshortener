package api

import(
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/villegasl/urlshortener/models"

	"github.com/gorilla/mux"
)

func RedirectByShortURL(DB_Handler *models.DBHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		var jsonRes []byte
		var err error

		vars := mux.Vars(r)
		shortURL := vars["url"]

		fmt.Println("\nshort url:",shortURL)

		status := models.GetUrl(shortURL, DB_Handler)

		if status.Error != nil {
			fmt.Println("Error while trying to query the database:", status.Error)
			jsonRes, err = json.Marshal(status.FailureStatus)
			if err != nil {
				fmt.Println("Error while trying to marshal the JSON response:", 
					err)
				return
			}
			w.Write(jsonRes)
			return
		}

		jsonRes, err = json.Marshal(status.SuccessStatus)
		if err != nil {
				fmt.Println("Error while trying to marshal the JSON response:", 
					err)
				return
			}
		fmt.Println("long url:",status.SuccessStatus.OriginalUrl)
		w.Write(jsonRes)
		return
	})
}

func NewShortURL(DB_Handler *models.DBHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var jsonRes []byte
		originalUrl := r.FormValue("url")

		// Please uncomment this url validation when good internet speed be available
		/*_, err := http.Get(originalUrl)
		if err != nil {
			fmt.Println("Error: Invalid URL:", err.Error())
			jsonRes, err := json.Marshal(models.FailureResponse { Error: "Invalid URL" })
			if err != nil {
				fmt.Println("Error while trying to marshal the JSON response:", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonRes)
			return
		}*/

		status := models.SaveUrl(originalUrl, DB_Handler)

		if status.Error != nil {
			fmt.Println("could not update the database:",status.Error)
			jsonRes, err := json.Marshal(status.FailureStatus)
			if err != nil {
				fmt.Println("Error while trying to marshal the JSON response:", 
					err)
				return
			}
			w.Write(jsonRes)
			return
		}
		// At this point the long url was shortened successfully

		// Respond to the client with appropiate JSON
		jsonRes, err := json.Marshal(status.SuccessStatus)
		if err != nil {
			fmt.Println("error while trying to marshal the JSON response:", err)
			return
		}
		w.Write(jsonRes)
	})
}