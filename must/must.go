package must

func String(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func ByteArray(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
