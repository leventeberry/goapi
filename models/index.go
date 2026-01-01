package models

type User struct {
    ID        int    `gorm:"primaryKey" json:"user_id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    PassHash  string `json:"-"` // Excluded from JSON responses for security
    PhoneNum  string `json:"phone_number"`
    Role      string `json:"role"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
