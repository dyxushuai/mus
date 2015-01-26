package middlewares

import (
	"net/http"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/models"
	"fmt"
	"time"
	"go/token"
)

type Store interface {

}



func Auth(store *db.Storage) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				if err != nil {
					data, _ := models.NewErr("authentication failed").JSON()
					w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "application/json", "utf-8"))
					w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
					fmt.Fprintf(w, string(data))
				}
			}()
			params := r.URL.Query()
			token := params.Get("token")

			fmt.Println(time.Now())
			str, err := store.GetStr("mus:token")
			fmt.Println(time.Now())
			if err != nil {
				return
			}


			if str == token {
				h.ServeHTTP(w, r)
			}
			return
		})
	}
}
