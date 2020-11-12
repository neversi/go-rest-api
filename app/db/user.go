package db 

import 	(
	"github.com/neversi/go-rest-api/app/objects"
	"golang.org/x/crypto/bcrypt"
)



// UserRep is needed to 
type UserRep struct {
	dbUser *DB
}


// Create func creates user in user repository
func (user *UserRep) Create(currentUser *objects.User) (*objects.User, error) {
	encrypted, err := encryptString(currentUser.Password)
	if err != nil {
		return nil, nil
	}

	if currentUser.Validation() != nil {
		return nil, err
	}

	if err = user.dbUser.db.QueryRow(
		"INSERT INTO users (login, password, firstName, surName, email) values ($1, $2, $3, $4, $5) returning id", 
		currentUser.Login,
		encrypted,
		currentUser.FirstName,
		currentUser.SurName,
		currentUser.Email,
		).Scan(&currentUser.ID);
		err != nil {
			println("Haha")
			return nil, err
		}
	

		return currentUser, nil

}

// UserByLogin here
func (user  *UserRep) UserByLogin(login string) (*objects.User, error) {
	u := &objects.User{}

	if err := user.dbUser.db.QueryRow(
		"SELECT id, login, password FROM users WHERE login = $1",
		 login).Scan(
			 &u.ID,
			 &u.Login, 
			 &u.Password,
	); err != nil {
		return nil, err
	}
	
	return u, nil
		
}

func encryptString(str string) (string, error) {
	hBytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hBytes), nil
}

//CompareHash ..
func compareHash(hashed string, raw string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashed), 
		[]byte(raw),
		); err != nil {
			return false, err
		}
	
	return true, nil
}

