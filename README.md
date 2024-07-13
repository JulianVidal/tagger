# Tagger

## Description

Tag everything in your computer, notes, websites, images and documents with attribute-value tags. This will allow for easy and quick searching for anything you need. This will be a CLIutility that you can pass in paths and tag it into the system. Everything will be stored locally with easy migration between computers.


Inspiration:
https://www.youtube.com/watch?v=wTQeMkYRMcw

Sources:
https://hci.ucsd.edu/hollan/hotos.pdf

The program will look for files in $XDG_HOME_DATA/tagger, defaults to ~/.local/share/tagger. All the files in this folder will be automatically added, you can symlink folders and files into this directory and have it read. This allows you to move the files wherever you want, fix the symlinks and the program should mantain the same tags.

The program only read 1 directory, won't read symlinked directories within a symlinked directory. This is to avoid loops and it is not supported by golang.

## TODO
 - Create cli tool
     - Returns the files that are in the engine, piped to fuzzy search?

