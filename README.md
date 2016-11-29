# Documentation

__stobs__ is a software which gives you the possibility to reorganized your scanned pdf files.
It can be used on Linux, Windows and MacOsX platform.
From a selected pdf file you choose the first back page and stobs reorganized all the pages for an eventually print


## stobs.deb package creation
* mkdir package
* mkdir stobs 
* cd package/stobs
* mkdir DEBIAN
* touch DEBIAN/control
* touch DEBIAN/postinst
* touch DEBIAN/postrm
* touch DEBIAN/preinst
* touch DEBIAN/prerm
* chmod 755 DEBIAN/post*
* chmod 755 DEBIAN/pre*
* mkdir usr
* mkdir usr/bin
  
### Fill the __Control__ file

    Package: stobs
    Version: 1.0
    Section: base
    Priority: optional
    Architecture: all
    Depends: bash,evince, pdftk
    Maintainer: Dubosts Renaud
    Description: "Organisez vos pdf en recto/verso"

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
