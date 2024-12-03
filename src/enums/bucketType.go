package enums

type BucketType uint

const (
	BannersBucket BucketType = iota + 1
	SessionsBucket
	PodcastsBucket
	ProfileBucket
)

// possibly do not need to this at all maybe
func (bt BucketType) String() string {
	switch bt {
	case BannersBucket:
		return "BannersBucket"
	case SessionsBucket:
		return "SessionsBucket"
	case PodcastsBucket:
		return "PodcastsBucket"
	case ProfileBucket:
		return "ProfileBucket"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		BannersBucket,
		SessionsBucket,
		PodcastsBucket,
		ProfileBucket,
	}
}
