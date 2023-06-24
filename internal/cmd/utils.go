package cmd

import "os"

// ExposeToJail will call mknod on the block device to ensure
// visibility of the device
func ExposeToJail(dst string, uid, gid int) error {

	if err := os.Chmod(dst, 0600); err != nil {
		return err
	}

	if err := os.Chown(dst, uid, gid); err != nil {
		return err
	}

	return nil
}
