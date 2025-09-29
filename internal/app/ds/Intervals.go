package ds

type Interval struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	IsDelete    bool    `gorm:"type:boolean not null;default:false"`
	Photo       string  `gorm:"type:varchar(100)"`
	Title       string  `gorm:"type:varchar(255) not null"`
	Description string  `gorm:"type:varchar(255) not null"`
	Tone        float64 `gorm:"type:numeric(10,1)"`
}
