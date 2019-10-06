package routers

import (
	"fmt"
	"net/http"
	"time"
	// "../models"
)

// logRequest - логгер/принтер в консоль данных о запросе
func logRequest(r *http.Request) {
	const dtFormat = "2006-01-02 03:04:05"
	correntTime := time.Now().Format(dtFormat)
	fmt.Println(r.Method, r.RemoteAddr, correntTime, r.RequestURI, r.Header["User-Agent"][0])
}

// IndexRoute = is '/'
// Передача параметров в template
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	// logRequest(r)
	// bearer := strings.ContainsAny(r.Header.Get("Authorization"), "Bearer")
	// if bearer {
	// 	bearer := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	// 	user, err := ValidateToken(bearer)
	// 	if err != nil {
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}
	// 	fmt.Println(user)
	// }

	validRoles := []string{"moderator", "admin"}
	session, _ := store.Get(r, "user")
	user := getUser(session)
	user.checkRole(validRoles, w, r)
	message := fmt.Sprintf("%s you are %s!", "Hello", user.Role)
	w.Write([]byte(message))

}
