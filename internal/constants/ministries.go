package constants

// RegvorlMinistries maps CLI ministry values to their API values
// for the regvorl command (format: "ABBR (Full Name)").
var RegvorlMinistries = map[string]string{
	"bka":    "BKA (Bundeskanzleramt)",
	"bmkoes": "BMKOES (Bundesministerium für Kunst, Kultur, öffentlichen Dienst und Sport)",
	"bmeia":  "BMEIA (Bundesministerium für europäische und internationale Angelegenheiten)",
	"bmaw":   "BMAW (Bundesministerium für Arbeit und Wirtschaft)",
	"bmbwf":  "BMBWF (Bundesministerium für Bildung, Wissenschaft und Forschung)",
	"bmf":    "BMF (Bundesministerium für Finanzen)",
	"bmi":    "BMI (Bundesministerium für Inneres)",
	"bmj":    "BMJ (Bundesministerium für Justiz)",
	"bmk":    "BMK (Bundesministerium für Klimaschutz, Umwelt, Energie, Mobilität, Innovation und Technologie)",
	"bmlv":   "BMLV (Bundesministerium für Landesverteidigung)",
	"bml":    "BML (Bundesministerium für Land- und Forstwirtschaft, Regionen und Wasserwirtschaft)",
	"bmsgpk": "BMSGPK (Bundesministerium für Soziales, Gesundheit, Pflege und Konsumentenschutz)",
	"bmffim": "BMFFIM (Bundesministerin für Frauen, Familie, Integration und Medien im Bundeskanzleramt)",
	"bmeuv":  "BMEUV (Bundesministerin für EU und Verfassung im Bundeskanzleramt)",
}

// ErlMinistries maps CLI ministry values to their API values
// for the erlaesse sub-command (full name only, no abbreviation prefix).
var ErlMinistries = map[string]string{
	"bka":    "Bundeskanzleramt",
	"bmkoes": "Bundesministerium für Kunst, Kultur, öffentlichen Dienst und Sport",
	"bmeia":  "Bundesministerium für europäische und internationale Angelegenheiten",
	"bmaw":   "Bundesministerium für Arbeit und Wirtschaft",
	"bmbwf":  "Bundesministerium für Bildung, Wissenschaft und Forschung",
	"bmf":    "Bundesministerium für Finanzen",
	"bmi":    "Bundesministerium für Inneres",
	"bmj":    "Bundesministerium für Justiz",
	"bmk":    "Bundesministerium für Klimaschutz, Umwelt, Energie, Mobilität, Innovation und Technologie",
	"bmlv":   "Bundesministerium für Landesverteidigung",
	"bml":    "Bundesministerium für Land- und Forstwirtschaft, Regionen und Wasserwirtschaft",
	"bmsgpk": "Bundesministerium für Soziales, Gesundheit, Pflege und Konsumentenschutz",
}
