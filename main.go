package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Uso: %s <arquivo1> <arquivo2> ...\n", os.Args[0])
		return
	}

	outputDir := "output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Erro ao criar diretório de saída: %v\n", err)
		return
	}

	for _, filePath := range os.Args[1:] {
		err := processFile(filePath, outputDir)
		if err != nil {
			fmt.Printf("Erro ao processar arquivo %s: %v\n", filePath, err)
		}
	}
}

func processFile(filePath, outputDir string) error {

	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("erro ao calcular o hash: %w", err)
	}

	hashSum := hash.Sum(nil)
	hashHexadecimal := fmt.Sprintf("%x", hashSum)

	fileName := filepath.Base(filePath)
	labelPath := remove(strings.Split(fileName, "."), 1)

	outputFileDir := filepath.Join(outputDir, labelPath[0])
	err = os.MkdirAll(outputFileDir, 0755)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório de saída para %s: %w", labelPath[0], err)
	}

	outputFilePath := filepath.Join(outputFileDir, "hash.txt")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de saída: %w", err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(hashDateWithString(hashHexadecimal))
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de saída: %w", err)
	}

	fmt.Printf("SHA-256 hash do arquivo %s: %s\n", filePath, hashHexadecimal)
	return nil
}

func hashDateWithString(hashHexadecimal string) string {
	return fmt.Sprintf("%s : %s", hashHexadecimal, time.Now().Local().Format("02/01/2006 15:04:05"))
}

func remove(slice []string, s int) []string {
	slice = append(slice[:s], slice[s+1:]...)
	return slice
}
