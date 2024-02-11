package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

// Suggest represents data about a suggestion
type Suggest struct {
	Fname string `bson:"fname" json:"fname"`
	Lname string `bson:"lname" json:"lname"`
	Email string `bson:"email" json:"email"`
	Sugg  string `bson:"sugg" json:"sugg"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	//Establish a connection to MongoDB
	client, err := connection()
	if err != nil {
		log.Fatal(err)
		return
	}
	mongoClient = client
	defer mongoClient.Disconnect(context.TODO())

	router := gin.Default()

	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5500/api/submit"} // Replace with your frontend's domain
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	//define API routes
	router.GET("/api/Suggest", getSuggestions)
	router.POST("/api/submit", submitSuggestion)
	// Serve static assets directly using Gin

	router.StaticFile("/", "./index.html")
	router.StaticFile("/edit", "./playground/edit.html")
	router.StaticFile("/view", "./view.html")
	router.Static("/assets", "./assets")

	/*
		// Define dynamic routes for handling views and edits
		router.GET("Projects/Go-Wonder/playground/view/", makeHandler(viewHandler))
		router.GET("/playground/edit/", makeHandler(editHandler))
		router.POST("/save/", makeHandler(saveHandler))
		router.GET("/", indexHandler)
	*/

	// Create a custom http.Server with the Gin router as the handler
	server := &http.Server{
		Addr:    ":5500",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
		}
	}()
	/// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}

func connection() (*mongo.Client, error) {
	// MongoDB connection string
	connectionString := "mongodb+srv://umbra:password1995@umbramortis.m70inkz.mongodb.net/?retryWrites=true&w=majority"

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	//connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}

// getSuggestions responds with the list of all suggestions from MongoDB as JSON.
func getSuggestions(c *gin.Context) {
	db := mongoClient.Database("Suggestions")
	coll := db.Collection("submissions")

	// Query the MongoDB collection for all documents (suggestions).
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.TODO())

	// Decode the retrieved documents into a slice of Suggest structs.
	var suggestions []Suggest
	for cursor.Next(context.TODO()) {
		var suggestion Suggest
		if err := cursor.Decode(&suggestion); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		suggestions = append(suggestions, suggestion)
	}

	// Respond with the retrieved suggestions as JSON.
	c.JSON(http.StatusOK, suggestions)
}

// PostSuggestions adds a suggestion from JSON received in the request body.

func submitSuggestion(c *gin.Context) {

	db := mongoClient.Database("Suggestions")
	coll := db.Collection("submissions")

	var newSugg Suggest

	// Parse JSON data from the request body into the newSugg struct
	if err := c.ShouldBindJSON(&newSugg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	_, err := coll.InsertOne(context.TODO(), newSugg)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, newSugg)
}

/*
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(c *gin.Context) {
	// Extract the "title" parameter from the URL
	title := c.Param("title")

	// Load the page based on the "title" parameter
	p, err := loadPage(title)
	if err != nil {
		// Handle the error, for example:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Render the "view" template with the page data
	renderTemplate(c, "view", p)
}
func editHandler(c *gin.Context) {
	title := c.Param("title")
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(c, "edit", p)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(c *gin.Context, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(c.Writer, tmpl+".html", p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}
}

func saveHandler(c *gin.Context) {
	title := c.Param("title")
	body := c.PostForm("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/view/"+title)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := validPath.FindStringSubmatch(c.Request.URL.Path)
		if m == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		fn(c)
	}
}
func indexHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/playground/view/")
}
*/
