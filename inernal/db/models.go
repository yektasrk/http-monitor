package db

type User struct {
	ID       uint
	Username string `gorm:"unique"`
	Password string
}

func (Client Client) SaveUser(user User) error {
	result := Client.db.Create(&user)
	return result.Error
}

func (Client Client) GetUser(username string) (User, error) {
	user := User{}
	result := Client.db.Where("username = ?", username).First(&user)
	return user, result.Error
}

