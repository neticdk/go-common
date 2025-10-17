package types

type (
	Regions []Region
	Region  string
)

const (
	NeticRegionDKNorth            Region = "dk-north"
	AzureRegionAustraliacentral   Region = "australiacentral"
	AzureRegionAustraliaeast      Region = "australiaeast"
	AzureRegionAustraliasoutheast Region = "australiasoutheast"
	AzureRegionAustriaeast        Region = "austriaeast"
	AzureRegionBelgiumcentral     Region = "belgiumcentral"
	AzureRegionBrazilsouth        Region = "brazilsouth"
	AzureRegionCanadacentral      Region = "canadacentral"
	AzureRegionCanadaeast         Region = "canadaeast"
	AzureRegionCentralindia       Region = "centralindia"
	AzureRegionCentralus          Region = "centralus"
	AzureRegionChilecentral       Region = "chilecentral"
	AzureRegionChinaeast          Region = "chinaeast"
	AzureRegionChinaeast2         Region = "chinaeast2"
	AzureRegionChinanorth         Region = "chinanorth"
	AzureRegionChinanorth2        Region = "chinanorth2"
	AzureRegionChinanorth3        Region = "chinanorth3"
	AzureRegionDenmarkeast        Region = "denmarkeast"
	AzureRegionEastasia           Region = "eastasia"
	AzureRegionEastus             Region = "eastus"
	AzureRegionEastus2            Region = "eastus2"
	AzureRegionEastus3            Region = "eastus3"
	AzureRegionFinlandcentral     Region = "finlandcentral"
	AzureRegionFrancecentral      Region = "francecentral"
	AzureRegionGermanywestcentral Region = "germanywestcentral"
	AzureRegionGreececentral      Region = "greececentral"
	AzureRegionIndiasouthcentral  Region = "indiasouthcentral"
	AzureRegionIndonesiacentral   Region = "indonesiacentral"
	AzureRegionIsraelcentral      Region = "israelcentral"
	AzureRegionItalynorth         Region = "italynorth"
	AzureRegionJapaneast          Region = "japaneast"
	AzureRegionJapanwest          Region = "japanwest"
	AzureRegionKoreacentral       Region = "koreacentral"
	AzureRegionMalaysiawest       Region = "malaysiawest"
	AzureRegionMexicocentral      Region = "mexicocentral"
	AzureRegionNewzealandnorth    Region = "newzealandnorth"
	AzureRegionNorthcentralus     Region = "northcentralus"
	AzureRegionNortheurope        Region = "northeurope"
	AzureRegionNorwayeast         Region = "norwayeast"
	AzureRegionPolandcentral      Region = "polandcentral"
	AzureRegionQatarcentral       Region = "qatarcentral"
	AzureRegionSaudiarabiacentral Region = "saudiarabiacentral"
	AzureRegionSouthafricanorth   Region = "southafricanorth"
	AzureRegionSouthcentralus     Region = "southcentralus"
	AzureRegionSoutheastasia      Region = "southeastasia"
	AzureRegionSouthindia         Region = "southindia"
	AzureRegionSpaincentral       Region = "spaincentral"
	AzureRegionSwedencentral      Region = "swedencentral"
	AzureRegionSwitzerlandnorth   Region = "switzerlandnorth"
	AzureRegionTaiwannorth        Region = "taiwannorth"
	AzureRegionUaenorth           Region = "uaenorth"
	AzureRegionUksouth            Region = "uksouth"
	AzureRegionUkwest             Region = "ukwest"
	AzureRegionUsdodcentral       Region = "usdodcentral"
	AzureRegionUsdodeast          Region = "usdodeast"
	AzureRegionUsgovarizona       Region = "usgovarizona"
	AzureRegionUsgovtexas         Region = "usgovtexas"
	AzureRegionUsgovvirginia      Region = "usgovvirginia"
	AzureRegionUsseceast          Region = "usseceast"
	AzureRegionUssecwest          Region = "ussecwest"
	AzureRegionUssecwestcentral   Region = "ussecwestcentral"
	AzureRegionWestcentralus      Region = "westcentralus"
	AzureRegionWesteurope         Region = "westeurope"
	AzureRegionWestus             Region = "westus"
	AzureRegionWestus2            Region = "westus2"
	AzureRegionWestus3            Region = "westus3"
	AWSRegionAPEast1              Region = "ap-east-1"
	AWSRegionAPNortheast1         Region = "ap-northeast-1"
	AWSRegionAPNortheast3         Region = "ap-northeast-3"
	AWSRegionAPSouth1             Region = "ap-south-1"
	AWSRegionAPSouth2             Region = "ap-south-2"
	AWSRegionAPSoutheast1         Region = "ap-southeast-1"
	AWSRegionAPSoutheast2         Region = "ap-southeast-2"
	AWSRegionAPSoutheast3         Region = "ap-southeast-3"
	AWSRegionAPSoutheast4         Region = "ap-southeast-4"
	AWSRegionCACentral1           Region = "ca-central-1"
	AWSRegionEUCentral1           Region = "eu-central-1"
	AWSRegionEUCentral2           Region = "eu-central-2"
	AWSRegionEUNorth1             Region = "eu-north-1"
	AWSRegionEUSouth1             Region = "eu-south-1"
	AWSRegionEUSouth2             Region = "eu-south-2"
	AWSRegionEUWest1              Region = "eu-west-1"
	AWSRegionEUWest3              Region = "eu-west-3"
	AWSRegionMESouth1             Region = "me-south-1"
	AWSRegionSAEast1              Region = "sa-east-1"
)

var neticRegions = Regions{
	NeticRegionDKNorth,
}

var awsRegions = Regions{
	AWSRegionAPEast1,
	AWSRegionAPNortheast1,
	AWSRegionAPNortheast3,
	AWSRegionAPSouth1,
	AWSRegionAPSouth2,
	AWSRegionAPSoutheast1,
	AWSRegionAPSoutheast2,
	AWSRegionAPSoutheast3,
	AWSRegionAPSoutheast4,
	AWSRegionCACentral1,
	AWSRegionEUCentral1,
	AWSRegionEUCentral2,
	AWSRegionEUNorth1,
	AWSRegionEUSouth1,
	AWSRegionEUSouth2,
	AWSRegionEUWest1,
	AWSRegionEUWest3,
	AWSRegionMESouth1,
	AWSRegionSAEast1,
}

var azureRegions = Regions{
	AzureRegionAustraliacentral,
	AzureRegionAustraliaeast,
	AzureRegionAustraliasoutheast,
	AzureRegionAustriaeast,
	AzureRegionBelgiumcentral,
	AzureRegionBrazilsouth,
	AzureRegionCanadacentral,
	AzureRegionCanadaeast,
	AzureRegionCentralindia,
	AzureRegionCentralus,
	AzureRegionChilecentral,
	AzureRegionChinaeast,
	AzureRegionChinaeast2,
	AzureRegionChinanorth,
	AzureRegionChinanorth2,
	AzureRegionChinanorth3,
	AzureRegionDenmarkeast,
	AzureRegionEastasia,
	AzureRegionEastus,
	AzureRegionEastus2,
	AzureRegionEastus3,
	AzureRegionFinlandcentral,
	AzureRegionFrancecentral,
	AzureRegionGermanywestcentral,
	AzureRegionGreececentral,
	AzureRegionIndiasouthcentral,
	AzureRegionIndonesiacentral,
	AzureRegionIsraelcentral,
	AzureRegionItalynorth,
	AzureRegionJapaneast,
	AzureRegionJapanwest,
	AzureRegionKoreacentral,
	AzureRegionMalaysiawest,
	AzureRegionMexicocentral,
	AzureRegionNewzealandnorth,
	AzureRegionNorthcentralus,
	AzureRegionNortheurope,
	AzureRegionNorwayeast,
	AzureRegionPolandcentral,
	AzureRegionQatarcentral,
	AzureRegionSaudiarabiacentral,
	AzureRegionSouthafricanorth,
	AzureRegionSouthcentralus,
	AzureRegionSoutheastasia,
	AzureRegionSouthindia,
	AzureRegionSpaincentral,
	AzureRegionSwedencentral,
	AzureRegionSwitzerlandnorth,
	AzureRegionTaiwannorth,
	AzureRegionUaenorth,
	AzureRegionUksouth,
	AzureRegionUkwest,
	AzureRegionUsdodcentral,
	AzureRegionUsdodeast,
	AzureRegionUsgovarizona,
	AzureRegionUsgovtexas,
	AzureRegionUsgovvirginia,
	AzureRegionUsseceast,
	AzureRegionUssecwest,
	AzureRegionUssecwestcentral,
	AzureRegionWestcentralus,
	AzureRegionWesteurope,
	AzureRegionWestus,
	AzureRegionWestus2,
	AzureRegionWestus3,
}

var allRegions = Regions{}

var (
	neticRegionsFieldMap = map[string]Region{}
	azureRegionsFieldMap = map[string]Region{}
	awsRegionsFieldMap   = map[string]Region{}
	allRegionsFieldMap   = map[string]Region{}
)

func (r Region) String() string {
	return string(r)
}

func init() {
	allRegions = append(allRegions, neticRegions...)
	allRegions = append(allRegions, azureRegions...)
	allRegions = append(allRegions, awsRegions...)

	for _, f := range allRegions {
		allRegionsFieldMap[string(f)] = f
	}

	for _, f := range neticRegions {
		neticRegionsFieldMap[string(f)] = f
	}

	for _, f := range azureRegions {
		azureRegionsFieldMap[string(f)] = f
	}

	for _, f := range awsRegions {
		awsRegionsFieldMap[string(f)] = f
	}
}

func AllRegions() Regions {
	return allRegions
}

func AllRegionsString() (regions []string) {
	for _, r := range allRegions {
		regions = append(regions, r.String())
	}
	return regions
}

func PartitionRegions(p Partition) Regions {
	switch p {
	case NeticPartition:
		return neticRegions
	case AzurePartition:
		return azureRegions
	case AWSPartition:
		return awsRegions
	default:
		return Regions{}
	}
}

func ParseRegion(name string) (Region, bool) {
	field, ok := allRegionsFieldMap[name]
	return field, ok
}

func ParseNeticRegion(name string) (Region, bool) {
	field, ok := neticRegionsFieldMap[name]
	return field, ok
}

func ParseAzureRegion(name string) (Region, bool) {
	field, ok := azureRegionsFieldMap[name]
	return field, ok
}

func ParseAWSRegion(name string) (Region, bool) {
	field, ok := awsRegionsFieldMap[name]
	return field, ok
}
