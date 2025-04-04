package user

import (
	"database/sql"
	"fmt"
	"forum/server"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	SEPARATOR = "|"
)

type User struct {
	UUID              string    `json:"user_uuid"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
	Role              string    `json:"role"`
	ProfilePicture    string    `json:"profile_picture"`
	Gender            string    `json:"gender"`
	Age               int       `json:"age"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewUser(uuid, username, email, encryptedPassword string, createdAt time.Time, role string, profile_picture string, gender string, last_name string, first_name string, age int) User {
	newUser := User{uuid, username, email, encryptedPassword, createdAt, role, profile_picture, gender, age, first_name, last_name}
	return newUser
}

// Trouver un utilisateur par son email et le renvoyer ( pour login )

func FetchUserByEmail(email string) (User, error) {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if re.FindAllString(email, -1) != nil {
		return User{}, fmt.Errorf("injection detected")
	}

	fetchUserQuery := `SELECT * FROM users WHERE email= ?`
	params := []interface{}{email}

	rows, err := server.RunQuery(fetchUserQuery, params...)
	if err != nil {
		return User{}, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même email. C'est normal ça ?")
	} else if len(rows) == 0 {
		return User{}, nil
	}

	newUser := User{}
	result := rows[0]

	// Utiliser des assertions de type avec vérification de valeur nulle
	if v, ok := result["user_uuid"]; ok && v != nil {
		newUser.UUID = v.(string)
	}
	if v, ok := result["username"]; ok && v != nil {
		newUser.Username = v.(string)
	}
	if v, ok := result["password"]; ok && v != nil {
		newUser.EncryptedPassword = v.(string)
	}
	if v, ok := result["profile_picture"]; ok && v != nil {
		newUser.ProfilePicture = v.(string)
	}
	if v, ok := result["email"]; ok && v != nil {
		newUser.Email = v.(string)
	}
	if v, ok := result["role"]; ok && v != nil {
		newUser.Role = v.(string)
	}
	if v, ok := result["created_at"]; ok && v != nil {
		// Assurez-vous que c'est bien un type time.Time
		if parsedTime, ok := v.(time.Time); ok {
			newUser.CreatedAt = parsedTime
		} else {
			fmt.Fprintln(os.Stderr, "Erreur : created_at n'est pas un time.Time")
		}
	}

	return newUser, nil
}

// Savoir si un utilisateur existe par son nom d'utilisateur ( pour register )

func IsUsernameTaken(username string) (bool, error) {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if re.FindAllString(username, -1) != nil {
		return false, fmt.Errorf("injection detected")
	}

	fetchUserQuery := `SELECT * FROM users WHERE username= ?`
	params := []interface{}{username}

	rows, err := server.RunQuery(fetchUserQuery, params...)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	return len(rows) >= 1, nil
}

// Trouver l'image de profil utilisateur avec ID
func FetchPPByID(id string) (string, error) {
	fetchUserQuery := `SELECT profile_picture FROM users WHERE user_uuid= ?`
	params := []interface{}{id}

	rows, err := server.RunQuery(fetchUserQuery, params...)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même ID. C'est normal ça ?")
	} else if len(rows) == 0 {
		return "", nil
	}

	usrFound := User{}
	result := rows[0]

	// Utiliser une assertion de type avec vérification de valeur nulle
	if v, ok := result["profile_picture"]; ok && v != nil {
		usrFound.ProfilePicture = v.(string)
	}

	return usrFound.ProfilePicture, nil
}

func (u *User) ToMap() map[string]interface{} {
	usrMap := make(map[string]interface{}, 0)

	fmt.Println("user: ", u.EncryptedPassword)

	usrMap["user_uuid"] = u.UUID
	usrMap["username"] = u.Username
	usrMap["email"] = u.Email
	usrMap["password"] = u.EncryptedPassword
	usrMap["created_at"] = u.CreatedAt.Format("2006-01-02")
	usrMap["role"] = u.Role
	usrMap["profile_picture"] = u.ProfilePicture
	usrMap["last_name"] = u.LastName
	usrMap["first_name"] = u.FirstName
	usrMap["gender"] = u.Gender
	usrMap["age"] = u.Age

	return usrMap
}

func (u *User) ToCookieValue() string {
	return u.UUID + SEPARATOR +
		u.Username + SEPARATOR +
		u.Email + SEPARATOR +
		u.Role
}

// Enregistrer un user complet ( Register )

func RegisterUser(params map[string]interface{}) error {

	profile_picture, _ := params["profile_picture"].(string)

	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if params["password"] != "" {
		for key, value := range params {
			if (key == "username" || key == "email" || key == "password") && re.FindAllString(value.(string), -1) != nil {
				return fmt.Errorf("injection detected")
			}
		}
	}

	if profile_picture == "" {
		profile_picture = RandomProfilPicture()
	}

	registerUserQuery := `INSERT INTO users (user_uuid, username, first_name, last_name, gender, age, email, password, created_at, profile_picture, role )  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var err error

	_, err = server.RunQuery(registerUserQuery, params["user_uuid"], params["username"], params["first_name"], params["last_name"], params["gender"], params["age"], params["email"], params["password"], params["created_at"], profile_picture, params["role"])

	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion de l'utilisateur : %v", err)
	}

	return nil
}

// Mettre à jour les valeurs d'un utilisateur ( Update )
func (u *User) UpdateUser(params map[string]interface{}) error {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if params["password"] != "" {
		for key, value := range params {
			if (key == "username" || key == "email" || key == "password") && re.FindAllString(value.(string), -1) != nil {
				return fmt.Errorf("injection detected")
			}
		}
	}

	updateUserQuery := `UPDATE users SET username = ?, email = ?, password = ?, profile_picture = ? WHERE user_uuid = ?`
	_, err := server.RunQuery(updateUserQuery, params["username"], params["email"], params["password"], params["profile_picture"], params["user_uuid"])

	if err != nil {
		return err
	}

	return nil
}

// FetchAllComments récupère tous les commentaires de la base de données
func FetchAllUsers(db *sql.DB) ([]User, error) {
	results, err := server.RunQuery("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []User

	// ce qu'on veut renvoyer
	for _, row := range results {
		user := User{}

		if userUUID, ok := row["user_uuid"].(string); ok {
			user.UUID = userUUID
		}
		if createdAt, ok := row["created_at"].(time.Time); ok {
			user.CreatedAt = createdAt
		}
		if username, ok := row["username"].(string); ok {
			user.Username = username
		}
		if profilePicture, ok := row["profile_picture"].(string); ok {
			user.ProfilePicture = profilePicture
		}

		if role, ok := row["role"].(string); ok {
			user.Role = role
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUserRole(user_UUID, action string) error {
	getUserRole := `SELECT role FROM users WHERE user_uuid = ?`
	results, err := server.RunQuery(getUserRole, user_UUID)
	if err != nil {
		return fmt.Errorf("error fetching user role: %v", err)
	}

	if len(results) == 0 {
		return fmt.Errorf("no user found with the specified UUID")
	}

	// Récupère le rôle actuel de l'utilisateur
	roleValue, ok := results[0]["role"].(string)
	if !ok {
		return fmt.Errorf("invalid role type for user")
	}

	// Détermine le nouveau rôle en fonction de l'action et du rôle actuel
	var newRole string
	if action == "promote" {
		if roleValue == "user" {
			newRole = "modo"
		} else if roleValue == "modo" {
			newRole = "admin"
		}
	} else if action == "demote" && roleValue == "modo" {
		newRole = "user"
	} else if action == "demote" && roleValue == "admin" {
		newRole = "modo"
	}

	if newRole == "" {
		return nil
	}

	updateUserRole := `
    UPDATE users 
    SET role = ?
    WHERE user_uuid = ?`

	_, err = server.RunQuery(updateUserRole, newRole, user_UUID)
	if err != nil {
		return fmt.Errorf("failed to update role: %v", err)
	}

	return nil
}
