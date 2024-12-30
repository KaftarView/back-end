package enums

type BucketType uint

const (
	EventsBucket BucketType = iota + 1
	PodcastsBucket
	NewsBucket
	JournalsBucket
	ProfilesBucket
)

func (bt BucketType) String() string {
	switch bt {
	case EventsBucket:
		return "EventsBucket"
	case PodcastsBucket:
		return "PodcastsBucket"
	case NewsBucket:
		return "NewsBucket"
	case JournalsBucket:
		return "JournalsBucket"
	case ProfilesBucket:
		return "ProfilesBucket"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		EventsBucket,
		PodcastsBucket,
		NewsBucket,
		JournalsBucket,
		ProfilesBucket,
	}
}
