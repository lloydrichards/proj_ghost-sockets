package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	migrate() error

	// GetUsers returns a list of Users.
	GetUsers() ([]User, error)

	CreateUser(User User) error
	UpdateUser(User User) error
	GetUser(username string) (User, error)

	CreateSession(session Session) error
	UpdateSession(sessionId uuid.UUID) error
	GetLatestSession(username string) (Session, error)
	ResetAllSessions() error
}

type service struct {
	db *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	fmt.Println("Creating new database connection...")
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	// migrate the database
	dbInstance.migrate()

	fmt.Println("Connected to database:", database)

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	db, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// Ping the database
	err = db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (s *service) GetUsers() ([]User, error) {
	var users []User
	result := s.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (s *service) CreateUser(user User) error {
	result := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)

	return result.Error
}

func (s *service) UpdateUser(User User) error {
	result := s.db.Save(&User)
	return result.Error
}

func (s *service) GetUser(username string) (User, error) {
	var user User
	result := s.db.Where("name = ?", username).First(&user)

	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (s *service) CreateSession(session Session) error {
	result := s.db.Create(&session)
	return result.Error
}

func (s *service) UpdateSession(sessionId uuid.UUID) error {
	result := s.db.Model(&Session{}).Where("id = ?", sessionId).Update("is_active", false)
	return result.Error
}

func (s *service) GetLatestSession(username string) (Session, error) {
	var session Session
	result := s.db.Where("user_name = ?", username).Order("created_at desc").First(&session)

	if result.Error != nil {
		return Session{}, result.Error
	}
	return session, nil
}

func (s *service) ResetAllSessions() error {
	result := s.db.Model(&Session{}).Update("is_active", false)
	return result.Error
}
