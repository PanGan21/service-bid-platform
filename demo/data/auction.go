package data

import "time"

type AuctionData struct {
	Title    string
	Postcode string
	Info     string
	Deadline int64
}

var today = time.Now()
var twoDaysAgo = today.AddDate(0, 0, -2).UTC().UnixMilli()
var yesterday = today.AddDate(0, 0, -1).UTC().UnixMilli()

var YesterdatAuction = AuctionData{
	Title:    "Καθαρισμός οικοπέδου",
	Postcode: "19007",
	Info:     "150 τ.μ. κόψιμο κλαδιών, όργωμα",
	Deadline: yesterday,
}

var TwoDaysAgoAuction = AuctionData{
	Title:    "Κόψιμο δέντρου",
	Postcode: "15387",
	Info:     "Πεύκο που καλύπτει το σπίτι πρέπει να κοπεί. Περίπου 12 μέτρα ύψος.",
	Deadline: twoDaysAgo,
}
