package main

import (
	"bufio"
	"fmt"
	"github.com/notJoon/pcg"
	"os"
)

func main() {
	rng := pcg.NewPCG32()

	file, err := os.Create("random_numbers.txt")
	if err != nil {
		fmt.Println("failed to generate numbers:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// write random numbers to the file
	for i := 0; i < 1200; i++ {
		number := rng.Uintn32(256) // 0 ~ 255
		_, err := writer.WriteString(fmt.Sprintf("%d\n", number))
		if err != nil {
			fmt.Println("파일 쓰기 실패:", err)
			return
		}
	}

	writer.Flush()

	fmt.Println("FINISHED!")
}
