package base

import "os"

func init() {
	err := os.RemoveAll(BASE_TMP_DIR)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(BASE_TMP_DIR, 0755)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(BASE_TMP_DIR)
	if err != nil {
		panic(err)
	}

}
