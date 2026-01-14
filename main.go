package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Country struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Capital    []string `json:"capital"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Population int64    `json:"population"`
	Area       float64  `json:"area"`
	Flags      struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
	} `json:"flags"`
	Languages  map[string]string `json:"languages"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Timezones  []string `json:"timezones"`
	Continents []string `json:"continents"`
	Borders    []string `json:"borders"`
	Landlocked bool     `json:"landlocked"`
	Maps       struct {
		GoogleMaps     string `json:"googleMaps"`
		OpenStreetMaps string `json:"openStreetMaps"`
	} `json:"maps"`
}

type CountryEvent struct {
	Name       string `json:"name"`
	Official   string `json:"official"`
	Capital    string `json:"capital"`
	Region     string `json:"region"`
	Subregion  string `json:"subregion"`
	Population string `json:"population"`
	Area       string `json:"area"`
	Flag       string `json:"flag"`
	Languages  string `json:"languages"`
	Currency   string `json:"currency"`
	Timezone   string `json:"timezone"`
	Continent  string `json:"continent"`
	Borders    int    `json:"borders"`
	Landlocked bool   `json:"landlocked"`
	MapURL     string `json:"mapUrl"`
	Timestamp  string `json:"timestamp"`
}

var countries []Country

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Println("🌍 Fetching countries data...")
	if err := fetchCountries(); err != nil {
		log.Fatal("Failed to fetch countries: ", err)
	}
	log.Printf("✅ Loaded %d countries\n", len(countries))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleHome)
	r.Get("/api/stream", handleSSE)
	r.Get("/api/random", handleRandom)
	r.Get("/api/stats", handleStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Country Pulse running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func fetchCountries() error {
	client := &http.Client{Timeout: 30 * time.Second}
	
	// Get all independent countries with full data
	resp, err := client.Get("https://restcountries.com/v3.1/independent?status=true")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(&countries)
}

func getRandomCountry() CountryEvent {
	c := countries[rand.Intn(len(countries))]

	capital := "N/A"
	if len(c.Capital) > 0 {
		capital = c.Capital[0]
	}

	var languages string
	for _, lang := range c.Languages {
		if languages != "" {
			languages += ", "
		}
		languages += lang
	}
	if languages == "" {
		languages = "N/A"
	}

	var currency string
	for _, curr := range c.Currencies {
		if curr.Symbol != "" {
			currency = fmt.Sprintf("%s (%s)", curr.Name, curr.Symbol)
		} else {
			currency = curr.Name
		}
		break
	}
	if currency == "" {
		currency = "N/A"
	}

	timezone := "N/A"
	if len(c.Timezones) > 0 {
		timezone = c.Timezones[0]
	}

	continent := "N/A"
	if len(c.Continents) > 0 {
		continent = c.Continents[0]
	}

	flag := c.Flags.Svg
	if flag == "" {
		flag = c.Flags.Png
	}

	return CountryEvent{
		Name:       c.Name.Common,
		Official:   c.Name.Official,
		Capital:    capital,
		Region:     c.Region,
		Subregion:  c.Subregion,
		Population: formatNumber(c.Population),
		Area:       formatArea(c.Area),
		Flag:       flag,
		Languages:  languages,
		Currency:   currency,
		Timezone:   timezone,
		Continent:  continent,
		Borders:    len(c.Borders),
		Landlocked: c.Landlocked,
		MapURL:     c.Maps.GoogleMaps,
		Timestamp:  time.Now().Format("15:04:05"),
	}
}

func formatNumber(n int64) string {
	if n >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(n)/1000000000)
	}
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func formatArea(a float64) string {
	if a >= 1000000 {
		return fmt.Sprintf("%.2fM km²", a/1000000)
	}
	if a >= 1000 {
		return fmt.Sprintf("%.1fK km²", a/1000)
	}
	return fmt.Sprintf("%.0f km²", a)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	country := getRandomCountry()
	data, _ := json.Marshal(country)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	for {
		select {
		case <-ticker.C:
			country := getRandomCountry()
			data, _ := json.Marshal(country)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func handleRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRandomCountry())
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	regions := make(map[string]int)
	for _, c := range countries {
		regions[c.Region]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":   len(countries),
		"regions": regions,
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	tmpl.Execute(w, nil)
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Country Pulse - Real-time Country Discovery</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        body { font-family: 'Inter', sans-serif; }
        .globe-pulse { animation: pulse 2s ease-in-out infinite; }
        @keyframes pulse {
            0%, 100% { transform: scale(1); opacity: 1; }
            50% { transform: scale(1.05); opacity: 0.8; }
        }
        .card-enter { animation: slideUp 0.5s ease-out; }
        @keyframes slideUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .flag-shadow { box-shadow: 0 10px 40px rgba(0,0,0,0.3); }
        .stat-glow { box-shadow: 0 0 30px rgba(59, 130, 246, 0.1); }
    </style>
</head>
<body class="bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 min-h-screen text-white">
    <div class="container mx-auto px-4 py-8">
        <header class="text-center mb-12">
            <div class="inline-flex items-center gap-4 mb-4">
                <div class="globe-pulse text-6xl">🌍</div>
                <div>
                    <h1 class="text-4xl font-bold bg-gradient-to-r from-blue-400 via-cyan-400 to-emerald-400 bg-clip-text text-transparent">
                        Country Pulse
                    </h1>
                    <p class="text-slate-400 mt-1">Real-time Country Discovery Stream</p>
                </div>
            </div>
            <div id="connection" class="inline-flex items-center gap-2 px-4 py-2 bg-slate-800/50 rounded-full text-sm">
                <span class="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></span>
                <span class="text-slate-300">Streaming live discoveries</span>
            </div>
        </header>

        <div id="stats" class="flex justify-center gap-6 mb-10 flex-wrap">
            <div class="stat-glow bg-slate-800/60 backdrop-blur px-6 py-3 rounded-2xl border border-slate-700/50">
                <span class="text-slate-400 text-sm">Total Countries</span>
                <span id="totalCountries" class="ml-2 text-xl font-bold text-cyan-400">---</span>
            </div>
            <div class="stat-glow bg-slate-800/60 backdrop-blur px-6 py-3 rounded-2xl border border-slate-700/50">
                <span class="text-slate-400 text-sm">Discovered</span>
                <span id="discovered" class="ml-2 text-xl font-bold text-emerald-400">0</span>
            </div>
            <div class="stat-glow bg-slate-800/60 backdrop-blur px-6 py-3 rounded-2xl border border-slate-700/50">
                <span class="text-slate-400 text-sm">Current Region</span>
                <span id="currentRegion" class="ml-2 text-xl font-bold text-amber-400">---</span>
            </div>
        </div>

        <div id="countryCard" class="max-w-4xl mx-auto">
            <div class="text-center py-20 text-slate-500">
                <div class="text-5xl mb-4">🌐</div>
                <p>Waiting for first discovery...</p>
            </div>
        </div>

        <div class="mt-16 max-w-4xl mx-auto">
            <h2 class="text-xl font-semibold text-slate-300 mb-4 flex items-center gap-2">
                <span>📜</span> Recent Discoveries
            </h2>
            <div id="history" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3"></div>
        </div>

        <footer class="text-center mt-16 text-slate-500 text-sm">
            <p>Powered by <a href="https://restcountries.com" class="text-cyan-400 hover:underline" target="_blank">REST Countries API</a></p>
            <p class="mt-1">Built with Go + SSE + Tailwind</p>
        </footer>
    </div>

    <script>
        let discoveredCount = 0;
        let history = [];
        const maxHistory = 12;

        function formatCard(c) {
            return ` + "`" + `
                <div class="card-enter bg-gradient-to-br from-slate-800/90 to-slate-900/90 backdrop-blur-xl rounded-3xl border border-slate-700/50 overflow-hidden">
                    <div class="relative bg-gradient-to-br from-slate-700/30 to-slate-800/30 p-8 flex justify-center">
                        <img src="${c.flag}" alt="${c.name} flag" class="h-40 w-auto flag-shadow rounded-lg object-cover">
                        <div class="absolute top-4 right-4 bg-slate-900/80 backdrop-blur px-3 py-1 rounded-full text-xs text-slate-300">${c.timestamp}</div>
                    </div>
                    <div class="p-8">
                        <div class="text-center mb-6">
                            <h2 class="text-3xl font-bold text-white mb-2">${c.name}</h2>
                            <p class="text-slate-400 text-sm">${c.official}</p>
                            <div class="mt-3 inline-flex items-center gap-2 flex-wrap justify-center">
                                <span class="px-3 py-1 bg-blue-500/20 text-blue-300 rounded-full text-sm">${c.region || 'N/A'}</span>
                                <span class="px-3 py-1 bg-cyan-500/20 text-cyan-300 rounded-full text-sm">${c.subregion || 'N/A'}</span>
                            </div>
                        </div>
                        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                            <div class="bg-slate-800/50 rounded-xl p-4 text-center">
                                <div class="text-2xl mb-1">🏛️</div>
                                <div class="text-xs text-slate-400 mb-1">Capital</div>
                                <div class="text-white font-semibold text-sm">${c.capital}</div>
                            </div>
                            <div class="bg-slate-800/50 rounded-xl p-4 text-center">
                                <div class="text-2xl mb-1">👥</div>
                                <div class="text-xs text-slate-400 mb-1">Population</div>
                                <div class="text-white font-semibold text-sm">${c.population}</div>
                            </div>
                            <div class="bg-slate-800/50 rounded-xl p-4 text-center">
                                <div class="text-2xl mb-1">📐</div>
                                <div class="text-xs text-slate-400 mb-1">Area</div>
                                <div class="text-white font-semibold text-sm">${c.area}</div>
                            </div>
                            <div class="bg-slate-800/50 rounded-xl p-4 text-center">
                                <div class="text-2xl mb-1">🌐</div>
                                <div class="text-xs text-slate-400 mb-1">Continent</div>
                                <div class="text-white font-semibold text-sm">${c.continent}</div>
                            </div>
                        </div>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                            <div class="bg-slate-800/30 rounded-xl p-4">
                                <div class="flex items-center gap-2 text-slate-400 text-sm mb-2"><span>🗣️</span> Languages</div>
                                <div class="text-white">${c.languages}</div>
                            </div>
                            <div class="bg-slate-800/30 rounded-xl p-4">
                                <div class="flex items-center gap-2 text-slate-400 text-sm mb-2"><span>💰</span> Currency</div>
                                <div class="text-white">${c.currency}</div>
                            </div>
                            <div class="bg-slate-800/30 rounded-xl p-4">
                                <div class="flex items-center gap-2 text-slate-400 text-sm mb-2"><span>🕐</span> Timezone</div>
                                <div class="text-white">${c.timezone}</div>
                            </div>
                            <div class="bg-slate-800/30 rounded-xl p-4">
                                <div class="flex items-center gap-2 text-slate-400 text-sm mb-2"><span>🗺️</span> Borders</div>
                                <div class="text-white">${c.borders} neighbor${c.borders === 1 ? '' : 's'} ${c.landlocked ? '(🏔️ Landlocked)' : ''}</div>
                            </div>
                        </div>
                        <div class="flex justify-center">
                            <a href="${c.mapUrl}" target="_blank" class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-500 hover:to-cyan-500 rounded-xl font-medium transition-all transform hover:scale-105">
                                <span>🗺️</span> View on Map
                            </a>
                        </div>
                    </div>
                </div>
            ` + "`" + `;
        }

        function addToHistory(c) {
            history.unshift(c);
            if (history.length > maxHistory) history.pop();
            document.getElementById('history').innerHTML = history.map(x => ` + "`" + `
                <div class="bg-slate-800/60 rounded-xl p-3 text-center hover:bg-slate-700/60 transition-all border border-slate-700/30">
                    <img src="${x.flag}" alt="${x.name}" class="h-8 w-auto mx-auto rounded shadow-md mb-2">
                    <div class="text-xs text-white truncate">${x.name}</div>
                    <div class="text-xs text-slate-500">${x.timestamp}</div>
                </div>
            ` + "`" + `).join('');
        }

        async function loadStats() {
            const resp = await fetch('/api/stats');
            const data = await resp.json();
            document.getElementById('totalCountries').textContent = data.total;
        }

        function connect() {
            const source = new EventSource('/api/stream');
            source.onmessage = (e) => {
                const c = JSON.parse(e.data);
                discoveredCount++;
                document.getElementById('countryCard').innerHTML = formatCard(c);
                document.getElementById('discovered').textContent = discoveredCount;
                document.getElementById('currentRegion').textContent = c.region || 'Unknown';
                addToHistory(c);
            };
            source.onerror = () => {
                document.getElementById('connection').innerHTML = ` + "`" + `
                    <span class="w-2 h-2 bg-red-500 rounded-full"></span>
                    <span class="text-red-400">Reconnecting...</span>
                ` + "`" + `;
                setTimeout(connect, 3000);
            };
        }

        loadStats();
        connect();
    </script>
</body>
</html>`
