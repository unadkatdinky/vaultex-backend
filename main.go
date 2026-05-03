package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 1. Define the Database Model (GORM will auto-generate the table from this)
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ItemName  string         `gorm:"type:varchar(100);not null" json:"itemname"`
	Category  string         `gorm:"type:varchar(50)" json:"category"`
	Image     string         `gorm:"type:text" json:"image"`
	ItemProps datatypes.JSON `json:"itemprops"` 
}

var DB *gorm.DB

func main() {
	// 2. Connect to PostgreSQL
	// UPDATE "user" to match what worked in your DB App!
	dsn := "host=localhost user=dinkyunadkat dbname=vaultex port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. Check your DSN string:", err)
	}

	log.Println("Database connection successful!")

	// 3. Auto-Migrate: This tells GORM to look at the Product struct and build the SQL table automatically!
	DB.AutoMigrate(&Product{})

	// 4. Set up the Router
	r := gin.Default()

	// 🔓 Relaxed CORS for Local Development
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Tells Go to trust your React app no matter what IP it uses
		AllowMethods:    []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// 5. Define API Routes
	api := r.Group("/api")
	{
		// GET all products
		api.GET("/products", func(c *gin.Context) {
			var products []Product
			DB.Find(&products)
			c.JSON(http.StatusOK, products)
		})

		// POST a new product
		api.POST("/products", func(c *gin.Context) {
			var newProduct Product
			if err := c.ShouldBindJSON(&newProduct); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			DB.Create(&newProduct)
			c.JSON(http.StatusCreated, newProduct)
		})
	}

	// 6. Start the Server
	log.Println("Go server running on http://localhost:8080")
	r.Run(":8080")
}