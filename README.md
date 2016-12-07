# Documentation

__stobs__ is a software which gives you the possibility to reorganized your scanned pdf files.
From a selected pdf file you choose the first back page and stobs reorganized all the pages for an eventually print

### MakeFile description
* make all : launch the make install
* make install : install glide and build the go program
* make build : generate the project binary
* make glide-install : install all project dependancies
* make glide-update : update all project dependancies
* make vendor : install all dependancies
* make deb : create the project .deb package

### Fill the __Control__ file

    Package: stobs
    Version: 1.0
    Section: base
    Priority: optional
    Architecture: all
    Depends: bash, pdftk
    Maintainer: Ovya
    Description: "reorganize your pdf in two-sided format"

### Generate stobs binary

* go build main.go
* cp stobs usr/bin/.


### build the package 

* cd ..
* su dpkg-deb --build stobs

### installation 

* su dpkg -i stobs.deb

## use

* select a pdf file 
* display the file for checking (optional)
* select the first back page  _the total number of pages is indicate_
* display the new two-sided file
* save the file __the file will be saved with extension _recto_verso_ok.pdf in original folder file__


# Development environment

## installation of go-gtk lib

* go get github.com/mattn/go-gtk/gtk
__installation des dependances__
* apt-get install libgtk2.0-dev libglib2.0-dev libgtksourceview2.0-dev
you can follow instruction on the github project page at https://github.com/mattn/go-gtk
