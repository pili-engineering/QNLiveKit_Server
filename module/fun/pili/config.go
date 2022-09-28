package pili

type Config struct {
	Hub string `mapstructure:"hub"`

	RoomTokenExpireS int64  `mapstructure:"room_token_expire_s"`
	PlayBackUrl      string `mapstructure:"playback_url"`
	StreamPattern    string `mapstructure:"stream_pattern"`
	PublishUrl       string `mapstructure:"publish_url"`
	PublishDomain    string `mapstructure:"publish_domain"`
	RtmpPlayUrl      string `mapstructure:"rtmp_play_url"`
	FlvPlayUrl       string `mapstructure:"flv_play_url"`
	HlsPlayUrl       string `mapstructure:"hls_play_url"`
	SecurityType     string `mapstructure:"security_type"`    //expiry, expiry_sk, none
	PublishKey       string `mapstructure:"publish_key"`      //推流key
	PublishExpireS   int64  `mapstructure:"publish_expire_s"` //推流URL 过期时间，单位：秒
}
