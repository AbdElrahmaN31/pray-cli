// Package prayer provides prayer times calculation helpers and data
package prayer

// MethodDetails contains detailed information about a calculation method
type MethodDetails struct {
	ID          int
	Name        string
	Description string
	FajrAngle   float64
	IshaAngle   float64
	Region      string
}

// Methods contains detailed information about all calculation methods
var Methods = map[int]MethodDetails{
	0: {
		ID:          0,
		Name:        "Shia Ithna-Ashari",
		Description: "Shia Ithna-Ashari, Leva Institute, Qum",
		FajrAngle:   16.0,
		IshaAngle:   14.0,
		Region:      "Iran",
	},
	1: {
		ID:          1,
		Name:        "University of Islamic Sciences, Karachi",
		Description: "University of Islamic Sciences, Karachi",
		FajrAngle:   18.0,
		IshaAngle:   18.0,
		Region:      "Pakistan, Bangladesh, India, Afghanistan",
	},
	2: {
		ID:          2,
		Name:        "Islamic Society of North America",
		Description: "Islamic Society of North America (ISNA)",
		FajrAngle:   15.0,
		IshaAngle:   15.0,
		Region:      "North America",
	},
	3: {
		ID:          3,
		Name:        "Muslim World League",
		Description: "Muslim World League (MWL)",
		FajrAngle:   18.0,
		IshaAngle:   17.0,
		Region:      "Europe, Far East, parts of USA",
	},
	4: {
		ID:          4,
		Name:        "Umm Al-Qura University, Makkah",
		Description: "Umm Al-Qura University, Makkah",
		FajrAngle:   18.5,
		IshaAngle:   0, // 90 minutes after Maghrib
		Region:      "Arabian Peninsula",
	},
	5: {
		ID:          5,
		Name:        "Egyptian General Authority of Survey",
		Description: "Egyptian General Authority of Survey",
		FajrAngle:   19.5,
		IshaAngle:   17.5,
		Region:      "Africa, Syria, Iraq, Lebanon, Malaysia",
	},
	6: {
		ID:          6,
		Name:        "Institute of Geophysics, University of Tehran",
		Description: "Institute of Geophysics, University of Tehran",
		FajrAngle:   17.7,
		IshaAngle:   14.0,
		Region:      "Iran",
	},
	7: {
		ID:          7,
		Name:        "Gulf Region",
		Description: "Gulf Region",
		FajrAngle:   19.5,
		IshaAngle:   0, // 90 minutes after Maghrib
		Region:      "Gulf Countries",
	},
	8: {
		ID:          8,
		Name:        "Kuwait",
		Description: "Kuwait",
		FajrAngle:   18.0,
		IshaAngle:   17.5,
		Region:      "Kuwait",
	},
	9: {
		ID:          9,
		Name:        "Qatar",
		Description: "Qatar",
		FajrAngle:   18.0,
		IshaAngle:   0, // 90 minutes after Maghrib
		Region:      "Qatar",
	},
	10: {
		ID:          10,
		Name:        "Majlis Ugama Islam Singapura",
		Description: "Majlis Ugama Islam Singapura, Singapore",
		FajrAngle:   20.0,
		IshaAngle:   18.0,
		Region:      "Singapore",
	},
	11: {
		ID:          11,
		Name:        "Union Organization Islamic de France",
		Description: "Union Organization Islamic de France",
		FajrAngle:   12.0,
		IshaAngle:   12.0,
		Region:      "France",
	},
	12: {
		ID:          12,
		Name:        "Diyanet İşleri Başkanlığı",
		Description: "Diyanet İşleri Başkanlığı, Turkey",
		FajrAngle:   18.0,
		IshaAngle:   17.0,
		Region:      "Turkey",
	},
	13: {
		ID:          13,
		Name:        "Spiritual Administration of Muslims of Russia",
		Description: "Spiritual Administration of Muslims of Russia",
		FajrAngle:   16.0,
		IshaAngle:   15.0,
		Region:      "Russia",
	},
}

// GetMethod returns the method details for a given ID
func GetMethod(id int) *MethodDetails {
	if method, ok := Methods[id]; ok {
		return &method
	}
	return nil
}

// GetAllMethods returns all available methods
func GetAllMethods() []MethodDetails {
	methods := make([]MethodDetails, 0, len(Methods))
	for i := 0; i <= 13; i++ {
		if method, ok := Methods[i]; ok {
			methods = append(methods, method)
		}
	}
	return methods
}
