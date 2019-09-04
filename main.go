package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

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

	// TITLE
	label := gtk.NewLabel("Convertir un PDF en Recto/Verso")
	label.ModifyFontEasy("DejaVu Serif 14")

	// label info fichier
	lblPage := gtk.NewLabel("")

	// txtBox file
	lblFileName := gtk.NewLabel("Nom du fichier: ")
	fileName := gtk.NewLabel("")

	lblVerso := gtk.NewLabel("Indiquer la première page Verso: ")
	txtVerso := gtk.NewEntry()
	txtVerso.Connect("changed", func() {
		_, err := strconv.Atoi(txtVerso.GetText())
		if err != nil {
			txtVerso.SetText("")
		}
	})

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
				txtVerso.SetText(strconv.Itoa(int(math.Ceil(float64(nbPages/2)) + 1)))
			}

			//normalizeFileName(fileName.GetText())
			createTempFile()
			saveFile(inFile, fileName.GetText())
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})

	//btSave
	btSave := gtk.NewButtonWithLabel("Enregistrer")
	btSave.Clicked(func() {

		cat = catFile(fileName, txtVerso)

		if cat == true {
			file := strings.Replace(fileName.GetText(), ".pdf", "_recto_verso.pdf", 1)
			saveFile(file, outFile)
			dial := gtk.NewMessageDialog(window, 1, 1, 1, "Le PDF recto/verso généré est "+file)
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
		if fileName.GetText() != "" {
			cat = catFile(fileName, txtVerso)

			if cat {
				displayfile(outFile)
			} else {
				displayfile(inFile)
			}
		}
	})

	//btn Visualiser le fichier original
	btShowOriginalFile := gtk.NewButtonWithLabel("Visualiser l'original")
	btShowOriginalFile.Clicked(func() {
		if fileName.GetText() != "" {
			displayfile(inFile)
		}
	})

	//btSave
	btCancel := gtk.NewButtonWithLabel("Annuler")
	btCancel.Clicked(func() {
		gtk.MainQuit()
	})

	// creation du conteneur Principal
	vbox := gtk.NewVBox(false, 1)
	vpaned := gtk.NewVPaned()
	vpaned.SetSizeRequest(600, 530)
	vboxContent := gtk.NewVBox(false, 1)

	frameFile := gtk.NewFrame("Fichier")
	fBoxFile := gtk.NewVBox(false, 1)
	frameFile.Add(fBoxFile)

	frameOpt := gtk.NewFrame("Options")
	fBoxOpt := gtk.NewVBox(false, 1)
	frameOpt.Add(fBoxOpt)

	vboxSys := gtk.NewHBox(false, 1)
	vboxSys.Add(btCancel)
	vboxSys.Add(btSave)

	vboxContent.Add(frameFile)
	vboxContent.Add(frameOpt)

	fileBox1 := gtk.NewHBox(false, 10)
	fileBox1.Add(lblFileName)
	fileBox1.Add(fileName)

	fileBox2 := gtk.NewHBox(false, 10)
	fileBox2.Add(lblPage)

	fileBox3 := gtk.NewHBox(false, 10)
	fileBox3.Add(btShowOriginalFile)
	fileBox3.Add(btOpenFile)

	fBoxFile.PackStart(fileBox1, false, false, 10)
	fBoxFile.PackStart(fileBox2, false, false, 10)
	fBoxFile.PackStart(fileBox3, false, false, 10)

	optBox1 := gtk.NewHBox(false, 10)
	optBox1.Add(lblVerso)
	optBox1.Add(txtVerso)

	fBoxOpt.PackStart(optBox1, false, false, 10)
	fBoxOpt.PackStart(btShowFile, false, false, 10)

	vpaned.Pack1(vboxContent, false, false)

	vbox.Add(label)
	vbox.Add(vpaned)
	vbox.Add(vboxSys)

	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()

}

func getNumberOfPAges(filename string) int {

	cmd := exec.Command("qpdf", filename, "--show-npages")
	var out, err = cmd.Output()

	nb := strings.Split(string(out), "\n")

	ret, err := strconv.Atoi(nb[0])
	if err != nil {
		return 0
	}

	return ret
}

func catFile(fileName *gtk.Label, txtVerso *gtk.Entry) bool {
	var cmd *exec.Cmd

	if fileName.GetText() != "" {

		nbVerso, _ := strconv.Atoi(txtVerso.GetText())

		if nbVerso < nbPages {
			j := nbVerso

			cmdText := ""
			for i := 1; i < nbPages; i++ {
				if i < nbVerso {
					cmdText += strconv.Itoa(i) + ","
				}
				if j < nbPages {
					cmdText += strconv.Itoa(j) + ","
				} else if j == nbPages {
					cmdText += strconv.Itoa(j)
				}
				j++
			}

			cmd = exec.Command("bash", "-c", "qpdf --empty --pages "+inFile+" "+cmdText+" -- "+outFile)
			ko := cmd.Run()

			if ko != nil {
				return false
			}
			return true
		}
	}
	return false
}

func displayfile(file string) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", file)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("start", file)
	}
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

func createTempFile() {

	//creating temp file
	in, err := ioutil.TempFile("", "in")

	out, err := ioutil.TempFile("", "out")

	if err != nil {
		log.Fatal(err)
	}

	if inFile != "" && outFile != "" {
		os.Remove(inFile)
		os.Remove(outFile)
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
func saveFile(filename string, input string) {
	in, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filename, in, 0666)
}
