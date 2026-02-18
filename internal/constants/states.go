package constants

// LandesrechtStates maps CLI state values to their API parameter names
// for the Landesrecht endpoint (Bundesland.SucheIn* format).
var LandesrechtStates = map[string]string{
	"wien":              "Bundesland.SucheInWien",
	"niederoesterreich": "Bundesland.SucheInNiederoesterreich",
	"oberoesterreich":   "Bundesland.SucheInOberoesterreich",
	"salzburg":          "Bundesland.SucheInSalzburg",
	"tirol":             "Bundesland.SucheInTirol",
	"vorarlberg":        "Bundesland.SucheInVorarlberg",
	"kaernten":          "Bundesland.SucheInKaernten",
	"steiermark":        "Bundesland.SucheInSteiermark",
	"burgenland":        "Bundesland.SucheInBurgenland",
}

// BezirkeStates maps CLI state values to their API display names
// for the Bezirke endpoint (with Umlauts).
var BezirkeStates = map[string]string{
	"wien":              "Wien",
	"niederoesterreich": "Niederösterreich",
	"oberoesterreich":   "Oberösterreich",
	"salzburg":          "Salzburg",
	"tirol":             "Tirol",
	"vorarlberg":        "Vorarlberg",
	"kaernten":          "Kärnten",
	"steiermark":        "Steiermark",
	"burgenland":        "Burgenland",
}

// VerordnungenStates maps CLI state values to their API display names
// for the Verordnungen endpoint (direct Bundesland values, NOT SucheIn* format).
var VerordnungenStates = map[string]string{
	"wien":              "Wien",
	"niederoesterreich": "Niederösterreich",
	"oberoesterreich":   "Oberösterreich",
	"salzburg":          "Salzburg",
	"tirol":             "Tirol",
	"vorarlberg":        "Vorarlberg",
	"kaernten":          "Kärnten",
	"steiermark":        "Steiermark",
	"burgenland":        "Burgenland",
}
