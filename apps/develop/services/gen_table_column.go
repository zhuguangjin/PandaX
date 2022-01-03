package services

import (
	"errors"
	"pandax/apps/develop/entity"
	"pandax/base/biz"
	"pandax/base/config"
	"pandax/base/global"
)

/**
 * @Description
 * @Author Panda
 * @Date 2021/12/31 14:44
 **/

type (
	SysGenTableColumnModel interface {
		FindDbTablesColumnListPage(page, pageSize int, data entity.DBColumns) (*[]entity.DBColumns, int64)
		FindDbTableColumnList(tableName string) *[]entity.DBColumns

		Insert(data entity.DevGenTableColumn) *entity.DevGenTableColumn
		FindList(data entity.DevGenTableColumn, exclude bool) *[]entity.DevGenTableColumn
		Update(data entity.DevGenTableColumn) *entity.DevGenTableColumn
	}

	devTableColumnModelImpl struct {
		table string
	}
)

var DevTableColumnModelDao SysGenTableColumnModel = &devTableColumnModelImpl{
	table: "dev_gen_table_columns",
}

func (m *devTableColumnModelImpl) FindDbTablesColumnListPage(page, pageSize int, data entity.DBColumns) (*[]entity.DBColumns, int64) {
	list := make([]entity.DBColumns, 0)
	var total int64 = 0
	offset := pageSize * (page - 1)
	if config.Conf.Server.DbType != "mysql" && config.Conf.Server.DbType == "postgresql" {
		biz.ErrIsNil(errors.New("只支持mysql和postgresql数据库"), "只支持mysql和postgresql数据库")
	}

	db := global.Db.Table("information_schema.COLUMNS")
	db = db.Where("table_schema= ? ", config.Conf.Gen.Dbname)

	if data.TableName != "" {
		db = db.Where("table_name = ?", data.TableName)
	}

	err := db.Count(&total).Error
	err = db.Limit(pageSize).Offset(offset).Find(&list).Error
	biz.ErrIsNil(err, "查询生成代码列表信息失败")
	return &list, total
}

func (m *devTableColumnModelImpl) FindDbTableColumnList(tableName string) *[]entity.DBColumns {
	resData := make([]entity.DBColumns, 0)
	if config.Conf.Server.DbType != "mysql" && config.Conf.Server.DbType == "postgresql" {
		biz.ErrIsNil(errors.New("只支持mysql和postgresql数据库"), "只支持mysql和postgresql数据库")
	}
	db := global.Db.Table("information_schema.columns")
	db = db.Where("table_schema = ? ", config.Conf.Gen.Dbname)
	biz.IsTrue(tableName != "", "table name cannot be empty！")

	db = db.Where("table_name = ?", tableName)
	err := db.First(&resData).Error
	biz.ErrIsNil(err, err.Error())
	return &resData
}

func (m *devTableColumnModelImpl) Insert(dgt entity.DevGenTableColumn) *entity.DevGenTableColumn {
	err := global.Db.Table(m.table).Create(&dgt).Error
	biz.ErrIsNil(err, "新增生成代码字段表失败")
	return &dgt
}

func (m *devTableColumnModelImpl) FindList(data entity.DevGenTableColumn, exclude bool) *[]entity.DevGenTableColumn {
	list := make([]entity.DevGenTableColumn, 0)
	db := global.Db.Table(m.table).Where("table_id = ?", data.TableId)
	if exclude {
		notIn := make([]string, 6)
		notIn = append(notIn, "id")
		notIn = append(notIn, "create_by")
		notIn = append(notIn, "update_by")
		notIn = append(notIn, "created_at")
		notIn = append(notIn, "updated_at")
		notIn = append(notIn, "deleted_at")
		db = db.Where(" column_name not in(?)", notIn)
	}
	err := db.Find(&list).Error
	biz.ErrIsNil(err, "查询生成代码字段表信息失败")
	return &list
}

func (m *devTableColumnModelImpl) Update(data entity.DevGenTableColumn) *entity.DevGenTableColumn {
	err := global.Db.Table(m.table).Model(&data).Updates(&data).Error
	biz.ErrIsNil(err, "修改生成代码字段表失败")
	return &data
}
