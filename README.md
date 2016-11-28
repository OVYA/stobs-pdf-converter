# Documentation Stobs

stobs est un outil qui vous permet de réorganiser vos fichiers pdf scannés en recto/verso
Il peut etre utilisé sous linux et windows 
A partir d'un fichier pdf séléctionné, vous indiquez la première page du fichier en verso 
et il s'occupe de réorganiser le fichier pour une eventuelle impression

## Creation du package stobs.deb

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
  
### Remplir le fichier control

    Package: stobs
    Version: 1.0
    Section: base
    Priority: optional
    Architecture: all
    Depends: bash
    Maintainer: Dubosts Renaud
    Description: "Organisez vos pdf en recto/verso"

### Creer le binaire stobs

* go build main.go
* cp stobs usr/bin/.

### Construire le package 

* remonter dans le repertoire parent du dossier stobs 
  cd ..
  
* su dpkg-deb --build stobs

### installation 

* su dpkg -i stobs.deb

## utilisation

* selectionner un fichier 
* visualiser le dossier pour verification (optionnel)
* le nombre de pages total est indiqué
* selectionner la première page du fichier qui doit etre verso
* visualiser le nouveau fichier (en recto/verso)
* enregistrer le fichier  le fichier sera enregistré avec l'extension _recto_verso_ok.pdf dans le dossier de l'original
