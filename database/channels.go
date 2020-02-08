package database

import (
	"github.com/TicketsBot/GoPanel/database"
)

type Channel struct {
	ChannelId int64  `gorm:"column:CHANNELID"`
	GuildId   int64  `gorm:"column:GUILDID"`
	Name      string `gorm:"column:NAME;type:VARCHAR(32)"`
	Type      int    `gorm:"column:CHANNELTYPE;type:TINYINT(1)"`
}

func (Channel) TableName() string {
	return "Channel"
}

func StoreChannel(channelId, guildId int64, name string, channelType int) {
	channel := Channel{
		ChannelId: channelId,
		GuildId:   guildId,
		Name:      name,
		Type:      channelType,
	}

	database.Database.Where(&Channel{ChannelId:channelId}).Assign(&channel).FirstOrCreate(&Channel{})
}

func DeleteChannel(channelId int64) {
	var node Channel
	Db.Where(Channel{ChannelId: channelId}).Take(&node)
	Db.Delete(&node)
}