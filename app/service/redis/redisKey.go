package redis

type redisKey string

const (
	CONFIG_ENV_KEY redisKey = "CONFIG_ENV"
	USER_INFO_KEY  redisKey = "USER_INFO"
	OTP_KEY        redisKey = "OTP"
<<<<<<< HEAD
<<<<<<< HEAD
	PET_INFO_KEY   redisKey = "PET_INFO"
=======
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
	PET_INFO_KEY   redisKey = "PET_INFO"
>>>>>>> 98e9e45 (ratelimit and recovery function)
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
