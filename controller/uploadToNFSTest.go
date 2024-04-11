package controller

import (
	"fmt"
	"os"
)

// Test connection to NFS server
func TestNFSConnection() error {

	localPath := "/mnt/nfs"

	testFilePath := localPath + "/test.txt"
	testFile, err := os.Create(testFilePath)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer testFile.Close()

	fmt.Println("Successfully connected to NFS server and created test file in " + testFilePath)
	return nil
}