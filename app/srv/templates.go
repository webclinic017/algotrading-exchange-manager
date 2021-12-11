package srv

import (
	"io"
	"os"
)

func CheckFiles() {

	//FileCopyIfMissing("app/templates/ENV_Settings.env", "app/config/ENV_Settings.env")
	FileCopyIfMissing("app/templates/ENV_accesstoken.env", "app/config/ENV_accesstoken.env")
	FileCopyIfMissing("app/templates/trackSymbols.txt", "app/config/trackSymbols.txt")
}

// Copy the src file to dst. Skipped if file exists
func FileCopyIfMissing(src, dst string) error {

	_, err := os.Open(dst)
	if err != nil {
		// File does not exist, copy it

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return out.Close()
	}
	return nil
}
