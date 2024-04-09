package controller

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

// Test connection to NFS server
func TestNFSConnection() error {
	ip := "192.168.1.11"
	// Check NFS server
	err := pingNFS(ip)
	if err != nil {
		return fmt.Errorf("failed to connect to NFS server (Probably unreachable): %w", err)
	}

	// Verificar si localPath existe, y si no, crearlo
	localPath := "./nfs"
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		if err := os.MkdirAll(localPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", localPath, err)
		}
	}

	// Cambiar el propietario y los permisos del directorio de montaje
	if err := changeDirectoryPermissions(localPath); err != nil {
		return fmt.Errorf("failed to change directory permissions: %w", err)
	}

	// Mount the NFS server as a local disk
	nfsPath := "/var/nfs/testing"
	mountCmd := exec.Command("sudo", "mount", "-t", "nfs", ip+":"+nfsPath, localPath)
	err = mountCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to mount NFS server: %w", err)
	}

	testFilePath := localPath + "/test.txt"
	testFile, err := os.Create(testFilePath)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer testFile.Close()

	fmt.Println("Successfully connected to NFS server and created test file in " + testFilePath)
	return nil
}

// Ping the NFS server to check if it is reachable
func pingNFS(ip string) error {
	pingCmd := exec.Command("ping", "-c", "4", ip)
	err := pingCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to ping NFS server: %w", err)
	}

	return nil
}

// Cambiar el propietario y los permisos del directorio
func changeDirectoryPermissions(path string) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	chownCmd := exec.Command("chown", "-R", currentUser.Username+":", path)
	err = chownCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change directory owner: %w", err)
	}

	chmodCmd := exec.Command("chmod", "-R", "755", path)
	err = chmodCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change directory permissions: %w", err)
	}

	return nil
}
