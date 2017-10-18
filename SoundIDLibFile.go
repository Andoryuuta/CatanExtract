package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func readSoundIDHeaderInfo(f *os.File) (bool, uint32, uint32, uint32, error) {
	var unk, xorCryptEnabled uint16
	var minSoundID, maxSoundID uint32
	var entriesCount uint32

	err := binary.Read(f, binary.LittleEndian, &unk)
	if err != nil {
		return false, 0, 0, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &xorCryptEnabled)
	if err != nil {
		return false, 0, 0, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &minSoundID)
	if err != nil {
		return false, 0, 0, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &maxSoundID)
	if err != nil {
		return false, 0, 0, 0, err
	}

	err = binary.Read(f, binary.LittleEndian, &entriesCount)
	if err != nil {
		return false, 0, 0, 0, err
	}

	return xorCryptEnabled == 1, minSoundID, maxSoundID, entriesCount, nil
}

func readSoundIDEntryInfo(f *os.File) (*Entry, error) {
	var unk1, fileSize, unk2 uint32

	err := binary.Read(f, binary.LittleEndian, &unk1)
	if err != nil {
		return nil, err
	}

	err = binary.Read(f, binary.LittleEndian, &fileSize)
	if err != nil {
		return nil, err
	}

	err = binary.Read(f, binary.LittleEndian, &unk2)
	if err != nil {
		return nil, err
	}

	return &Entry{Size: fileSize}, nil
}

func getSoundIDEntries(f *os.File) ([]*Entry, bool, error) {
	// Read the header info.
	log.Println("Reading SID-based file header.")
	xorEnabled, minSID, _, entriesCount, err := readSoundIDHeaderInfo(f)
	if err != nil {
		return nil, false, err
	}

	log.Printf("%s contains %d files. XOR enabled:%v\n", os.Args[1], entriesCount, xorEnabled)

	// Read in all the entries basic info (size).
	log.Println("Reading SID-based file entries information.")
	var entries []*Entry
	for i := uint32(0); i < entriesCount; i++ {
		entry, err := readSoundIDEntryInfo(f)
		if err != nil {
			return nil, false, err
		}

		// Assume proper sound ID is the index from the minimum.
		entry.Name = fmt.Sprintf("%d.mp3", minSID+i)

		// Check if sound ID has any data.
		if entry.Size == 0 {
			log.Printf("Entry %s contains no data, skipping.\n", entry.Name)
			continue
		}

		entries = append(entries, entry)
	}

	// Read entry data.
	log.Println("Reading SID-based file entries data.")
	for _, entry := range entries {
		entry.Data = make([]byte, entry.Size)

		err = binary.Read(f, binary.LittleEndian, &entry.Data)
		if err != nil {
			return nil, false, err
		}
	}

	return entries, xorEnabled, nil
}
