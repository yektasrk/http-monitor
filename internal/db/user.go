package db

type User struct {
	ID       uint
	Username string `gorm:"unique"`
	Password string
}

func (client Client) SaveUser(user User) error {
	result := client.db.Create(&user)
	return result.Error
}

func (client Client) GetUser(username string) (User, error) {
	user := User{}
	result := client.db.Where("username = ?", username).First(&user)
	return user, result.Error
}
