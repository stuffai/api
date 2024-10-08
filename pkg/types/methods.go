package types

type Signable interface {
	GetBucket() Bucket
	SetURL(string)
}

func (i *Image) GetBucket() Bucket {
	return i.Bucket
}

func (i *Image) SetURL(url string) {
	i.URL = url
}

func (l ImageList) SignableUserProfiles() []Signable {
	out := make([]Signable, len(l))
	for i, v := range l {
		out[i] = v.User
	}
	return out
}

func (u *UserProfile) GetBucket() Bucket {
	return u.PPBucket
}

func (u *UserProfile) SetURL(url string) {
	u.PPURL = url
}

func (u *ImageUser) GetBucket() Bucket {
	return u.PPBucket
}

func (u *ImageUser) SetURL(url string) {
	u.PPURL = url
}

func (u *LeaderboardEntry) GetBucket() Bucket {
	return u.PPBucket
}

func (u *LeaderboardEntry) SetURL(url string) {
	u.PPURL = url
}

func (m SignableMap) GetBucket() Bucket {
	bucket, ok := m["bucket"].(SignableMap)
	if !ok {
		return Bucket{}
	}
	return Bucket{
		Name: bucket["name"].(string),
		Key:  bucket["key"].(string),
	}
}

func (m SignableMap) SetURL(url string) {
	m["imgURL"] = url
	delete(m, "bucket")
}

func (n *Notification) GetBucket() Bucket {
	return n.Data.GetBucket()
}

func (n *Notification) SetURL(url string) {
	n.Data.SetURL(url)
}

func (c *Comment) GetBucket() Bucket {
	return c.PPBucket
}

func (c *Comment) SetURL(url string) {
	c.PPURL = url
}
