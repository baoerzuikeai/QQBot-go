package my_dto

// file_type	int	是	媒体类型：1 图片，2 视频，3 语音，4 文件（暂不开放） 资源格式要求 图片：png/jpg，视频：mp4，语音：silk
// url	string	是	需要发送媒体资源的url
// srv_send_msg	bool	是	设置 true 会直接发送消息到目标端，且会占用主动消息频次
// file_data	string   base64
type PostMedia struct {
	FileType   int    `json:"file_type"`
	Url        string `json:"url,omitempty"`
	SrvSendMsg bool   `json:"srv_send_msg"`
	FileData   string `json:"file_data,omitempty"`
}
