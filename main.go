package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Check command line args.
	if len(os.Args) <= 1 {
		fmt.Printf("Example usage: %s xspeak.lib\n", os.Args[0])
		os.Exit(1)
	}

	// Open the file.
	f, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Figure out which type of lib file this is.
	var entries []*Entry
	var xorEnabled bool
	if strings.HasPrefix(os.Args[1], "x") {
		log.Println("Assuming named lib from filename.")
		entries, xorEnabled, err = getNamedEntries(f)
	} else {
		log.Println("Assuming SoundID-based lib from filename.")
		entries, xorEnabled, err = getSoundIDEntries(f)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt if the file data is XOR crypted.
	if xorEnabled {
		log.Println("Decrypting file entries.")
		for _, entry := range entries {
			XORCrypt(entry)
		}
	}

	// Output files.
	outputDir := fmt.Sprintf("./extract_%s/", os.Args[1])
	err = os.MkdirAll(outputDir, 0666)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		log.Println("Writing file", outputDir+entry.Name)
		of, err := os.Create(outputDir + entry.Name)
		if err != nil {
			log.Fatal(err)
		}

		n, err := of.Write(entry.Data)
		if err != nil {
			log.Fatal(err)
		}
		if n != len(entry.Data) {
			log.Fatal("Failed to write whole file.")
		}
	}

}
