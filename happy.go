package happy

import "net/http"

func Run() {
	http.ListenAndServe(":8044", &RouterMux {})
}

