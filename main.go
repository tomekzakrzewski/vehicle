package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()

	if err != nil {
		log.Fatalln(err)
	}
}

type Model struct {
	title     string
	textinput textinput.Model

	vinResponse VinResponse
	err         error
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter vin here"
	ti.Focus()

	return Model{
		title:     "hujk",
		textinput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEnter:
			v := m.textinput.Value()
			return m, handleVinSearch(v)
		}

	case VinResponseMsg:
		if msg.Err != nil {
			m.err = msg.Err
		}

		m.vinResponse = msg.VinResponse
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := m.textinput.View() + "\n\n"

	if len(m.vinResponse.Results) > 0 {
		s += m.vinResponse.Results[0].ModelYear + "\n\n"
		s += m.vinResponse.Results[0].Make + " " + m.vinResponse.Results[0].Model + "\n\n"
		s += "HP: " + m.vinResponse.Results[0].EngineHP + "\n\n"
		s += "L: " + m.vinResponse.Results[0].DisplacementL + "\n\n"
		s += "Drive type: " + m.vinResponse.Results[0].DriveType + "\n\n"
		s += "Fuel: " + m.vinResponse.Results[0].FuelTypePrimary + "\n\n"
		s += "Seats: " + m.vinResponse.Results[0].Seats + "\n\n"
		s += "Transmission" + m.vinResponse.Results[0].TransmissionStyle + "\n\n"
	}

	return s
}

func handleVinSearch(q string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://vpic.nhtsa.dot.gov/api/vehicles/decodevinvaluesextended/%s?format=json", q)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		if err != nil {
			return VinResponseMsg{
				Err: err,
			}
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return VinResponseMsg{
				Err: err,
			}
		}

		defer res.Body.Close()

		var vinResponse VinResponse

		err = json.NewDecoder(res.Body).Decode(&vinResponse)

		if err != nil {
			return VinResponseMsg{
				Err: err,
			}
		}

		return VinResponseMsg{
			VinResponse: vinResponse,
		}
	}
}

type VinResponse struct {
	Results []struct {
		BodyClass         string `json:"BodyClass"`
		DisplacementL     string `json:"DisplacementL"`
		Doors             string `json:"Doors"`
		DriveType         string `json:"DriveType"`
		EngineCycles      string `json:"EngineCycles"`
		EngineCylinders   string `json:"EngineCylinders"`
		EngineHP          string `json:"EngineHP"`
		FuelTypePrimary   string `json:"FuelTypePrimary"`
		Make              string `json:"Make"`
		Manufacturer      string `json:"Manufacturer"`
		Model             string `json:"Model"`
		ModelYear         string `json:"ModelYear"`
		SeatRows          string `json:"SeatRows"`
		Seats             string `json:"Seats"`
		VehicleType       string `json:"VehicleType"`
		TransmissionStyle string `json:"TransmissionStyle"`
	} `json:"Results"`
}

type VinResponseMsg struct {
	VinResponse VinResponse
	Err         error
}
