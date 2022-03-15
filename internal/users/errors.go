package users

type WrongUsernameOrPassword struct{}

func (w *WrongUsernameOrPassword) Error() string {
	return "wrong username or password"
}
