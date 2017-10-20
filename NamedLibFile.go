package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

const (
	namedHeaderSize = 0x8
	namedEntrySize  = 0x3C
)

func readNamedHeaderInfo(f *os.File) (bool, uint32, error) {
	var unk, xorCryptEnabled uint16
	var entriesCount uint32

	err := binary.Read(f, binary.LittleEndian, &unk)
	if err != nil {
		return false, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &xorCryptEnabled)
	if err != nil {
		return false, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &entriesCount)
	if err != nil {
		return false, 0, err
	}

	return xorCryptEnabled == 1, entriesCount, nil
}

func readNamedEntryInfo(f *os.File) (*Entry, error) {
	var nameBytes [52]byte
	var fileOffset, fileSize uint32

	err := binary.Read(f, binary.LittleEndian, &nameBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Read(f, binary.LittleEndian, &fileOffset)
	if err != nil {
		return nil, err
	}

	err = binary.Read(f, binary.LittleEndian, &fileSize)
	if err != nil {
		return nil, err
	}

	// Name is null-terminated, so...
	nameLength := 0

	// C-Style strlen thingy-ma-bobber.
	for ; nameBytes[nameLength] != 0; nameLength++ {
	}

	return &Entry{Name: string(nameBytes[:nameLength]), Offset: fileOffset, Size: fileSize}, nil
}

func getNamedEntries(f *os.File) ([]*Entry, bool, error) {
	// Read the header info.
	log.Println("Reading named file header.")
	xorEnabled, entriesCount, err := readNamedHeaderInfo(f)
	if err != nil {
		return nil, false, err
	}

	log.Printf("%s contains %d files. XOR enabled:%v\n", os.Args[1], entriesCount, xorEnabled)

	// Read in all the entries basic info (Name and size).
	log.Println("Reading named file entries information.")
	var entries []*Entry
	for i := uint32(0); i < entriesCount; i++ {
		entry, err := readNamedEntryInfo(f)
		if err != nil {
			return nil, false, err
		}

		entries = append(entries, entry)
	}

	// Calculate start of data section.
	startOfDataSection := namedHeaderSize + (namedEntrySize * entriesCount)

	// Read entry data.
	log.Println("Reading named file entries data.")
	for _, entry := range entries {
		// Seek to the entry data.
		log.Printf("at f.Seek - entry.Offset 0x%X\n", entry.Offset)
		_, err = f.Seek(int64(startOfDataSection+entry.Offset), io.SeekStart)
		if err != nil {
			return nil, false, err
		}

		// Make a buffer for the data and read into it.
		entry.Data = make([]byte, entry.Size)

		err = binary.Read(f, binary.LittleEndian, &entry.Data)
		if err != nil {
			return nil, false, err
		}
	}

	return entries, xorEnabled, nil
}
