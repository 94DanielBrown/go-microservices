package main

const webPort = "80"

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main () {
	log.Println("Starting authentication service")

	// TODO connect to DB

	app := Config{}

	srv := &http.Server{
		Addr : fmt.Sprintf(":%s", webPort)
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}