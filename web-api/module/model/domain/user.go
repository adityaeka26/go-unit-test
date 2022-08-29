package domain

type InsertUser struct {
	Username     string `bson:"username"`
	Password     string `bson:"password"`
	Name         string `bson:"name"`
	MobileNumber string `bson:"mobileNumber"`
}

type User struct {
	Id           string `bson:"_id"`
	Username     string `bson:"username"`
	Password     string `bson:"password"`
	Name         string `bson:"name"`
	MobileNumber string `bson:"mobileNumber"`
}

type UserRedis struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Otp          string `json:"otp"`
	MobileNumber string `json:"mobileNumber"`
}

type RegisterOtpKafka struct {
	Name         string `json:"name"`
	Otp          string `json:"otp"`
	MobileNumber string `json:"mobileNumber"`
}
