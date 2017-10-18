# CatanExtract
An extractor for Catan's (https://archive.org/details/DeKolonistenVanCatanZeevaarders) .lib archive file format. 


### Installation from source

```
go get -u github.com/Andoryuuta/CatanExtract
```

### Download binary (Windows users)
You can simply download the latest release binary from [here](https://github.com/Andoryuuta/CatanExtract/releases).

### Usage Example
```
CatanExtract xspeak.lib
```
OR
```
CatanExtract 2speak.lib
```
The files will be extracted to a new folder in the current working directory.

### Notes
CatanExtract differentiates between two different (though similar) types of .lib files by looking at the file prefix. Named-based, and SoundID-based files. If a certain file doesn't work, try simply adding (or removing, it already exists) a 'x' to the front of the .lib file name.
