package ds

type CompositorInterval struct {
	ComIntID      uint `gorm:"primaryKey;not null;autoIncrement"`
	CompositionID uint `gorm:"primaryKey;not null"`
	IntervalID    uint `gorm:"not null"`
	Amount        uint `gorm:"type:integer"`

	Interval    Interval    `gorm:"foreignKey:IntervalID;references:ID;constraint:OnUpdate:CASCADE"`
	Composition Composition `gorm:"foreignKey:CompositionID;references:ID;constraint:OnUpdate:CASCADE"`
}
