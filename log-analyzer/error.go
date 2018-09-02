package main

// if err1 is not nil, return err1. else, return err2.
func setErr(err1, err2 error) error {
	if err1 != nil {
		return err1
	}

	return err2
}
