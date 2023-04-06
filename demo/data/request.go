package data

type RequestData struct {
	Title    string
	Postcode string
	Info     string
}

var YesterdayRequestNew = RequestData{
	Title:    "Καθαρισμός οικοπέδου",
	Postcode: "19007",
	Info:     "150 τ.μ. κόψιμο κλαδιών, όργωμα",
}

var TwoDaysAgoRequestNew = RequestData{
	Title:    "Κόψιμο δέντρου",
	Postcode: "15387",
	Info:     "Πεύκο που καλύπτει το σπίτι πρέπει να κοπεί. Περίπου 5 μέτρα ύψος.",
}

var TwoDaysAgoRequestOpen = RequestData{
	Title:    "Καθαρισμός σκουπιδιών",
	Postcode: "14660",
	Info:     "Σκουπίδια και μπετά μέσα σε οικόπεδο πρέπει να απομακρυνθούν. Η εργασία απαιτεί φορτηγό",
}

var TwoDaysAgoRequestAssigned = RequestData{
	Title:    "Κούρεμα γκαζόν και κλαδιών",
	Postcode: "14660",
	Info:     "Κούρεμα γκαζόν 50 τ.μ. και κλάδεμα 2 πεύκων ύψους 3 μέτρων",
}
