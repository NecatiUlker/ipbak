package geo

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"ipbak/internal/config"
)

type Location struct {
	IP             string
	IPType         string
	Country        string
	CountryCode    string
	Region         string
	City           string
	PostalCode     string
	Continent      string
	ContinentCode  string
	Timezone       string
	ASN            uint
	Organization   string
	Latitude       float64
	Longitude      float64
	AccuracyRadius uint16
	NetworkClass   string // "Hosting Provider" or "Likely Residential/ISP"
}

type GeoService struct {
	asnDB  *geoip2.Reader
	cityDB *geoip2.Reader
}

func NewGeoService() (*GeoService, error) {
	dbDir, err := config.GetDatabaseDir()
	if err != nil {
		return nil, err
	}

	asnPath := filepath.Join(dbDir, "GeoLite2-ASN.mmdb")
	cityPath := filepath.Join(dbDir, "GeoLite2-City.mmdb")

	var asnDB *geoip2.Reader
	if _, err := os.Stat(asnPath); err == nil {
		asnDB, err = geoip2.Open(asnPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open ASN database: %w", err)
		}
	}

	var cityDB *geoip2.Reader
	if _, err := os.Stat(cityPath); err == nil {
		cityDB, err = geoip2.Open(cityPath)
		if err != nil {
			if asnDB != nil {
				asnDB.Close()
			}
			return nil, fmt.Errorf("failed to open City database: %w", err)
		}
	}
    
    if asnDB == nil && cityDB == nil {
        return nil, fmt.Errorf("no GeoLite2 databases found in %s", dbDir)
    }

	return &GeoService{
		asnDB:  asnDB,
		cityDB: cityDB,
	}, nil
}

func (s *GeoService) Close() {
	if s.asnDB != nil {
		s.asnDB.Close()
	}
	if s.cityDB != nil {
		s.cityDB.Close()
	}
}

func (s *GeoService) Lookup(ipStr string) (*Location, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	loc := &Location{
		IP:     ipStr,
		IPType: getIPType(ip),
	}

	if s.cityDB != nil {
		record, err := s.cityDB.City(ip)
		if err == nil {
			loc.City = record.City.Names["en"]
			loc.PostalCode = record.Postal.Code
			if len(record.Subdivisions) > 0 {
				loc.Region = record.Subdivisions[0].Names["en"]
			}
			loc.Country = record.Country.Names["en"]
			loc.CountryCode = record.Country.IsoCode
			loc.Continent = record.Continent.Names["en"]
			loc.ContinentCode = record.Continent.Code
			loc.Timezone = record.Location.TimeZone
			loc.Latitude = record.Location.Latitude
             loc.Longitude = record.Location.Longitude
			loc.AccuracyRadius = record.Location.AccuracyRadius
		}
	}

	if s.asnDB != nil {
		record, err := s.asnDB.ASN(ip)
		if err == nil {
			loc.ASN = record.AutonomousSystemNumber
			loc.Organization = record.AutonomousSystemOrganization
		}
	}

    // Fix for ASN struct field access - typically Reader returns a struct with these fields
    // Ensure we are using the correct fields from geoip2-golang

	loc.NetworkClass = classifyNetwork(loc.Organization)

	return loc, nil
}

func getIPType(ip net.IP) string {
	if ip.IsLoopback() {
		return "Loopback"
	}
	if ip.IsPrivate() {
		return "Private"
	}
	if ip.IsUnspecified() {
		return "Unspecified" // 0.0.0.0
	}
    // Simple check for now. Could be more robust.
	if ip.To4() != nil {
		return "Public IPv4"
	}
	return "Public IPv6"
}

func classifyNetwork(org string) string {
	hostingKeywords := []string{
		"Amazon", "Google", "DigitalOcean", "Hetzner", "OVH", "Azure", "Microsoft", "Alibaba", "Tencent", "Linode", "Vultr", "Oracle",
	}
	orgLower := strings.ToLower(org)
	for _, kw := range hostingKeywords {
		if strings.Contains(orgLower, strings.ToLower(kw)) {
			return "Hosting Provider"
		}
	}
	return "Likely Residential/ISP"
}
