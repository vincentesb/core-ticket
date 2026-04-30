package error

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
