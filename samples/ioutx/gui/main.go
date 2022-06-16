package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Issued response struct
type Issued struct {
	Tokens []struct {
		ID struct {
			TxID string `json:"tx_id"`
		} `json:"id"`
		Owner struct {
			Raw string `json:"raw"`
		} `json:"owner"`
		Type     string `json:"type"`
		Quantity string `json:"quantity"`
		Issuer   struct {
			Raw string `json:"raw"`
		} `json:"Issuer"`
	} `json:"tokens"`
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())

	myWindow := myApp.NewWindow("Digital Token Transactions")
	myWindow.Resize(fyne.NewSize(700, 750))

	// Multiline text
	responsetxt := widget.NewLabel("")
	responsetxt.Wrapping = fyne.TextWrapWord

	// Initialize a list view
	responseList := widget.NewList(func() int {
		return 0
	}, func() fyne.CanvasObject {
		listText := widget.NewEntry()
		listText.MultiLine = true
		listText.Disabled()
		listText.TextStyle.Italic = true

		return listText
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
	})

	// For vertical spacing
	blank := widget.NewLabel("")

	/* ####################### Issuer Query Button ####################### */
	issuedQuery := `"{}"` // Get everything
	// Issuer Query total tokens issued
	btn := widget.NewButton("Total Amount Issued", func() {
		out, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/issuer/client-config.yaml -f issued -i "+issuedQuery).Output()
		if err != nil {
			log.Println(err)
			dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up", myWindow)
			// log.Fatal(err)
		}
		// Map response into struct
		var issued Issued
		err = json.Unmarshal(out, &issued)
		if err != nil {
			dialog.ShowInformation("Error", "Something went wrong! \n"+err.Error(), myWindow)
		}
		// Update list items with content
		responseList.Length = func() int { return len(issued.Tokens) }
		responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			tmp := "Tx_Id:\t" + issued.Tokens[id].ID.TxID +
				"\nCurrency:\t" + issued.Tokens[id].Type + "\nAmount:\t" + hexaNumberToInteger(issued.Tokens[id].Quantity)
			item.(*widget.Entry).SetText(tmp)
		}
		responseList.Refresh()
	})
	// Apply color to Button
	btn_color := canvas.NewRectangle(color.NRGBA{R: 135, G: 206, B: 235, A: 255})
	btn_container := container.New(
		layout.NewMaxLayout(),
		btn_color,
		btn,
	)

	/* ####################### Alice Query Buttons ####################### */
	aliceUnspentQuery := `"{}"`

	btnAliceSpent := widget.NewButton("Total Amount Spent", func() {
		// Update list items with content
		responseList.Length = func() int { return 1 }
		responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Entry).SetText("Fabric token sdk currently does not support queries for amount spent.\nNote: Tokens are deleted from sender's vault as soon as the transaction is finalized.")
		}
		responseList.Refresh()
	})

	btn_alice_unspent := widget.NewButton("Balance (Unspent Amount)", func() {
		outA, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/alice/client-config.yaml -f unspent -i "+aliceUnspentQuery).Output()
		if err != nil {
			log.Println(err)
			dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up", myWindow)
		}
		// Map response into struct
		var issued Issued
		err = json.Unmarshal(outA, &issued)
		if err != nil {
			dialog.ShowInformation("Error", "Something went wrong! \n"+err.Error(), myWindow)
		}

		// Update list items with content
		responseList.Length = func() int { return len(issued.Tokens) }
		responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			tmp := "Tx_Id:\t" + issued.Tokens[id].ID.TxID +
				"\nCurrency:\t" + issued.Tokens[id].Type + "\nAmount:\t" + hexaNumberToInteger(issued.Tokens[id].Quantity)
			item.(*widget.Entry).SetText(tmp)
		}
		responseList.Refresh()
	})
	// Apply color to Button
	btnColorA := canvas.NewRectangle(color.NRGBA{R: 135, G: 206, B: 235, A: 255})
	btnAliceUnspent := container.New(
		layout.NewMaxLayout(),
		btnColorA,
		btn_alice_unspent,
	)

	/* ####################### Bob Query Buttons ####################### */
	bobUnspentQuery := `"{}"`

	btnBobSpent := widget.NewButton("Total Amount Spent", func() {
		// Update list items with content
		responseList.Length = func() int { return 1 }
		responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Entry).SetText("Fabric token sdk currently does not support queries for amount spent.\nNote: Tokens are deleted from sender's vault as soon as the transaction is finalized.")
		}
		responseList.Refresh()
	})

	btn_bob_unspent := widget.NewButton("Balance (Unspent Amount)", func() {
		outB, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/bob/client-config.yaml -f unspent -i "+bobUnspentQuery).Output()
		if err != nil {
			log.Println(err)
			dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up", myWindow)
		}
		// Map response into struct
		var issued Issued
		err = json.Unmarshal(outB, &issued)
		if err != nil {
			dialog.ShowInformation("Error", "Something went wrong! \n"+err.Error(), myWindow)
		}
		// Update list items with content
		responseList.Length = func() int { return len(issued.Tokens) }
		responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			tmp := "Tx_Id:\t" + issued.Tokens[id].ID.TxID +
				"\nCurrency:\t" + issued.Tokens[id].Type + "\nAmount:\t" + hexaNumberToInteger(issued.Tokens[id].Quantity)
			item.(*widget.Entry).SetText(tmp)
		}
		responseList.Refresh()
	})
	// Apply color to Button
	btnColorB := canvas.NewRectangle(color.NRGBA{R: 135, G: 206, B: 235, A: 255})
	btnBobUnspent := container.New(
		layout.NewMaxLayout(),
		btnColorB,
		btn_bob_unspent,
	)

	// Title texts
	issuerTitle := createTitleText("Issue New Tokens")
	sendTitle := createTitleText("Send/Transfer Token")

	// Issuer Form
	currency1, amount1, recipient1 := widget.NewEntry(), widget.NewEntry(), widget.NewEntry()
	currency1.SetPlaceHolder("USD")
	amount1.SetPlaceHolder("10")
	recipient1.SetPlaceHolder("alice")

	issuerForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Currency", Widget: currency1},
			{Text: "Amount", Widget: amount1},
			{Text: "Recipient", Widget: recipient1}},
		OnSubmit: func() { // optional, handle form submission
			if currency1.Text == "" || amount1.Text == "" || recipient1.Text == "" {
				dialog.ShowInformation("Error", "All fields are requied", myWindow)
			} else {
				res := fmt.Sprintf("\\\"TokenType\\\":\\\"%v\\\",\\\"Quantity\\\":%v,\\\"Recipient\\\":\\\"%v\\\"", currency1.Text, amount1.Text, recipient1.Text)
				res = `"{` + res + `}"`

				out, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/issuer/client-config.yaml -f issue -i "+res).Output()
				if err != nil {
					log.Println(err)
					dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up and there is enough fund", myWindow)
				}

				// Update list items with content
				responseList.Length = func() int { return 1 }
				responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*widget.Entry).SetText("Transaction Successful\nTx_Id:\t" + string(out))
				}
				responseList.Refresh()
			}
		},
	}

	// Alice's Form
	currency2, amount2, recipient2 := widget.NewEntry(), widget.NewEntry(), widget.NewEntry()
	currency2.SetPlaceHolder("USD")
	amount2.SetPlaceHolder("5")
	recipient2.SetPlaceHolder("bob")

	aliceForm := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Currency", Widget: currency2},
			{Text: "Amount", Widget: amount2},
			{Text: "Recipient", Widget: recipient2}},
		OnSubmit: func() { // optional, handle form submission
			if currency2.Text == "" || amount2.Text == "" || recipient2.Text == "" {
				dialog.ShowInformation("Error", "All fields are requied", myWindow)
			} else {
				res := fmt.Sprintf("\\\"TokenType\\\":\\\"%v\\\",\\\"Quantity\\\":%v,\\\"Recipient\\\":\\\"%v\\\"", currency2.Text, amount2.Text, recipient2.Text)
				res = `"{` + res + `}"`

				out, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/alice/client-config.yaml -f transfer -i "+res).Output()
				if err != nil {
					log.Println(err)
					dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up and there is enough fund", myWindow)
				}

				// Update list items with content
				responseList.Length = func() int { return 1 }
				responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*widget.Entry).SetText("Transaction Successful\nTx_Id:\t" + string(out))
				}
				responseList.Refresh()
			}
		},
	}

	// Bob's Form
	currency3, amount3, recipient3 := widget.NewEntry(), widget.NewEntry(), widget.NewEntry()
	currency3.SetPlaceHolder("USD")
	amount3.SetPlaceHolder("3")
	recipient3.SetPlaceHolder("charlie")

	bobForm := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Currency", Widget: currency3},
			{Text: "Amount", Widget: amount3},
			{Text: "Recipient", Widget: recipient3}},
		OnSubmit: func() { // optional, handle form submission
			if currency3.Text == "" || amount3.Text == "" || recipient3.Text == "" {
				dialog.ShowInformation("Error", "All fields are requied", myWindow)
			} else {
				res := fmt.Sprintf("\\\"TokenType\\\":\\\"%v\\\",\\\"Quantity\\\":%v,\\\"Recipient\\\":\\\"%v\\\"", currency3.Text, amount3.Text, recipient3.Text)
				res = `"{` + res + `}"`

				out, err := exec.Command("bash", "-c", "../ioutx view -c ../testdata/fsc/nodes/bob/client-config.yaml -f transfer -i "+res).Output()
				if err != nil {
					log.Println(err)
					dialog.ShowInformation("Error", "Something went wrong! \nConfirm network is up and there is enough fund", myWindow)
				}
				// Update list items with content
				responseList.Length = func() int { return 1 }
				responseList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*widget.Entry).SetText("Transaction Successful\nTx_Id:\t" + string(out))
				}
				responseList.Refresh()
			}
		},
	}

	// Arrange view
	issuerContent := container.NewVSplit(container.NewVBox(blank, btn_container, blank, issuerTitle, issuerForm), responseList)
	aliceContent := container.NewVSplit(container.NewVBox(btnAliceSpent, btnAliceUnspent, blank, sendTitle, aliceForm), responseList)
	bobContent := container.NewVSplit(container.NewVBox(btnBobSpent, btnBobUnspent, blank, sendTitle, bobForm), responseList)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Token Issuer", theme.StorageIcon(), issuerContent),
		container.NewTabItemWithIcon("Alice's Wallet", theme.AccountIcon(), aliceContent),
		container.NewTabItemWithIcon("Bob's Wallet", theme.AccountIcon(), bobContent),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

// Convert hex number to integer equivalent
func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	output, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		log.Println(err)
		return ""
	}
	str := strconv.FormatInt(output, 10)

	return str
}

// Create title text
func createTitleText(str string) *canvas.Text {
	title := canvas.NewText(str, color.Black)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 18

	return title
}
