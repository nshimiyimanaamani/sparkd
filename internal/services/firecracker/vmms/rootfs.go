// rootfs file is used to generate root filesystem for the VM
// using init binary process and supplied tar file form podman supplied by user.
package vmms

import (
	"fmt"
	"os"

	"github.com/quarksgroup/sparkd/internal/cmd"
)

// generateRFs generates root filesystem for the VM according to the below steps:
// 1. create a directory for the rootfs
// 2. copy the init binary to the rootfs
// 3. copy the init base tar file to the rootfs
// 4. extract the init base tar file
// 5. copy the podman supplied tar file to the rootfs
// 6. extract the podman supplied tar file
// 7. delete the init base tar file
// 8. delete the podman supplied tar file
// 9. return the rootfs path or name
func (o *Config) generateRFs(dir, name string) (string, error) {

	fsName := fmt.Sprintf("%d-%s.ext4", o.vmIndex, name)

	// for creating the rootfs directory with 1024MB size
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("dd if=/dev/zero of=%s bs=1 count=0 seek=1G", fsName)); err != nil {
		return "", fmt.Errorf("failed to create rootfs file: %v", err)
	}

	//for making the rootfs file as ext4 file system
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("mkfs.ext4 %s", fsName)); err != nil {
		return "", fmt.Errorf("failed to create ext4 file system: %v", err)
	}

	//creating a temporary directory for mounting the rootfs file
	tmpDir, err := os.MkdirTemp(dir, fmt.Sprintf("%d-%s", o.vmIndex, name))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// unmout the created tmp dir from rootfs file
	defer cmd.RunSudo(fmt.Sprintf("umount %s", tmpDir))

	//for mounting the created rootfs file to tmp directory
	if _, err := cmd.RunSudo(fmt.Sprintf("mount -o loop %s %s", fsName, tmpDir)); err != nil {
		return "", fmt.Errorf("failed to mount rootfs file: %v", err)
	}

	imageTar := fmt.Sprintf("%d-%s.tar", o.vmIndex, name)
	imageName := fmt.Sprintf("%d-%s", o.vmIndex, name)

	cmd.RunNoneSudo(fmt.Sprintf("podman rm -f %s > /dev/null | true", imageName))

	// for exporting the podman tar file from supplied podman image
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman create --name %s %s", imageName, o.providedImage)); err != nil {
		return "", fmt.Errorf("podman failed to create tar file: %v", err)
	}
	defer cmd.RunNoneSudo(fmt.Sprintf("podman rm -f %s ", imageName))

	// start temp imageName container
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman start %s", imageName)); err != nil {
		return "", fmt.Errorf("podman failed to start container: %v", err)
	}

	// exectute update-get update in the container
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman exec -it %s apt update", imageName)); err != nil {
		return "", fmt.Errorf("podman failed to update container: %v", err)
	}

	// Run the install net-tools for ifconfig inside the container
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman exec -it %s apt install net-tools", imageName)); err != nil {
		return "", fmt.Errorf("podman failed to run setup ifconfig in container: %v", err)
	}

	// // Run the install iproute2 command for ip command inside the container
	// if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman exec -it %s apt install iproute2", imageName)); err != nil {
	// 	return "", fmt.Errorf("podman failed to run setup ip cmd  in container: %v", err)
	// }

	// // exectute update-get update in the container
	// if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman exec -it %s apt update", imageName)); err != nil {
	// 	return "", fmt.Errorf("podman failed to update container2: %v", err)
	// }

	// if _, err := cmd.RunSudo(fmt.Sprintf("podman exec -it %s ip tuntap add dev %s mode tap", imageName, o.tap)); err != nil {
	// 	return "", fmt.Errorf("podman failed to run setup tap device: %v", err)
	// }

	fmt.Println("Networks", o.tapGateWay, o.tapMask)
	// if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman exec -it %s ifconfig %s %s netmask %s", imageName, o.tap, o.tapGateWay, o.tapMask)); err != nil {
	// 	return "", fmt.Errorf("podman failed to run setup tap mask address: %v", err)
	// }

	//
	cmd.RunNoneSudo(fmt.Sprintf("podman stop %s ", imageName))

	// for exporting the podman tar file from supplied podman image
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("podman export %s -o %s%s", imageName, dir, imageTar)); err != nil {
		return "", fmt.Errorf("podman failed to export tar file: %v", err)
	}

	// for extracting the podman supplied tar file to the rootfs directory
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("tar -xvf %s%s -C %s", dir, imageTar, tmpDir)); err != nil {
		return "", fmt.Errorf("failed to extract podman supplied tar file: %v", err)
	}

	// include our init process into ext4 file system exported from podman
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("cp -r %sinit %s", dir, tmpDir)); err != nil {
		return "", fmt.Errorf("failed to cp init to tmp dir: %v", err)
	}

	if _, err := cmd.RunNoneSudo(fmt.Sprintf("cp -r %srun.json %s", dir, tmpDir)); err != nil {
		return "", fmt.Errorf("failed to cp run.json to tmp dir: %v", err)
	}

	//remove those created ext and tar files
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("rm -f %s", imageTar)); err != nil {
		return "", fmt.Errorf("failed to remove ext and tar files: %v", err)
	}

	return fsName, nil
}
