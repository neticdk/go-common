package types

type (
	Partition  string
	Partitions []Partition
)

const (
	NeticPartition Partition = "netic"
	AzurePartition Partition = "azure"
	AWSPartition   Partition = "aws"
)

func (p Partition) String() string {
	return string(p)
}

func AllPartitions() Partitions {
	return Partitions{NeticPartition, AzurePartition, AWSPartition}
}

func AllPartitionsString() (partitions []string) {
	for _, p := range AllPartitions() {
		partitions = append(partitions, p.String())
	}
	return partitions
}

func ParsePartition(name string) (Partition, bool) {
	switch name {
	case "netic":
		return NeticPartition, true
	case "azure":
		return AzurePartition, true
	case "aws":
		return AWSPartition, true
	default:
		return "", false
	}
}

func HasRegion(partition Partition, region Region) bool {
	regions := PartitionRegions(partition)
	for _, r := range regions {
		if r == region {
			return true
		}
	}
	return false
}
