package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles with LipGloss
var (
	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5DADE2")).
		Padding(0, 2)

	todoBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(78).
		BorderForeground(
			lipgloss.Color("#5DADE2")).
		Padding(0, 2)

	titleStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#F39C12"))

	weatherStyle = lipgloss.NewStyle().
		Width(20).
		Foreground(lipgloss.Color("#2ECC71"))

	repoStyle = lipgloss.NewStyle().
		Width(48).
		Foreground(lipgloss.Color("#2ECC71"))

	newsStyle = lipgloss.NewStyle().
		Width(74).
		Foreground(lipgloss.Color("#2ECC71"))

	quoteStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(74).
		Foreground(lipgloss.Color("#AF7AC5"))
)

type model struct {
	weatherData   WeatherResponse
	githubRepos   []Repo
	newsHeadlines []string
	choices       []string
	cursor        int
	selected      map[int]struct{}
	inputMode     bool
	inputBuffer   string
	CPUPercent    float32
}

type tickMsg struct{}

type WeatherResponse struct {
	Timezone       string `json:"timezone"`
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
	} `json:"current_weather"`
	Hourly struct {
		Time                     []string  `json:"time"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
	} `json:"hourly"`
}

type NewsResponse struct {
	Status   string `json:"status"`
	Total    int    `json:"totalResults"`
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
}

type Repo struct {
	Name      string    `json:"name"`
	HTMLURL   string    `json:"html_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

func initialModel() model {
	return model{
		choices:  []string{},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadData(),
		tick(),
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Hour*2, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func loadData() tea.Cmd {
	return func() tea.Msg {
		m := struct {
			Weather WeatherResponse
			Repos   []Repo
			News    []string
		}{}

		// Weather
		if data, err := fetchWeather(); err == nil {
			m.Weather = data
		}

		// GitHub Repos
		if repos, err := fetchGitHubRepos("EffexDev", ""); err == nil {
			m.Repos = repos
		}

		// News
		if news, err := fetchTechNews(); err == nil {
			m.News = news
		}

		return m
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			switch msg.Type {
			case tea.KeyEnter:
				// Save the item and exit input mode
				if len(m.inputBuffer) > 0 {
					m.choices = append(m.choices, m.inputBuffer)
				}
				m.inputMode = false
				m.inputBuffer = ""
			case tea.KeyEsc:
				// Cancel input
				m.inputMode = false
				m.inputBuffer = ""
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			default:
				if msg.String() != "" && len(msg.String()) == 1 {
					m.inputBuffer += msg.String()
				}
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "i":
				// Enter input mode
				m.inputMode = true
				m.inputBuffer = ""
			case "enter", " ":
				// Toggle selection
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		}
	case struct {
		Weather WeatherResponse
		Repos   []Repo
		News    []string
	}:
		m.weatherData = msg.Weather
		m.githubRepos = msg.Repos
		m.newsHeadlines = msg.News
	case tickMsg:
		data, err := fetchWeather()
		if err == nil {
			m.weatherData = data
		}
		repos, err := fetchGitHubRepos("EffexDev", "")
		if err == nil {
			m.githubRepos = repos
		}
		headlines, err := fetchTechNews()
		if err == nil {
			m.newsHeadlines = headlines
		}
		return m, tick()
	}
	return m, nil
}

func fetchWeather() (WeatherResponse, error) {
	var data WeatherResponse

	url := "https://api.open-meteo.com/v1/forecast?latitude=-31.9514&longitude=115.8617&current_weather=true&hourly=precipitation_probability"
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func fetchTechNews() ([]string, error) {
	apiKey := "af7c72b12c4d45d091f800eeb33e295f"
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?sources=techcrunch&apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("News API error: %s", resp.Status)
	}

	var news NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	var headlines []string
	for _, a := range news.Articles {
		headlines = append(headlines, fmt.Sprintf("%s", a.Title))
	}

	return headlines, nil
}

func fetchGitHubRepos(username, token string) ([]Repo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Go-GitHub-Client")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (m model) View() string {
	weather := m.weatherData.CurrentWeather
	wind := fmt.Sprintf("%.1f", weather.Windspeed)
	temp := fmt.Sprintf("%.1f", weather.Temperature)
	rain := "N/A"
	if len(m.weatherData.Hourly.PrecipitationProbability) > 0 {
		rain = fmt.Sprintf("%.1f", m.weatherData.Hourly.PrecipitationProbability[0])
	}
	weatherLines := []string{
		titleStyle.Render("Weather Pane"),
		weatherStyle.Render("Temp: " + temp + " Celsius"),
		weatherStyle.Render("Wind: " + wind + " km/hr"),
		weatherStyle.Render("Rain: " + rain + " %"),
	}
	weatherBox := borderStyle.Render(lipgloss.JoinVertical(lipgloss.Top, weatherLines...))

	title := "Effex Dashboard"
	var styledQuotes []string
	styledQuotes = append(styledQuotes, quoteStyle.Render(title))

	repoTitle := titleStyle.Render("Updated Repos")
	var repoLines []string
	repoLines = append(repoLines, repoTitle)
	for i, r := range m.githubRepos {
		if i >= 3 {
			break
		}
		repoLines = append(repoLines, repoStyle.Render(r.Name))
	}
	reposBox := borderStyle.Render(lipgloss.JoinVertical(lipgloss.Top, repoLines...))

	var newsLines []string
	newsLines = append(newsLines, titleStyle.Render("TechCrunch Headlines"))
	if len(m.newsHeadlines) == 0 {
		newsLines = append(newsLines, "Loading news...")
	} else {
		for i, h := range m.newsHeadlines {
			if i >= 3 { // limit to 5 headlines
				break
			}
			newsLines = append(newsLines, newsStyle.Render(h+"\n"))
		}
	}
	newsBox := borderStyle.Render(lipgloss.JoinVertical(lipgloss.Top, newsLines...))

	todoTitle := titleStyle.Render("To-Do List")
	var todoList []string
	todoList = append(todoList, todoTitle)
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✓"
		}
		line := fmt.Sprintf("%s [%s] %s", cursor, checked, choice)
		todoList = append(todoList, line)
	}
	todoBox := todoBorderStyle.Render(lipgloss.JoinVertical(lipgloss.Top, todoList...))

	titleBox := borderStyle.Render(lipgloss.JoinVertical(lipgloss.Top, styledQuotes...))
	firstRow := lipgloss.JoinHorizontal(lipgloss.Top, newsBox)
	secondRow := lipgloss.JoinHorizontal(lipgloss.Top, weatherBox, reposBox)
	thirdRow := lipgloss.JoinHorizontal(lipgloss.Top, todoBox)
	footer := "Press q to exit app. Press i to add item."
	if m.inputMode {
		footer = "Adding new item: " + m.inputBuffer + "\n\n(Enter=Save, Esc=Cancel)\n\n"
	}
	return lipgloss.JoinVertical(lipgloss.Top, titleBox, firstRow, secondRow, thirdRow, footer)
}

func ClearTerminal() {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "cls"}
	default:
		cmd = "clear"
		args = []string{}
	}
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	tick()
	ClearTerminal()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
