package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhowden/tag"
	"gmusic/player"
)

var startDir  = flag.String("dir", "", "music directory to open (default: ~/Music or cwd)")
var showStats = flag.Bool("stats", false, "show CPU and memory usage")

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	purple     = lipgloss.Color("#7C3AED")
	darkBorder = lipgloss.Color("#4B5563")
	green      = lipgloss.Color("#10B981")
	yellow     = lipgloss.Color("#F59E0B")
	white      = lipgloss.Color("#F9FAFB")
	gray       = lipgloss.Color("#6B7280")

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkBorder).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().Foreground(purple).Bold(true)
	helpStyle  = lipgloss.NewStyle().Foreground(gray).Padding(0, 1)

	playingStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
	pausedStyle  = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	stoppedStyle = lipgloss.NewStyle().Foreground(gray)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(gray).
				Padding(0, 1)
)

// ── Data ──────────────────────────────────────────────────────────────────────

type track struct {
	path   string
	title  string
	artist string
	album  string
}

var audioExts = map[string]bool{
	".mp3": true, ".wav": true, ".flac": true, ".ogg": true,
}

func loadTracks(dir string) []track {
	var tracks []track
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if !audioExts[strings.ToLower(filepath.Ext(d.Name()))] {
			return nil
		}
		t := track{path: path}
		if f, err := os.Open(path); err == nil {
			if m, err := tag.ReadFrom(f); err == nil {
				t.title = m.Title()
				t.artist = m.Artist()
				t.album = m.Album()
			}
			f.Close()
		}
		if t.title == "" {
			base := filepath.Base(path)
			t.title = strings.TrimSuffix(base, filepath.Ext(base))
		}
		if t.artist == "" {
			t.artist = "Unknown Artist"
		}
		if t.album == "" {
			t.album = "Unknown Album"
		}
		tracks = append(tracks, t)
		return nil
	})
	return tracks
}

// ── List items ────────────────────────────────────────────────────────────────

type trackItem struct{ t track }

func (i trackItem) Title() string       { return i.t.title }
func (i trackItem) Description() string { return i.t.artist }
func (i trackItem) FilterValue() string { return i.t.title + " " + i.t.artist + " " + i.t.album }

type groupItem struct {
	name   string
	tracks []track
}

func (i groupItem) Title() string       { return i.name }
func (i groupItem) Description() string { return fmt.Sprintf("%d tracks", len(i.tracks)) }
func (i groupItem) FilterValue() string { return i.name }

// ── Tabs ──────────────────────────────────────────────────────────────────────

type tabID int

const (
	tabAll tabID = iota
	tabAlbums
	tabArtists
	tabPlaylist
	tabCount
)

var tabNames = [tabCount]string{"All", "Albums", "Artists", "Playlist"}

// ── Playlist persistence ──────────────────────────────────────────────────────

func playlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "gmusic", "playlist.txt")
}

func loadPlaylist() []string {
	data, err := os.ReadFile(playlistPath())
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

func savePlaylist(paths []string) {
	p := playlistPath()
	os.MkdirAll(filepath.Dir(p), 0755)
	os.WriteFile(p, []byte(strings.Join(paths, "\n")+"\n"), 0644)
}

// ── Messages ──────────────────────────────────────────────────────────────────

type tickMsg time.Time
type trackDoneMsg struct{}

// ── Model ─────────────────────────────────────────────────────────────────────

type model struct {
	activeTab tabID
	lists     [tabCount]list.Model

	// Drill-down state (Albums / Artists go one level deep)
	inDrill        [tabCount]bool
	drillBackItems [tabCount][]list.Item
	drillBackTitle [tabCount]string

	tracks       []track
	playlist     []string
	progress     progress.Model
	p            *player.Player
	currentDir   string
	currentTrack string
	queue        []string
	queueIndex   int
	width        int
	height       int
	showStats    bool
	stats        procStats
}

func newDelegate(showDesc bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = showDesc
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(purple).BorderLeftForeground(purple)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(purple)
	return d
}

func newList(showDesc bool) list.Model {
	l := list.New(nil, newDelegate(showDesc), 0, 0)
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	return l
}

func initialModel() model {
	dir := *startDir
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, "Music")
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			dir, _ = os.Getwd()
		}
	}

	tracks := loadTracks(dir)
	playlist := loadPlaylist()

	var lists [tabCount]list.Model
	lists[tabAll] = newList(false)
	lists[tabAll].Title = "All Tracks"
	lists[tabAlbums] = newList(true)
	lists[tabAlbums].Title = "Albums"
	lists[tabArtists] = newList(true)
	lists[tabArtists].Title = "Artists"
	lists[tabPlaylist] = newList(false)
	lists[tabPlaylist].Title = "Playlist"

	m := model{
		lists:      lists,
		tracks:     tracks,
		playlist:   playlist,
		progress:   progress.New(progress.WithDefaultGradient(), progress.WithSolidFill(string(purple))),
		p:          player.New(),
		currentDir: dir,
		showStats:  *showStats,
	}
	m.rebuildLists()
	return m
}

func (m *model) rebuildLists() {
	// All
	allItems := make([]list.Item, len(m.tracks))
	for i, t := range m.tracks {
		allItems[i] = trackItem{t}
	}
	m.lists[tabAll].SetItems(allItems)

	// Albums
	albumMap := map[string][]track{}
	for _, t := range m.tracks {
		albumMap[t.album] = append(albumMap[t.album], t)
	}
	albumItems := make([]list.Item, 0, len(albumMap))
	for _, name := range sortedKeys(albumMap) {
		albumItems = append(albumItems, groupItem{name: name, tracks: albumMap[name]})
	}
	m.lists[tabAlbums].SetItems(albumItems)

	// Artists
	artistMap := map[string][]track{}
	for _, t := range m.tracks {
		artistMap[t.artist] = append(artistMap[t.artist], t)
	}
	artistItems := make([]list.Item, 0, len(artistMap))
	for _, name := range sortedKeys(artistMap) {
		artistItems = append(artistItems, groupItem{name: name, tracks: artistMap[name]})
	}
	m.lists[tabArtists].SetItems(artistItems)

	// Playlist
	trackByPath := map[string]track{}
	for _, t := range m.tracks {
		trackByPath[t.path] = t
	}
	playlistItems := make([]list.Item, 0, len(m.playlist))
	for _, path := range m.playlist {
		t, ok := trackByPath[path]
		if !ok {
			t = track{path: path, title: filepath.Base(path), artist: "", album: ""}
		}
		playlistItems = append(playlistItems, trackItem{t})
	}
	m.lists[tabPlaylist].SetItems(playlistItems)
}

func sortedKeys(m map[string][]track) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tickMsg:
		cmds = append(cmds, tickCmd())
		progCmd := m.progress.SetPercent(m.p.Progress())
		cmds = append(cmds, progCmd)
		if m.showStats {
			m.stats.update()
		}

	case progress.FrameMsg:
		newProgress, cmd := m.progress.Update(msg)
		m.progress = newProgress.(progress.Model)
		cmds = append(cmds, cmd)

	case trackDoneMsg:
		cmds = append(cmds, m.playNext())

	case tea.KeyMsg:
		l := &m.lists[m.activeTab]

		// Let list handle filtering input
		if l.FilterState() == list.Filtering {
			var cmd tea.Cmd
			*l, cmd = l.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.p.Stop()
			return m, tea.Quit

		case "tab":
			m.activeTab = (m.activeTab + 1) % tabCount

		case "shift+tab":
			m.activeTab = (m.activeTab + tabCount - 1) % tabCount

		case "enter":
			switch m.activeTab {
			case tabAll, tabPlaylist:
				if _, ok := l.SelectedItem().(trackItem); ok {
					cmds = append(cmds, m.startQueue(l.Items(), l.Index()))
				}
			case tabAlbums, tabArtists:
				if m.inDrill[m.activeTab] {
					if _, ok := l.SelectedItem().(trackItem); ok {
						cmds = append(cmds, m.startQueue(l.Items(), l.Index()))
					}
				} else {
					if item, ok := l.SelectedItem().(groupItem); ok {
						m.drillInto(item)
					}
				}
			}

		case "esc":
			if m.inDrill[m.activeTab] {
				m.drillBack()
			}

		case " ":
			m.p.PlayPause()

		case "+", "=":
			m.p.Seek(15 * time.Second)

		case "-":
			m.p.Seek(-15 * time.Second)

		case "a":
			if m.currentTrack != "" {
				for _, p := range m.playlist {
					if p == m.currentTrack {
						goto skipAdd
					}
				}
				m.playlist = append(m.playlist, m.currentTrack)
				savePlaylist(m.playlist)
				m.rebuildLists()
			skipAdd:
			}

		case "s":
			m.showStats = !m.showStats
			if m.showStats {
				m.stats.update()
			}

		case "d":
			if m.activeTab == tabPlaylist {
				idx := l.Index()
				if idx >= 0 && idx < len(m.playlist) {
					m.playlist = append(m.playlist[:idx], m.playlist[idx+1:]...)
					savePlaylist(m.playlist)
					m.rebuildLists()
				}
			}

		default:
			var cmd tea.Cmd
			*l, cmd = l.Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		var cmd tea.Cmd
		l := &m.lists[m.activeTab]
		*l, cmd = l.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) drillInto(g groupItem) {
	tab := m.activeTab
	m.drillBackItems[tab] = m.lists[tab].Items()
	m.drillBackTitle[tab] = m.lists[tab].Title
	m.inDrill[tab] = true

	items := make([]list.Item, len(g.tracks))
	for i, t := range g.tracks {
		items[i] = trackItem{t}
	}
	m.lists[tab].SetItems(items)
	m.lists[tab].Title = g.name
}

func (m *model) drillBack() {
	tab := m.activeTab
	m.lists[tab].SetItems(m.drillBackItems[tab])
	m.lists[tab].Title = m.drillBackTitle[tab]
	m.inDrill[tab] = false
}

func (m *model) playFile(path string) tea.Cmd {
	done := make(chan struct{})
	if err := m.p.Load(path, func() { close(done) }); err != nil {
		return nil
	}
	m.currentTrack = path
	return func() tea.Msg {
		<-done
		return trackDoneMsg{}
	}
}

func (m *model) startQueue(items []list.Item, index int) tea.Cmd {
	m.queue = make([]string, 0, len(items))
	for _, item := range items {
		if ti, ok := item.(trackItem); ok {
			m.queue = append(m.queue, ti.t.path)
		}
	}
	m.queueIndex = index
	if index < len(m.queue) {
		return m.playFile(m.queue[index])
	}
	return nil
}

func (m *model) playNext() tea.Cmd {
	m.queueIndex++
	if m.queueIndex < len(m.queue) {
		return m.playFile(m.queue[m.queueIndex])
	}
	return nil
}

func (m *model) updateSizes() {
	leftWidth := (m.width * 60 / 100) - 4
	rightWidth := m.width - leftWidth - 8
	listHeight := m.height - 7 // panel border + tab bar + help bar

	for i := range m.lists {
		m.lists[i].SetSize(leftWidth, listHeight)
	}
	m.progress.Width = rightWidth - 4
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	leftWidth := (m.width * 60 / 100) - 4
	rightWidth := m.width - leftWidth - 8
	innerHeight := m.height - 5

	// Tab bar
	tabBar := m.renderTabBar()

	// Left panel: tab bar + active list
	leftContent := tabBar + "\n" + m.lists[m.activeTab].View()
	leftPanel := panelStyle.
		Width(leftWidth).
		Height(innerHeight).
		Render(leftContent)

	// Right panel: now playing
	rightPanel := panelStyle.
		Width(rightWidth).
		Height(innerHeight).
		Render(m.nowPlayingView(rightWidth - 4))

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)

	// Help bar
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, panels, help)
}

func (m model) renderTabBar() string {
	tabs := make([]string, tabCount)
	for i := tabID(0); i < tabCount; i++ {
		name := tabNames[i]
		if i == m.activeTab {
			tabs[i] = activeTabStyle.Render("● " + name)
		} else {
			tabs[i] = inactiveTabStyle.Render("○ " + name)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m model) renderHelp() string {
	var hint string
	switch {
	case m.inDrill[m.activeTab]:
		hint = "enter:play  esc:back  space:pause  +/-:seek15s  s:stats  q:quit"
	case m.activeTab == tabPlaylist:
		hint = "enter:play  d:remove  space:pause  +/-:seek15s  tab:switch  s:stats  q:quit"
	default:
		hint = "enter:play  space:pause  +/-:seek15s  tab:switch  /:filter  a:playlist  s:stats  q:quit"
	}
	parts := []string{hint}
	return helpStyle.Render(parts[0])
}

func (m model) nowPlayingView(width int) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Now Playing") + "\n\n")

	trackName := "—"
	artistName := ""
	albumName := ""
	if m.currentTrack != "" {
		for _, t := range m.tracks {
			if t.path == m.currentTrack {
				trackName = t.title
				artistName = t.artist
				albumName = t.album
				break
			}
		}
		if trackName == "—" {
			base := filepath.Base(m.currentTrack)
			trackName = strings.TrimSuffix(base, filepath.Ext(base))
		}
	}

	nameStyle := lipgloss.NewStyle().Foreground(white).Bold(true)
	sb.WriteString(nameStyle.Render(truncate(trackName, width)) + "\n")
	if artistName != "" && artistName != "Unknown Artist" {
		sb.WriteString(lipgloss.NewStyle().Foreground(gray).Render(truncate(artistName, width)) + "\n")
	}
	if albumName != "" && albumName != "Unknown Album" {
		sb.WriteString(lipgloss.NewStyle().Foreground(gray).Italic(true).Render(truncate(albumName, width)) + "\n")
	}
	sb.WriteString("\n")

	// State
	state := m.p.State()
	switch state {
	case player.StatePlaying:
		sb.WriteString(playingStyle.Render("▶  Playing"))
	case player.StatePaused:
		sb.WriteString(pausedStyle.Render("⏸  Paused"))
	default:
		sb.WriteString(stoppedStyle.Render("■  Stopped"))
	}
	sb.WriteString("\n\n")

	// Timecode
	pos := m.p.Position()
	dur := m.p.Duration()
	sb.WriteString(lipgloss.NewStyle().Foreground(gray).Render(
		fmt.Sprintf("%s / %s", fmtDuration(pos), fmtDuration(dur)),
	) + "\n\n")

	// Progress bar
	sb.WriteString(m.progress.View() + "\n\n")

	// Stats
	if m.showStats {
		dimStyle := lipgloss.NewStyle().Foreground(gray)
		sb.WriteString("\n" + titleStyle.Render("Stats") + "\n")
		sb.WriteString(dimStyle.Render(fmt.Sprintf("CPU:  %.1f%%\n", m.stats.CPU)))
		sb.WriteString(dimStyle.Render(fmt.Sprintf("RSS:  %s\n", formatBytes(m.stats.MemRSS))))
		sb.WriteString(dimStyle.Render(fmt.Sprintf("Heap: %s\n", formatBytes(m.stats.MemHeap))))
	}

	return sb.String()
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	flag.Parse()
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
