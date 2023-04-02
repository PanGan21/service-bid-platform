package data

import "time"

type RequestData struct {
	Title    string
	Postcode string
	Info     string
	Deadline int64
}

var today = time.Now()
var twoDaysAgo = today.AddDate(0, 0, -2).UTC().UnixMilli()
var yesterday = today.AddDate(0, 0, -1).UTC().UnixMilli()

var YesterdayRequestNew = RequestData{
	Title:    "Καθαρισμός οικοπέδου",
	Postcode: "19007",
	Info:     "150 τ.μ. κόψιμο κλαδιών, όργωμα",
	Deadline: yesterday,
}

var TwoDaysAgoRequestNew = RequestData{
	Title:    "Κόψιμο δέντρου",
	Postcode: "15387",
	Info:     "Πεύκο που καλύπτει το σπίτι πρέπει να κοπεί. Περίπου 5 μέτρα ύψος.",
	Deadline: twoDaysAgo,
}

var TwoDaysAgoRequestOpen = RequestData{
	Title:    "Καθαρισμός σκουπιδιών",
	Postcode: "14660",
	Info:     "Σκουπίδια και μπετά μέσα σε οικόπεδο πρέπει να απομακρυνθούν. Η εργασία απαιτεί φορτηγό",
	Deadline: twoDaysAgo,
}

var TwoDaysAgoRequestAssigned = RequestData{
	Title:    "Κούρεμα γκαζόν και κλαδιών",
	Postcode: "14660",
	Info:     "Κούρεμα γκαζόν 50 τ.μ. και κλάδεμα 2 πεύκων ύψους 3 μέτρων",
	Deadline: twoDaysAgo,
}
