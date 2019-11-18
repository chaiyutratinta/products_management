package utils

//Checker for check error and panic
func Checker(err error) {
	if err != nil {
		panic(err)
	}
}
