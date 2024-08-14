package sources

type ListChapter struct {
	Order   int
	Title   string
	Url     string
	Content string
}

type NovelInfo struct {
	Title    string
	Image    string
	Author   string
	Synopsis string
	Data     []ListChapter
}
