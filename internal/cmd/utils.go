package cmd

import (
	"fmt"
	"os"
	"syscall"
)

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

// mount mounts a filesystem to a target
func Mount(source, target, filesystemtype string, flags uintptr) error {

	if _, err := os.Stat(target); os.IsNotExist(err) {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return fmt.Errorf("error creating target folder: %s %s", target, err)
		}
	}

	fmt.Println("Mounting:", source, target, filesystemtype, flags)
	err := syscall.Mount(source, target, filesystemtype, flags, "")
	if err != nil {
		return fmt.Errorf("error mounting %s to %s, error: %v", source, target, err)
	}
	return nil
}
