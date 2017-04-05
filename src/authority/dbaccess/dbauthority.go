package dbaccess

import (
	"database/sql"
	"fmt"
)

type Authority struct {
	Id		uint64		`db:"id" json:"id"`
	Code	string 		`db:"code" json:"code"`
	Name	string		`db:"name" json:"name"`
	GroupId	uint64		`db:"group_id" json:"group_id"`
	SubAuth	[]*Authority	`db:"-" json:"subauth,omitempty"`
}

func insertAuthority(a []*Authority, item *Authority) {
	i := 0
	for _, v := range a {
		if item.Code < v.Code {
			i++
		} else {
			break
		}
	}
	a = append(a, nil)
	copy(a[i+1:], a[i:])
	a[i] = item
}

//要求输入的item是有序的
//TODO 防止意外，应当做检查
func authorityTreeAdd(group *Authority, item *Authority) {
	if item.GroupId == group.Id {
		group.SubAuth = append(group.SubAuth, item)
		return
	}
	index := len(group.SubAuth)-1
	authorityTreeAdd(group.SubAuth[index], item)
}

//input list must be sorted
func ConstructAuthorityTree(sorted_list []Authority) []*Authority {
	root := &Authority{Id:0}
	//必须用index
	for i, _ := range sorted_list {
		authorityTreeAdd(root, &sorted_list[i])
	}
	return root.SubAuth
}

func QueryAuthorityById(id uint64) (*Authority, error) {
	sqlstr := "SELECT * FROM tb_authority WHERE id=?"
	auth := &Authority{}
	err := db.Get(auth, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func QueryAuthorityByCode(code string) (*Authority, error) {
	sqlstr := "SELECT * FROM tb_authority WHERE code=?"
	auth := &Authority{}
	err := db.Get(auth, sqlstr, code)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

//output list was sorted
func QueryAuthorityGroup(code string) ([]Authority, error) {
	sqlstr := fmt.Sprintf(`SELECT * FROM tb_authority WHERE code LIKE "%s%%" ORDER BY code`, code)
	//sqlstr := `SELECT * FROM tb_authority WHERE code LIKE ?%% ORDER BY code`
	auths := []Authority{}
	err := db.Select(&auths, sqlstr)
	if err != nil {
		return nil, err
	}
	return auths, nil
}

func QueryAuthoritysByGroupId(id uint64) ([]Authority, error) {
	sqlstr := "SELECT * FROM tb_authority WHERE group_id=?"
	auths := []Authority{}
	err := db.Select(&auths, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return auths, nil
}

func CreateAuthority(auth *Authority) error {
	//TODO 检查前缀是否一致
	sqlstr := "INSERT INTO tb_authority(code, name, group_id) VALUES(:code, :name, :group_id)"
	_, err := db.NamedExec(sqlstr, auth)
	return err
}

func UpdateAuthorityName(id uint64, name string) error {
	sqlstr := "UPDATE tb_authority SET name=? WHERE id=?"
	_, err := db.Exec(sqlstr, name, id)
	return err
}

//func DeleteAuthorityById(id uint64) error {
//	//TODO 删除整个组（分支）
//	tx, err := db.Begin()
//	if err != nil {
//		return err
//	}
//	_, err = tx.Exec("DELETE FROM tb_authority WHERE id=?", id)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//	_, err = tx.Exec("DELETE FROM tb_user_authority WHERE authority_id=?", id)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//	tx.Commit()
//	return nil
//}

func DeleteAuthorityById(id uint64) error {
	sqlstr := "SELECT code FROM tb_authority WHERE id=?"
	var code string
	err := db.QueryRow(sqlstr, id).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return DeleteAuthorityByCode(code)
}

func DeleteAuthorityByCode(code string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf(`DELETE FROM tb_user_authority 
		WHERE authority_id IN (SELECT id FROM	tb_authority WHERE CODE LIKE "%s%%")`, code))
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(fmt.Sprintf(`DELETE FROM tb_authority WHERE code like "%s%%"`, code))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//func DeleteAuthorityByCode(code string) error {
//	sqlstr := "SELECT id FROM tb_authority WHERE code=?"
//	var id uint64
//	err := db.QueryRow(sqlstr, code).Scan(&id)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil
//		}
//		return err
//	}
//	return DeleteAuthorityById(id)
//}

//func DeleteAuthorityGroup(code string) error {
//	sqlstr := fmt.Sprintf("SELECT * FROM tb_authority WHERE code like ?", code)
//	var id uint64
//	err := db.QueryRow(sqlstr, code).Scan(&id)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil
//		}
//		return err
//	}
//	return DeleteAuthorityById(id)
//}