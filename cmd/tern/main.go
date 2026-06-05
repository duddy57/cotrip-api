package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmdArgs := []string{
		"migrate",
		"--migrations", "./internal/pgstore/migrations",
		"--config", "./internal/pgstore/migrations/tern.conf",
	}

	var output []byte
	var err error
	maxAttempts := 5

	for i := 1; i <= maxAttempts; i++ {
		cmd := exec.Command("tern", cmdArgs...)
		output, err = cmd.CombinedOutput()
		if err == nil {
			break
		}

		fmt.Errorf("Tentativa %d de %d falhou. Banco pode estar iniciando. Aguardando...", i, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		fmt.Println("Command execution failed definitivamente: ", err)
		fmt.Println("Output: ", string(output))
		os.Exit(1) // Fecha de forma limpa sem panic gordo na tela
	}

	fmt.Println("Migrações executadas com sucesso!")
}
