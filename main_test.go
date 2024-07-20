package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func deleteTestFile(fileName string) error {
	if err := os.RemoveAll(fileName); err != nil {
		return fmt.Errorf("could not remove temp file at:%s", fileName)
	} else {
		return nil
	}
}

func TestWhenKubeconfigExistsThenDoesKubeConfigExistShouldReturnTrue(t *testing.T) {
	kubeConfigDirPath := filepath.Dir(KubeConfigPath)
	var testConfig *os.File
	tempFileCreated := false
	// create ~/.kube folder if it does not exist
	if _, err := os.Stat(kubeConfigDirPath); os.IsNotExist(err) {
		err := os.Mkdir(kubeConfigDirPath, 0770)
		if err != nil {
			t.Fatalf("could not create temp kubeconfig dir: %s", err)
		}
	}

	// create config file if it does not exist
	if _, err := os.Stat(KubeConfigPath); os.IsNotExist(err) {
		// create a temp file at the expected default .kube/config location
		tempFileCreated = true
		testConfig, err = os.Create(KubeConfigPath)
		if err != nil {
			t.Fatalf("Error creating temporary kubeconfig file: %s", err)
		}
	} else {
		testConfig, err = os.Open(KubeConfigPath)
	}
	defer testConfig.Close()

	if tempFileCreated {
		defer func() {
			err := deleteTestFile(kubeConfigDirPath)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	if doesKubeConfigExist(testConfig.Name()) == false {
		t.Error("doesKubeConfigExist function returned false even though the config file existed")
	}
}

func TestWhenKubeconfigDoesNotExistThenDoesKubeConfigExistShouldReturnFalse(t *testing.T) {
	brokenKubeConfigPath := filepath.Join(KubeConfigPath, "nonsense")
	if doesKubeConfigExist(brokenKubeConfigPath) == true {
		t.Error("doesKubeConfigExist function returned true even though the config file did not exist")
	}
}
