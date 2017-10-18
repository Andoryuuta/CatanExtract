package main

type Entry struct {
	Name string
	Size uint32
	Data []byte
}

func XORCrypt(entry *Entry) {
	xkey := byte(0x47)
	ykey := byte(0x7E)
	for i := uint32(0); i < entry.Size; i++ {
		outputByte := xkey ^ (entry.Data[i] - 1)
		entry.Data[i] = outputByte
		xkey += ykey
		ykey += 0x21
	}
}
