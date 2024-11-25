package redis

type redisKey string

const (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	CONFIG_ENV_KEY redisKey = "CONFIG_ENV"
	USER_INFO_KEY  redisKey = "USER_INFO"
	OTP_KEY        redisKey = "OTP"
	PET_INFO_KEY   redisKey = "PET_INFO"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	CONFIG_ENV_KEY      redisKey = "CONFIG_ENV"
	USER_INFO_KEY       redisKey = "USER_INFO"
	TOKEN_USER_INFO_KEY redisKey = "TOKEN_USER_INFO"
>>>>>>> dff4498 (calendar api)
=======
	CONFIG_ENV_KEY redisKey = "CONFIG_ENV"
	USER_INFO_KEY  redisKey = "USER_INFO"
	OTP_KEY        redisKey = "OTP"
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
>>>>>>> 98e9e45 (ratelimit and recovery function)
=======
	CONFIG_ENV_KEY      redisKey = "CONFIG_ENV"
	USER_INFO_KEY       redisKey = "USER_INFO"
	TOKEN_USER_INFO_KEY redisKey = "TOKEN_USER_INFO"
>>>>>>> dff4498 (calendar api)
=======
	CONFIG_ENV_KEY redisKey = "CONFIG_ENV"
	USER_INFO_KEY  redisKey = "USER_INFO"
	OTP_KEY        redisKey = "OTP"
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
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
