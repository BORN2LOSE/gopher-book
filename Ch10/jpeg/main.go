// Команда jpeg читает изображение PNG из стандартного ввода
// и записывает его как изображение JPEG на стандартный вывод.
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Регистрация PNG-декодера
	"io"
	"os"
)

func main() {
	if err := convert(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

// convert to JPEG
func convert(ior io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(ior)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Входной формат =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
