package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"ipbak/internal/geo"
)

type Printer interface {
	Print(loc *geo.Location) error
}

type TextPrinter struct {
	Advanced bool
}

func (p *TextPrinter) Print(loc *geo.Location) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "IP Address:\t%s\n", loc.IP)
	fmt.Fprintf(w, "IP Type:\t%s\n", loc.IPType)
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Location")
	fmt.Fprintf(w, "\tCountry:\t%s (%s)\n", loc.Country, loc.CountryCode)
	if loc.Region != "" {
		fmt.Fprintf(w, "\tRegion:\t%s\n", loc.Region)
	}
	if loc.City != "" {
		fmt.Fprintf(w, "\tCity:\t%s\n", loc.City)
	}
	if p.Advanced {
		fmt.Fprintf(w, "\tPostal Code:\t%s\n", loc.PostalCode)
		fmt.Fprintf(w, "\tContinent:\t%s (%s)\n", loc.Continent, loc.ContinentCode)
	}
	fmt.Fprintf(w, "\tTimezone:\t%s\n", loc.Timezone)
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Network")
	if loc.ASN != 0 {
		if p.Advanced {
			fmt.Fprintf(w, "\tASN:\tAS%d\n", loc.ASN)
		} else {
			fmt.Fprintf(w, "\tASN:\tAS%d\n", loc.ASN)
		}
	}
	fmt.Fprintf(w, "\tOrganization:\t%s\n", loc.Organization)
	if p.Advanced {
		fmt.Fprintf(w, "\tClassification:\t%s\n", loc.NetworkClass)
	}
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Coordinates")
	fmt.Fprintf(w, "\tLatitude:\t%.4f\n", loc.Latitude)
	fmt.Fprintf(w, "\tLongitude:\t%.4f\n", loc.Longitude)
	if p.Advanced {
		fmt.Fprintf(w, "\tAccuracy:\t%dkm\n", loc.AccuracyRadius)
	}
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Map:")
	fmt.Fprintf(w, "\thttps://www.google.com/maps?q=%.4f,%.4f\n", loc.Latitude, loc.Longitude)
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Note: Location data is approximate.")
	return w.Flush()
}

type JsonPrinter struct{}

func (p *JsonPrinter) Print(loc *geo.Location) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(loc)
}
