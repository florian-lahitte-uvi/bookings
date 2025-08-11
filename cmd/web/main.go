package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/florian-lahitte-uvi/bookings/helpers"
	"github.com/florian-lahitte-uvi/bookings/internal/config"
	"github.com/florian-lahitte-uvi/bookings/internal/driver"
	"github.com/florian-lahitte-uvi/bookings/internal/handlers"
	"github.com/florian-lahitte-uvi/bookings/internal/models"
	"github.com/florian-lahitte-uvi/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting Mail Listener")
	listenForMail()

	// msg := models.MailData{
	// 	From:    "from@example.com",
	// 	To:      "to@example.com",
	// 	Subject: "Test Email",
	// 	Content: "Hello, this is a test email",
	// }

	// app.MailChan <- msg

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// Read flags
	inProduction := flag.Bool("production", true, "Application in production mode")
	useCache := flag.Bool("cache", true, "Use cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL mode")

	flag.Parse()

	if *dbName == "" || *dbUser == "" || *dbPass == "" {
		fmt.Println("Database name, user, and password must be provided.")
		os.Exit(1)
	}

	// create a channel for mail data
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s",
		*dbUser, *dbPass, *dbName, *dbPort, *dbHost, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	log.Println("Connected to database!")

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
