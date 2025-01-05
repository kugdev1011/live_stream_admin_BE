package cache

const (
	VIDEO_ENCODING_PREFIX = "video:encoding:%s" // expect: videoname. example: key = fmt.Sprintf(cachekeys.VIDEO_ENCODING_PREFIX, "2f2ccfc3-3bbb-45a6-9879-0c574576dda6.mp4"), value = boolean(true)
	IS_ENDING_LIVE_PREFIX = "stream:ending:%d"  // expect stream id. example : key = fmt.Sprintf(cachekeys.IS_ENDING_LIVE_PREFIX, "1"), value = boolean(true)
)
