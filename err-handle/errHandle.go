package errhandle

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
