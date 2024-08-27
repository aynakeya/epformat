package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func createGuiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gui",
		Short: "start gui",
		Long:  "start gui",
		Run: func(cmd *cobra.Command, args []string) {
			startGui()
		},
	}
	return cmd
}

func startGui() {
	a := app.NewWithID("com.aynakeya.epformat")
	w := a.NewWindow("Anime Episode Formatter")

	title := ""
	titleBinding := binding.BindString(&title)
	// Title Entry
	titleEntry := widget.NewEntryWithData(titleBinding)
	titleEntry.SetPlaceHolder("Enter Title")

	season := -1

	// Season Entry
	seasonEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&season)))
	seasonEntry.SetPlaceHolder("Enter Season Number")

	episode := -1

	// Episode Entry
	episodeEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&episode)))
	episodeEntry.SetPlaceHolder("Enter Episode Number")

	format := DefaultFormat
	// Format Entry
	formatEntry := widget.NewEntryWithData(binding.BindString(&format))
	formatEntry.SetPlaceHolder("Enter Format")

	filefolderPath := ""
	fildFolderPathBinding := binding.BindString(&filefolderPath)
	filefolderEntry := widget.NewEntryWithData(fildFolderPathBinding)
	formatEntry.SetPlaceHolder("Enter file path or folder path")

	// Table to display original and renamed filenames
	data := []fileInfo{}
	renamedInfo := []string{}

	applyRules := func() {
		renamedInfo = make([]string, len(data))
		for i, finfo := range data {
			fileName := finfo.info.Name()
			renamed, err := RenameEpInfo(MainExtractor, fileName, title, season, episode, format)
			if err != nil {
				renamedInfo[i] = ""
			} else {
				renamedInfo[i] = renamed
			}
		}
		if len(data) > 0 && title == "" {
			title = MainExtractor.Extract(data[0].info.Name()).Title
			_ = titleBinding.Reload()
		}
	}

	tableHeader := container.NewGridWithColumns(2,
		widget.NewLabel("Original Filename"),
		widget.NewLabel("Renamed"))

	table := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(2,
				widget.NewLabel("Original Filename"),
				widget.NewLabel("Renamed"))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			fileName := data[id].info.Name()
			object.(*fyne.Container).Objects[0].(*widget.Label).SetText(fileName)
			if renamedInfo[id] == "" {
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Error")
			} else {
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(renamedInfo[id])
			}
		},
	)

	// File selection button
	fileSelectionButton := widget.NewButton("Select", func() {
		var selectionDialog dialog.Dialog

		selectionDialog = dialog.NewCustom("Select File or Folder", "Cancel", container.NewHBox(
			widget.NewButton("Select File", func() {
				dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err != nil || uc == nil {
						return
					}
					filefolderPath = uc.URI().Path()
					_ = fildFolderPathBinding.Reload()
					selectionDialog.Hide()
				}, w).Show()
			}),
			widget.NewButton("Select Folder", func() {
				dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
					if err != nil || list == nil {
						return
					}
					filefolderPath = list.Path()
					_ = fildFolderPathBinding.Reload()
					selectionDialog.Hide()
				}, w).Show()
			}),
		), w)
		selectionDialog.Show()
	})
	filefolderBtns := container.NewHBox(
		fileSelectionButton,
		widget.NewButton("Reload", func() {
			title = ""
			data = getAllFiles(filefolderPath)
			applyRules()
			table.Refresh()
		}))
	var renameBtn *widget.Button
	renameBtn = widget.NewButton("Rename", func() {
		renameBtn.Disable()
		applyRules()
		table.Refresh()
		renamed := 0
		failed := 0
		for i, file := range data {
			if renamedInfo[0] == "" {
				continue
			}
			err := os.Rename(file.path, filepath.Join(filepath.Dir(file.path), removeSpecialChars(renamedInfo[i])))
			if err != nil {
				fmt.Printf("error occurs, %v", err)
				failed += 1
			} else {
				renamed += 1
			}
		}
		dialog.ShowConfirm("Rename finished", fmt.Sprintf("Total %d, successed %d, failed %d", renamed+failed, renamed, failed), func(v bool) {}, w)
		renameBtn.Enable()
	})
	formatBtns := container.NewHBox(
		widget.NewButton("Apply", func() {
			applyRules()
			table.Refresh()
		}), renameBtn,
	)

	// Main content layout
	inputForm := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("File/Folder"),
		container.NewBorder(nil, nil, nil, filefolderBtns, filefolderEntry),
		widget.NewLabel("Format"),
		container.NewBorder(nil, nil, nil, formatBtns, formatEntry),
		widget.NewLabel("Title"),
		titleEntry,
		widget.NewLabel("Season"),
		seasonEntry,
		widget.NewLabel("Episode"),
		episodeEntry,
	)
	content := container.NewVBox(
		inputForm,
		container.NewVBox(tableHeader, widget.NewSeparator(), table))

	w.SetContent(content)
	w.Resize(fyne.NewSize(720, 480))
	w.ShowAndRun()
}
