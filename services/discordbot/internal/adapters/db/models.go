package db

import (
	"time"
)

// ============================================
// AMBASSADOR SYSTEM MODELS
// ============================================

// AmbassadorApplicant stores users who have been approved but not yet accepted as full ambassadors
type AmbassadorApplicant struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Username    string    `gorm:"type:varchar(100);not null"`
	DisplayName string    `gorm:"type:varchar(100)"`
	ApprovedBy  string    `gorm:"type:varchar(20);not null"`
	ApprovedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ThreadID    string    `gorm:"type:varchar(20)"`
	Status      string    `gorm:"type:varchar(20);default:'pending';index"`
	Notes       string    `gorm:"type:text"`
}

// Ambassador stores fully accepted ambassadors with their stats and balance
type Ambassador struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Username           string    `gorm:"type:varchar(100);not null"`
	DisplayName        string    `gorm:"type:varchar(100)"`
	PrivateChannelID   string    `gorm:"type:varchar(20)"`
	AcceptedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AcceptedBy         string    `gorm:"type:varchar(20);not null"`
	
	// Stats
	TotalClaims         int `gorm:"default:0"`
	TotalUnclaims       int `gorm:"default:0"`
	TotalConversions    int `gorm:"default:0;index:idx_ambassador_conversions"`
	LifetimeConversions int `gorm:"default:0;index:idx_ambassador_lifetime"`
	CurrentClaims       int `gorm:"default:0"`
	
	// Balance
	Balance       float64 `gorm:"type:decimal(10,2);default:0.00"`
	TotalEarned   float64 `gorm:"type:decimal(10,2);default:0.00"`
	TotalPaidOut  float64 `gorm:"type:decimal(10,2);default:0.00"`
	
	// Tier system
	Tier string `gorm:"type:varchar(20);default:'bronze'"`
	
	// Cooldowns
	LastUnclaimAt *time.Time
	
	// Status
	IsActive bool `gorm:"default:true"`
}

// Claim tracks all claim relationships between ambassadors and free users
type Claim struct {
	ID               uint      `gorm:"primaryKey"`
	AmbassadorID     string    `gorm:"type:varchar(20);not null;index:idx_claims_ambassador"`
	ClaimedUserID    string    `gorm:"type:varchar(20);not null;index:idx_claims_claimed_user;uniqueIndex:idx_unique_claim"`
	ClaimedUsername  string    `gorm:"type:varchar(100)"`
	ClaimedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status           string    `gorm:"type:varchar(20);default:'active';index:idx_claims_status"`
	ConvertedAt      *time.Time
	UnclaimedAt      *time.Time
}

// Conversion logs all conversions for historical tracking
type Conversion struct {
	ID                 uint      `gorm:"primaryKey"`
	AmbassadorID       string    `gorm:"type:varchar(20);not null;index:idx_conversions_ambassador"`
	ConvertedUserID    string    `gorm:"type:varchar(20);not null"`
	ConvertedUsername  string    `gorm:"type:varchar(100)"`
	ConvertedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Commission         float64   `gorm:"type:decimal(10,2);default:0.00"`
	ClaimID            *uint
	WeekNumber         int       `gorm:"index:idx_conversions_week"`
	Year               int       `gorm:"index:idx_conversions_week"`
}

// PayoutRequest tracks payout requests and their status
type PayoutRequest struct {
	ID             uint      `gorm:"primaryKey"`
	AmbassadorID   string    `gorm:"type:varchar(20);not null;index:idx_payouts_ambassador"`
	Amount         float64   `gorm:"type:decimal(10,2);not null"`
	PaymentMethod  string    `gorm:"type:varchar(100)"`
	ThreadID       string    `gorm:"type:varchar(20)"`
	Status         string    `gorm:"type:varchar(20);default:'pending';index:idx_payouts_status"`
	RequestedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CompletedAt    *time.Time
	CompletedBy    string    `gorm:"type:varchar(20)"`
	ActualAmount   *float64  `gorm:"type:decimal(10,2)"`
	Notes          string    `gorm:"type:text"`
}

// BalanceTransaction audit log of all balance changes
type BalanceTransaction struct {
	ID           uint      `gorm:"primaryKey"`
	AmbassadorID string    `gorm:"type:varchar(20);not null;index:idx_transactions_ambassador"`
	Type         string    `gorm:"type:varchar(20);not null;index:idx_transactions_type"`
	Amount       float64   `gorm:"type:decimal(10,2);not null"`
	BalanceAfter float64   `gorm:"type:decimal(10,2);not null"`
	Description  string    `gorm:"type:text"`
	PerformedBy  string    `gorm:"type:varchar(20)"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// WeeklyLeaderboard snapshot of weekly stats
type WeeklyLeaderboard struct {
	ID           uint      `gorm:"primaryKey"`
	AmbassadorID string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_weekly_unique"`
	WeekNumber   int       `gorm:"not null;uniqueIndex:idx_weekly_unique"`
	Year         int       `gorm:"not null;uniqueIndex:idx_weekly_unique"`
	Username     string    `gorm:"type:varchar(100)"`
	Conversions  int       `gorm:"default:0"`
	Earned       float64   `gorm:"type:decimal(10,2);default:0.00"`
	ResetAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// ============================================
// FREEMIUM SYSTEM MODELS
// ============================================

// FreemiumModule represents an educational module
type FreemiumModule struct {
	ID              uint      `gorm:"primaryKey"`
	ModuleID        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Title           string    `gorm:"type:varchar(100);not null"`
	Description     string    `gorm:"type:text;not null"`
	BannerURL       string    `gorm:"type:varchar(500)"`
	ForumChannelID  string    `gorm:"type:varchar(20);not null"`
	ForumPostID     string    `gorm:"type:varchar(20);index:idx_modules_forum_post"`
	EmbedMessageID  string    `gorm:"type:varchar(20)"`
	CreatedBy       string    `gorm:"type:varchar(20);not null"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	IsPublished     bool      `gorm:"default:false;index:idx_modules_published,where:is_published = true"`
	IsArchived      bool      `gorm:"default:false"`
	
	// Relationships
	Lessons []FreemiumLesson `gorm:"foreignKey:ModuleID;references:ModuleID;constraint:OnDelete:CASCADE"`
}

// FreemiumLesson represents a lesson within a module
type FreemiumLesson struct {
	ID              uint      `gorm:"primaryKey"`
	ModuleID        string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_lesson_unique;index:idx_lessons_order"`
	LessonID        string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_lesson_unique"`
	Title           string    `gorm:"type:varchar(100);not null"`
	Description     string    `gorm:"type:text"`
	Content         string    `gorm:"type:text;not null"`
	ContentImageURL string    `gorm:"type:varchar(500)"`
	OrderIndex      int       `gorm:"not null;uniqueIndex:idx_lesson_order;index:idx_lessons_order"`
	IsFree          bool      `gorm:"->:false;<-:create;default:false"` // Computed: order_index = 1
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// FreemiumAccessLog tracks user access to lessons
type FreemiumAccessLog struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        string    `gorm:"type:varchar(20);not null;index:idx_access_user"`
	ModuleID      string    `gorm:"type:varchar(50);not null;index:idx_access_module"`
	LessonID      string    `gorm:"type:varchar(50);not null"`
	AccessGranted bool      `gorm:"not null"`
	UserRoleType  string    `gorm:"type:varchar(20);not null"`
	AccessedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;index:idx_access_time"`
}

// FreemiumSetting stores system-wide freemium settings
type FreemiumSetting struct {
	ID           uint      `gorm:"primaryKey"`
	SettingKey   string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	SettingValue string    `gorm:"type:text;not null"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedBy    string    `gorm:"type:varchar(20)"`
}

// TableName overrides (optional, GORM will use snake_case by default)
func (AmbassadorApplicant) TableName() string  { return "ambassador_applicants" }
func (Ambassador) TableName() string           { return "ambassadors" }
func (Claim) TableName() string                { return "claims" }
func (Conversion) TableName() string           { return "conversions" }
func (PayoutRequest) TableName() string        { return "payout_requests" }
func (BalanceTransaction) TableName() string   { return "balance_transactions" }
func (WeeklyLeaderboard) TableName() string    { return "weekly_leaderboard" }
func (FreemiumModule) TableName() string       { return "freemium_modules" }
func (FreemiumLesson) TableName() string       { return "freemium_lessons" }
func (FreemiumAccessLog) TableName() string    { return "freemium_access_log" }
func (FreemiumSetting) TableName() string      { return "freemium_settings" }

