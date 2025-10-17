package types

type (
	InfrastructureProvider  string
	InfrastructureProviders []InfrastructureProvider
)

const (
	NeticInfrastructureProvider InfrastructureProvider = "netic"
	AzureInfrastructureProvider InfrastructureProvider = "azure"
	AWSInfrastructureProvider   InfrastructureProvider = "aws"
)

func (p InfrastructureProvider) String() string {
	return string(p)
}

func AllInfrastructureProviders() InfrastructureProviders {
	return InfrastructureProviders{NeticInfrastructureProvider, AzureInfrastructureProvider, AWSInfrastructureProvider}
}

func AllInfrastructureProvidersString() (providers []string) {
	for _, p := range AllInfrastructureProviders() {
		providers = append(providers, p.String())
	}
	return providers
}

func ParseInfrastructureProvider(name string) (InfrastructureProvider, bool) {
	switch name {
	case "netic":
		return NeticInfrastructureProvider, true
	case "azure":
		return AzureInfrastructureProvider, true
	case "aws":
		return AWSInfrastructureProvider, true
	default:
		return "", false
	}
}
