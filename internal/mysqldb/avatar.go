package mysqldb

type Avatar struct {
	ID   int    `gorm:"type:bigint(20) auto_increment;primary_key" json:"id"`
	Data string `gorm:"type:longblob;" json:"data"`
}
