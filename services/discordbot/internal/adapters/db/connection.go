package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect initializes the database connection
func Connect() error {
	host := getEnv("POSTGRES_HOST", "postgres")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "")
	dbname := getEnv("POSTGRES_DB", "owlscapital")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established successfully")
	return nil
}

// AutoMigrate runs automatic migrations for all models
func AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	// Migrate all tables
	err := DB.AutoMigrate(
		// Ambassador System
		&AmbassadorApplicant{},
		&Ambassador{},
		&Claim{},
		&Conversion{},
		&PayoutRequest{},
		&BalanceTransaction{},
		&WeeklyLeaderboard{},

		// Freemium System
		&FreemiumModule{},
		&FreemiumLesson{},
		&FreemiumAccessLog{},
		&FreemiumSetting{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	log.Println("✅ Database schema migrated successfully!")

	// Insert default freemium settings if not exists
	if err := insertDefaultSettings(); err != nil {
		log.Printf("⚠️ Warning: failed to insert default settings: %v", err)
	}

	return nil
}

// insertDefaultSettings adds default freemium settings if they don't exist
func insertDefaultSettings() error {
	defaults := []FreemiumSetting{
		{SettingKey: "forum_channel_id", SettingValue: "1447780920940695655"},
		{SettingKey: "default_footer_url", SettingValue: "attachment://footer.png"},
		{SettingKey: "lifetime_role_id", SettingValue: "718643316786462772"},
		{SettingKey: "premium_role_id", SettingValue: "885910828086362132"},
		{SettingKey: "free_role_id", SettingValue: "718643370301325404"},
		{SettingKey: "upgrade_message", SettingValue: "Upgrade to Premium to unlock all lessons!"},
		{SettingKey: "upgrade_url", SettingValue: "https://your-upgrade-link.com"},
	}

	for _, setting := range defaults {
		result := DB.Where("setting_key = ?", setting.SettingKey).FirstOrCreate(&setting)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

