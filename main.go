package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"io/ioutil"
	"os"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var nbPages int
var cat = false
var outFile string
var inFile string

func main() {

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("STOBS")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		defer os.Remove(inFile)
		defer os.Remove(outFile)
		gtk.MainQuit()
	})

	// creation du conteneur Principal
	vbox := gtk.NewVBox(false, 1)
	vpaned := gtk.NewVPaned()
	vPaned := gtk.NewVPaned()
	vPaned2 := gtk.NewVPaned()
	vPaned.SetSizeRequest(600, 550)
	vpaned.Add(vPaned)
	vpaned.Add(vPaned2)
	vbox.Add(vpaned)

	//
	//CREATION DE L'INTERFACE
	//
	label := gtk.NewLabel("organisez vos PDF en Recto/Verso")
	label.ModifyFontEasy("DejaVu Serif 14")

	//creation du frame fichier
	//contient tout les elements label, textBox ...
	frame := gtk.NewFrame("Fichier")
	framebox := gtk.NewVBox(false, 5)
	frame.Add(framebox)

	frame2 := gtk.NewFrame("")
	framebox2 := gtk.NewVBox(false, 5)
	frame2.Add(framebox2)

	//ajout des elements au conteneur principal
	vPaned.Pack1(label, false, false)
	vPaned.Pack2(frame, false, false)
	vPaned2.Pack1(frame2, false, false)

	box1 := gtk.NewHBox(false, 5)
	box2 := gtk.NewHBox(false, 5)
	box3 := gtk.NewHBox(false, 5)
	box4 := gtk.NewHBox(false, 5)

	// label info fichier
	lblPage := gtk.NewLabel("")

	// txtBox file
	fileName := gtk.NewEntry()
	fileName.SetText("")

	lblVerso := gtk.NewLabel("Indiquer la première page Verso: ")
	txtVerso := gtk.NewEntry()

	// bouton ouverture file dialog
	btOpenFile := gtk.NewButtonWithLabel("Ouvrir")
	btOpenFile.Clicked(func() {
		//-----------------------------------------------------
		// boite de dialog selection du fichier
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			btOpenFile.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*.pdf")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			if filechooserdialog.GetFilename() != "" {
				fileName.SetText(filechooserdialog.GetFilename())
				nbPages = getNumberOfPAges(filechooserdialog.GetFilename())
				lblPage.SetText("Le fichier séléctionné contient " + strconv.Itoa(nbPages) + " pages.")
			}

			//normalizeFileName(fileName.GetText())
			createTempFile()
			saveFile(inFile,fileName.GetText())
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})

	//btSave
	btSave := gtk.NewButtonWithLabel("Enregistrer")
	btSave.Clicked(func() {
		if txtVerso.GetText() != "" && !cat {
			catFile(fileName, txtVerso)
		}

		if cat == true {
			file := strings.Replace(fileName.GetText(), ".pdf", "_recto_verso_ok.pdf", 1)
			saveFile(file, outFile)
			dial := gtk.NewMessageDialog(window, 1, 1, 1, "La modification du PDF a été effectué avec succès ")
			dial.Response(func() {
				dial.Destroy()
				cat = false
			})

			dial.Run()
		}
	})

	//btn Visualiser le fichier
	btShowFile := gtk.NewButtonWithLabel("Visualiser le fichier")
	btShowFile.Clicked(func() {
		if txtVerso.GetText() != "" {
			catFile(fileName, txtVerso)
			displayfile(outFile)
		}else{
			displayfile(inFile)
		}

	})

	// Ajout des elements dans leurs box correspondante
	// 1 box correspond à une ligne de l'interface
	box1.Add(fileName)
	box1.Add(btOpenFile)
	box2.Add(lblPage)
	box3.Add(btShowFile)
	box3.Add(btSave)
	box4.Add(lblVerso)
	box4.Add(txtVerso)
	//Ajout des boxs à la fenetre principal
	framebox.PackStart(box1, false, false, 10)
	framebox.PackStart(box2, false, false, 10)
	framebox.PackStart(box4, false, false, 10)
	framebox2.PackStart(box3, false, false, 10)
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}

func getNumberOfPAges(filename string) int {

	cmd := exec.Command("pdftk", filename, "dump_data")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	res := strings.Split(out.String(), "\n")
	var fileInfo map[string]string
	fileInfo = make(map[string]string)
	for i := 0; i < len(res); i++ {
		value := strings.Split(res[i], ": ")
		if len(value) == 2 {
			fileInfo[value[0]] = value[1]
		}
	}
	ret, err := strconv.Atoi(fileInfo["NumberOfPages"])
	if err != nil {
		return 0
	}
	return ret
}

func catFile(fileName, txtVerso *gtk.Entry) string {
	var cmd *exec.Cmd
	var ret string
	if fileName.GetText() != "" {
		if txtVerso.GetText() != "" {
			nbVerso, err := strconv.Atoi(txtVerso.GetText())
			if err != nil {
				log.Fatal(err)
			}
			if nbVerso < nbPages {
				j := nbVerso

				cmdText := " cat"
				for i := 1; i < nbPages; i++ {
					if i < nbVerso {
						cmdText += " " + strconv.Itoa(i)
					}
					if j <= nbPages {
						cmdText += " " + strconv.Itoa(j)
					}
					j++
				}

				if err != nil {
					log.Fatal(err)
				}
				cmd = exec.Command("sh", "-c", "pdftk " + inFile + " " + cmdText +" output " + outFile)
				ko := cmd.Run()

				if ko != nil {
					fmt.Println(ko)
					log.Fatal(ko)
				}
				cat = true
			}
		}
	}
	return ret
}

func displayfile(file string) {
	var cmd *exec.Cmd
	cmd = exec.Command("evince", file)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}


/*
**
**  Create in and out temporary Files for safety manipulations
**
**/

func createTempFile(){
	//creating temp file
	in, err := ioutil.TempFile("", "in.pdf")

	out, err := ioutil.TempFile("", "out.pdf")

	if err != nil {
		log.Fatal(err)
	}

	inFile = in.Name()
	outFile = out.Name()

}

/*
** Copy the content of the input file in the output file
**
** @Params filename string - output file
** @Params input string - input file for content
**/
func saveFile (filename string, input string) {
	in, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filename, in, 0666)
}
