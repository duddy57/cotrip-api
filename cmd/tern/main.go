package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/pgstore/migrations",
		"--config",
		"./internal/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command execution failed: ", err)
		fmt.Println("Output: ", string(output))
		panic(err)
	}

	fmt.Println("Command executed successfully ", string(output))
}
