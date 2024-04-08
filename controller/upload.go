package controller

import (
	"fmt"
	"os"
	"os/exec"
)

// Test connection to NFS server
func TestNFSConnection() error {
	ip := "192.168.0.1"
	// Check NFS server
	err := PingNFS(ip)
	if err != nil {
		return fmt.Errorf("failed to connect to NFS server (Probably unreachable): %w", err)
	}

	// Mount the NFS server as a local disk
	localPath := "./nfs"
	nfsPath := "/path/to/nfs"
	mountCmd := exec.Command("mount", "-t", "nfs", ip+nfsPath, localPath)
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
func PingNFS(ip string) error {
	pingCmd := exec.Command("ping", "-c", "1", ip)
	err := pingCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to ping NFS server: %w", err)
	}

	return nil
}
