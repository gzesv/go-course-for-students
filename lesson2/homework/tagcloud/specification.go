package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tagMap  map[string]int
	tagStat []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	tg := new(TagCloud)
	tg.tagMap = make(map[string]int)
	return *tg
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (tg *TagCloud) AddTag(tag string) {
	c := tg.tagMap[tag]
	if len(tg.tagStat) == 0 || tg.tagStat[c].Tag != tag {
		ts := TagStat{tag, 1}
		tg.tagStat = append(tg.tagStat, ts)
		tg.tagMap[tag] = len(tg.tagStat) - 1
	} else {
		tg.tagStat[c].OccurrenceCount += 1
		if c > 0 && tg.tagStat[c-1].OccurrenceCount < tg.tagStat[c].OccurrenceCount {
			tg.tagMap[tag], tg.tagMap[tg.tagStat[c-1].Tag] = tg.tagMap[tg.tagStat[c-1].Tag], tg.tagMap[tag]
			tg.tagStat[c-1], tg.tagStat[c] = tg.tagStat[c], tg.tagStat[c-1]
		}
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (tg *TagCloud) TopN(n int) []TagStat {
	var size int
	if n < len(tg.tagStat) {
		size = n
	} else {
		size = len(tg.tagStat)
	}
	return tg.tagStat[:size]
}
