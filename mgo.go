package goapis

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DAO 数据操作
type DAO struct {
	Session *mgo.Session
}

// InsertDo 插入操作
func (db *DAO) InsertDo(dbName, collName string, data bson.M) (bson.ObjectId, error) {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// new id
	newID := bson.NewObjectId()
	data["_id"] = newID
	// insert data
	err := c.Insert(data)
	return newID, err
}

// InsertDoID 插入操作，带有新的ID
func (db *DAO) InsertDoID(dbName, collName string, id bson.ObjectId, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// new id
	data["_id"] = id
	// insert data
	err := c.Insert(data)
	return err
}

// UpdateDo 更新操作
func (db *DAO) UpdateDo(dbName, collName string, selector, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert data
	err := c.Update(selector, bson.M{"$set": data})
	return err
}

// UpdateAllDo 更新全部操作
func (db *DAO) UpdateAllDo(dbName, collName string, selector, data bson.M) (int, error) {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert data
	info, err := c.UpdateAll(selector, bson.M{"$set": data})

	return info.Matched, err
}

// UpsertByID 更新插入操作
func (db *DAO) UpsertByID(dbName, collName string, id interface{}, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert or update data
	_, err := c.UpsertId(id, bson.M{"$set": data})

	return err
}

// UpsertDo 更新插入操作
func (db *DAO) UpsertDo(dbName, collName string, selector, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert or update data
	_, err := c.Upsert(selector, bson.M{"$set": data})

	return err
}

// UpdateNoSet 更新插入操作
func (db *DAO) UpdateNoSet(dbName, collName string, selector, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert or update data
	err := c.Update(selector, data)

	return err
}

// UpdateAllNoSet 更新全部插入操作
func (db *DAO) UpdateAllNoSet(dbName, collName string, selector, data bson.M) (int, error) {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert or update data
	info, err := c.UpdateAll(selector, data)

	return info.Updated, err
}

// IncDo 增量操作
func (db *DAO) IncDo(dbName, collName string, selector, data bson.M) error {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// insert data
	err := c.Update(selector, bson.M{"$inc": data})
	return err
}

// FindAll 查询操作
func (db *DAO) FindAll(dbName, collName string, query bson.M, skip, limit int, sort []string) []bson.M {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	var ret []bson.M
	c.Find(query).Skip(skip).Limit(limit).Sort(sort...).All(&ret)
	return ret
}

// FindAllSelector 查询操作
func (db *DAO) FindAllSelector(dbName, collName string, query, selector bson.M, skip, limit int, sort []string) []bson.M {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	var ret []bson.M
	c.Find(query).Select(selector).Skip(skip).Limit(limit).Sort(sort...).All(&ret)
	return ret
}

// FindOne 查询一条数据
func (db *DAO) FindOne(dbName, collName string, selector, query bson.M) bson.M {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	var ret bson.M
	c.Find(query).Select(selector).One(&ret)
	return ret
}

// RowsCount 查询记录数
func (db *DAO) RowsCount(dbName, collName string, query bson.M) int {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// query
	n, err := c.Find(query).Count()
	if err != nil {
		return 0
	}
	return n
}

// FindByGroup 聚合查询
func (db *DAO) FindByGroup(dbName, collName string, match, group, sort bson.M, skip, limit int) []bson.M {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	// pip
	pipes := []bson.M{
		bson.M{"$match": match},
		bson.M{"$group": group},
	}
	// 排序和数量
	if sort != nil {
		pipes = append(pipes, bson.M{"$sort": sort})
	}
	if skip > 0 {
		pipes = append(pipes, bson.M{"$skip": skip})
	}
	if limit > 0 {
		pipes = append(pipes, bson.M{"$limit": limit})
	}
	pipe := c.Pipe(pipes)
	iter := pipe.Iter()
	var list []bson.M
	iter.All(&list)
	iter.Close()
	return list
}

// RemoveAll 删除数据
func (db *DAO) RemoveAll(dbName, collName string, selector bson.M) {
	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(dbName).C(collName)
	c.RemoveAll(selector)
}
