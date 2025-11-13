package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/infrastructure"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// generateRandomPassword generates a random password with uppercase, lowercase, digits, and special characters
func generateRandomPassword(length int) (string, error) {
	const (
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		digits           = "0123456789"
		specialChars     = "!@#$%^&*"
		allChars         = uppercaseLetters + lowercaseLetters + digits + specialChars
	)

	if length < 12 {
		length = 12 // Minimum secure password length
	}

	password := make([]byte, length)

	// Ensure at least one character from each category
	categories := []string{uppercaseLetters, lowercaseLetters, digits, specialChars}
	for i, category := range categories {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(category))))
		if err != nil {
			return "", err
		}
		password[i] = category[n.Int64()]
	}

	// Fill the rest with random characters from all categories
	for i := len(categories); i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[n.Int64()]
	}

	// Shuffle the password to avoid predictable patterns
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(password))))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// Seed users for each role type
func main() {
	// Get database connection string from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "algorithmia")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "algorithmia")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Get all roles
	var roles []database.Role
	if err := db.WithContext(ctx).Preload("Permissions").Find(&roles).Error; err != nil {
		log.Fatalf("Failed to load roles: %v", err)
	}

	// Map roles by name
	roleMap := make(map[string]database.Role)
	for _, role := range roles {
		roleMap[role.Name] = role
	}

	fmt.Println("🌱 Starting user seeding...")
	fmt.Println("==========================================")

	// Define seed users for each role (passwords will be generated randomly)
	seedUserTemplates := []struct {
		Username string
		Email    string
		Roles    []string
		Display  string
	}{
		{
			Username: "admin",
			Email:    "admin@algorithmia.com",
			Roles:    []string{"super_admin"},
			Display:  "Super Administrator",
		},
		{
			Username: "setter1",
			Email:    "setter1@algorithmia.com",
			Roles:    []string{"setter"},
			Display:  "Alice Chen - Problem Setter",
		},
		{
			Username: "setter2",
			Email:    "setter2@algorithmia.com",
			Roles:    []string{"setter"},
			Display:  "Bob Wang - Problem Setter",
		},
		{
			Username: "reviewer1",
			Email:    "reviewer1@algorithmia.com",
			Roles:    []string{"reviewer"},
			Display:  "Carol Li - Problem Reviewer",
		},
		{
			Username: "reviewer2",
			Email:    "reviewer2@algorithmia.com",
			Roles:    []string{"reviewer"},
			Display:  "David Zhang - Problem Reviewer",
		},
		{
			Username: "tester1",
			Email:    "tester1@algorithmia.com",
			Roles:    []string{"tester"},
			Display:  "Eva Liu - Problem Tester",
		},
		{
			Username: "tester2",
			Email:    "tester2@algorithmia.com",
			Roles:    []string{"tester"},
			Display:  "Frank Wu - Problem Tester",
		},
		{
			Username: "tester3",
			Email:    "tester3@algorithmia.com",
			Roles:    []string{"tester"},
			Display:  "Grace Huang - Problem Tester",
		},
		{
			Username: "contest_mgr",
			Email:    "contest@algorithmia.com",
			Roles:    []string{"contest_manager"},
			Display:  "Henry Chen - Contest Manager",
		},
		{
			Username: "multi_role",
			Email:    "multi@algorithmia.com",
			Roles:    []string{"setter", "reviewer", "tester"},
			Display:  "Iris Yang - Multi-Role User",
		},
	}

	// Generate random passwords for each user
	seedUsers := make([]SeedUser, 0, len(seedUserTemplates))
	for _, template := range seedUserTemplates {
		password, err := generateRandomPassword(16)
		if err != nil {
			log.Fatalf("Failed to generate password for %s: %v", template.Username, err)
		}
		seedUsers = append(seedUsers, SeedUser{
			Username: template.Username,
			Email:    template.Email,
			Password: password,
			Roles:    template.Roles,
			Display:  template.Display,
		})
	}

	for _, seed := range seedUsers {
		// Check if user already exists
		var existingUser database.User
		err := db.WithContext(ctx).Where("email = ?", seed.Email).First(&existingUser).Error
		if err == nil {
			fmt.Printf("⏭️  User %s (%s) already exists, skipping...\n", seed.Username, seed.Email)
			continue
		}

		// Hash password using Argon2 (same as the backend)
		hasher := infrastructure.NewArgonPasswordHasher()
		hashedPassword, err := hasher.Hash(seed.Password)
		if err != nil {
			log.Printf("❌ Failed to hash password for %s: %v", seed.Username, err)
			continue
		}

		// Create user
		userID, err := uuid.NewV7()
		if err != nil {
			log.Printf("❌ Failed to generate UUID for %s: %v", seed.Username, err)
			continue
		}

		// Get roles for this user
		userRoles := make([]database.Role, 0)
		for _, roleName := range seed.Roles {
			if role, ok := roleMap[roleName]; ok {
				userRoles = append(userRoles, role)
			}
		}

		user := database.User{
			UserID:         userID,
			Username:       seed.Username,
			Email:          seed.Email,
			DisplayName:    seed.Display,
			HashedPassword: hashedPassword,
			Roles:          userRoles,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			log.Printf("❌ Failed to create user %s: %v", seed.Username, err)
			continue
		}

		fmt.Printf("✅ Created user: %s (%s)\n", seed.Username, seed.Email)
		fmt.Printf("   Display Name: %s\n", seed.Display)
		fmt.Printf("   Roles: %v\n", seed.Roles)
		fmt.Printf("   Password: %s\n", seed.Password)
		fmt.Println()
	}

	fmt.Println("==========================================")
	fmt.Println("🎉 User seeding completed!")

	// Write credentials to file
	outputFile := "test_users_credentials.txt"
	writeCredentialsToFile(seedUsers, outputFile)
	fmt.Printf("\n📄 Credentials saved to: %s\n", outputFile)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type SeedUser struct {
	Username string
	Email    string
	Password string
	Roles    []string
	Display  string
}

func writeCredentialsToFile(users []SeedUser, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create credentials file: %v", err)
		return
	}
	defer file.Close()

	// Write header
	file.WriteString("=" + strings.Repeat("=", 110) + "\n")
	file.WriteString("  ALGORITHMIA TEST USER CREDENTIALS\n")
	file.WriteString("  Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	file.WriteString("=" + strings.Repeat("=", 110) + "\n\n")

	// Write table header
	file.WriteString(fmt.Sprintf("%-15s %-30s %-20s %-40s\n", "USERNAME", "EMAIL", "PASSWORD", "ROLES"))
	file.WriteString(strings.Repeat("-", 110) + "\n")

	// Write user data
	for _, user := range users {
		rolesStr := strings.Join(user.Roles, ", ")
		file.WriteString(fmt.Sprintf("%-15s %-30s %-20s %-40s\n",
			user.Username,
			user.Email,
			user.Password,
			rolesStr))
	}

	file.WriteString("\n" + strings.Repeat("=", 110) + "\n")
	file.WriteString("LOGIN GUIDE:\n")
	file.WriteString(strings.Repeat("=", 110) + "\n\n")

	// Create a map for easy lookup
	userMap := make(map[string]SeedUser)
	for _, user := range users {
		userMap[user.Username] = user
	}

	// Group by role with actual passwords
	file.WriteString("SUPER ADMIN:\n")
	if user, ok := userMap["admin"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString("PROBLEM SETTERS:\n")
	if user, ok := userMap["setter1"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n", user.Username, user.Password))
	}
	if user, ok := userMap["setter2"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString("PROBLEM REVIEWERS:\n")
	if user, ok := userMap["reviewer1"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n", user.Username, user.Password))
	}
	if user, ok := userMap["reviewer2"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString("PROBLEM TESTERS:\n")
	if user, ok := userMap["tester1"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n", user.Username, user.Password))
	}
	if user, ok := userMap["tester2"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n", user.Username, user.Password))
	}
	if user, ok := userMap["tester3"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString("CONTEST MANAGER:\n")
	if user, ok := userMap["contest_mgr"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString("MULTI-ROLE USER (setter + reviewer + tester):\n")
	if user, ok := userMap["multi_role"]; ok {
		file.WriteString(fmt.Sprintf("  Username: %s | Password: %s\n\n", user.Username, user.Password))
	}

	file.WriteString(strings.Repeat("=", 110) + "\n")
	file.WriteString("\nNOTE: Passwords are randomly generated for security. Please use the passwords shown above.\n")
	file.WriteString(strings.Repeat("=", 110) + "\n")
}
