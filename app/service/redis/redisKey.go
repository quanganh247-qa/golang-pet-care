package redis

type redisKey string

const (
	CONFIG_ENV_KEY redisKey = "CONFIG_ENV"
	USER_INFO_KEY  redisKey = "USER_INFO"
	OTP_KEY        redisKey = "OTP"
)

// type keyType struct {
// 	PageTransKey redisKey
// 	PageInfoKey  redisKey
// }

// func initKey() *keyType {
// 	return &keyType{
// 		PageTransKey: PAGE_TRANS_KEY,
// 	}
// }

// func (key *keyType) GetPageTransKey(page int64, pageSize int64, sortField string, sortOrder string) string {
// 	return fmt.Sprintf("%s:{%d_%d_%s_%s}", key.PageTransKey, page, pageSize, sortField, sortOrder)
// }
