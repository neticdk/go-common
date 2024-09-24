package types

type (
	ResilienceZones []ResilienceZone
	ResilienceZone  string
)

const (
	PlatformResilienceZone       ResilienceZone = "platform"
	Internal1ResilienceZone      ResilienceZone = "internal-1"
	Innovators1ResilienceZone    ResilienceZone = "innovators-1"
	EarlyAdopters1ResilienceZone ResilienceZone = "early-adopters-1"
	EarlyMajority1ResilienceZone ResilienceZone = "early-majority-1"
	EarlyMajority2ResilienceZone ResilienceZone = "early-majority-2"
	LateMajority1ResilienceZone  ResilienceZone = "late-majority-1"
	LateMajority2ResilienceZone  ResilienceZone = "late-majority-2"
	Laggards1ResilienceZone      ResilienceZone = "laggards-1"
)

var allResilienceZones = ResilienceZones{
	PlatformResilienceZone,
	Internal1ResilienceZone,
	Innovators1ResilienceZone,
	EarlyAdopters1ResilienceZone,
	EarlyMajority1ResilienceZone,
	EarlyMajority2ResilienceZone,
	LateMajority1ResilienceZone,
	LateMajority2ResilienceZone,
	Laggards1ResilienceZone,
}

var resilienceZoneFieldMap = map[string]ResilienceZone{}

func (r ResilienceZone) String() string {
	return string(r)
}

func init() {
	for _, rz := range allResilienceZones {
		resilienceZoneFieldMap[string(rz)] = rz
	}
}

func AllResilienceZones() ResilienceZones {
	return allResilienceZones
}

func AllResilienceZonesString() (rzs []string) {
	for _, rz := range allResilienceZones {
		rzs = append(rzs, rz.String())
	}
	return
}

func ParseResilienceZone(name string) (ResilienceZone, bool) {
	field, ok := resilienceZoneFieldMap[name]
	return field, ok
}
