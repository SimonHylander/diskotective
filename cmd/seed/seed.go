package seed

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand/v2"
	"os"
	"path/filepath"
)

func Execute() {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, "seed")

	// check if seed dir exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	for i := 0; i < 5; i++ {
		createDirectory(path)
		createFile(path)
	}
}

func createDirectory(path string) {
	name := uuid.New().String()
	newDir := filepath.Join(path, name)
	err := os.Mkdir(newDir, 0755)

	if err != nil {
		panic(err)
	}

	createFile(newDir)
}

func createFile(path string) {
	name := uuid.New().String()
	file, err := os.Create(filepath.Join(path, name))

	if err != nil {
		panic(err)
	}

	// Randomize content
	for i := 0; i < rand.IntN(50); i++ {
		_, err = file.WriteString(fmt.Sprintf("%s\n", uuid.New().String()))
	}

	err = file.Close()
	if err != nil {
		return
	}
}
