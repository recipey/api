package postgres

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres driver -- never referenced but needed for database/sql to build

	"github.com/recipey/api"
)

type userRepository struct {
	db *sql.DB
}

// Application logic should take care of attributes validation. This method should just
// store the user as is.
func (ur *userRepository) Store(user *api.User) error {
	var err error
	if user.ID == "" {
		_, err = ur.db.Exec(`INSERT INTO users (username, email) VALUES ($1, $2)`,
			user.Username, user.Email)
	} else {
		_, err = ur.db.Exec(`UPDATE users SET username = $1, email = $2 WHERE id = $3`,
			user.Username, user.Email, user.ID)
	}
	return err
}

func (ur *userRepository) Find(id string) (*api.User, error) {
	var user *api.User
	row := ur.db.QueryRow(`SELECT * FROM users WHERE id = $1`, id)
	row.Scan(&user.ID, &user.Username, &user.Email)
	return user, nil
}

func (ur *userRepository) FindAll() ([]*api.User, error) {
	var users []*api.User
	rows, err := ur.db.Query(`SELECT * FROM users`)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user *api.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

// NewUserRepository returns postgres repository
func NewUserRepository(db *sql.DB) (api.UserRepository, error) {
	ur := &userRepository{db: db}

	return ur, nil
}
