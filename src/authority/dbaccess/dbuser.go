package dbaccess

import (
	"fmt"
	"database/sql"
//	"github.com/jmoiron/sqlx"
)

type User struct {
	Id	uint64	`db:"id" json:"id,omitempty"`
	Code string	`db:"code" json:"code,omitempty"`
	Name	string	`db:"name" json:"name,omitempty"`
}

type UserAuthority struct {
	Id uint64	`db:"id" json:"id,omitempty"`
	UserId uint64	`db:"user_id" json:"user_id,omitempty"`
	AuthorityId	uint64	`db:"authority_id" json:authority_id,omitempty"`
}

func QueryUserCount() (uint64, error) {
	sqlstr := "SELECT COUNT(*) FROM tb_user;"
	var count uint64
	err := db.QueryRow(sqlstr).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func QueryUsers(index int, size int) ([]User, error) {
	sqlstr := "SELECT * FROM tb_user LIMIT ?,?;"
	users := []User{}
	err := db.Select(&users, sqlstr, index, size)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func QueryUserById(id uint64) (*User, error) {
	sqlstr := "SELECT * FROM tb_user WHERE id=?"
	user := &User{}
	err := db.Get(user, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func QueryUserByCode(code string) (*User, error) {
	sqlstr := "SELECT * FROM tb_user WHERE code=?"
	user := &User{}
	err := db.Get(user, sqlstr, code)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUserByField(code string, name string) error {
	sqlstr := "INSERT INTO tb_user(code,name) VALUES(?,?)"
	_, err := db.Exec(sqlstr, code, name)
	return err
}

func CreateUser(user *User) error {
	return CreateUserByField(user.Code, user.Name);
}

func UpdateUser(user *User) error {
	//sqlstr := "UPDATE tb_user "
	return nil
}

func DeleteUserById(id uint64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM tb_user WHERE id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM tb_user_authority WHERE user_id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteUserByCode(code string) error {
	sqlstr := "SELECT id FROM tb_user WHERE code=?"
	var id uint64
	err := db.QueryRow(sqlstr, code).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return DeleteUserById(id)
}

func QueryUserAuthority(user_id uint64) ([]Authority, error) {
	sqlstr := "SELECT t2.* FROM tb_user_authority t1, tb_authority t2 WHERE t1.user_id=? AND t1.authority_id=t2.id"
	auths := []Authority{}
	err := db.Select(&auths, sqlstr, user_id)
	if err != nil {
		return nil, err
	}
	return auths, nil
}

//递归过程
//TODO 多人同时处理，并发同步问题
func TxGrantUserAuthority(user_id, auth_id uint64, tx *sql.Tx) error {
	var err error
//	if tx == nil {
//		tx, err = db.Begin()
//		if err != nil {
//			return err
//		}
//	}
	err = TxDeleteUserSubAuthority(user_id, auth_id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	err = TxCreateUserAuthority(user_id, auth_id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	auth, err := QueryAuthorityById(auth_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if auth.GroupId == 0 {
		//return tx.Commit()
		return nil
	}
	
	//query user auth
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM tb_user_authority t1,tb_authority t2 
		WHERE t1.authority_id=t2.id AND t1.user_id=? AND t2.group_id=?`, user_id, auth.GroupId).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}
	auth_count, err := QueryAuthoritysByGroupId(auth.GroupId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(auth_count) != count {
		//return tx.Commit()
		return nil
	} else {
		return TxGrantUserAuthority(user_id, auth.GroupId, tx)
	}
}

func TxCreateUserAuthority(user_id, auth_id uint64, tx *sql.Tx) error {
	sqlstr := `INSERT INTO tb_user_authority(user_id, authority_id) VALUES(?,?)`
	var err error
	if tx == nil {
		_, err = db.Exec(sqlstr, user_id, auth_id)
	} else {
		_, err = tx.Exec(sqlstr, user_id, auth_id)
	}
	//if err != nil {
	//	fmt.Printf("error:TxCreateUserAuthority")
	//}
	return err
}

func TxDeleteUserAuthorityGroup(user_id uint64, auth_code string, tx *sql.Tx) error {
	sqlstr := fmt.Sprintf(`DELETE FROM tb_user_authority 
		WHERE user_id=? AND authority_id IN 
		(SELECT id FROM tb_authority WHERE code LIKE "%s%%")`, auth_code)
	var err error
	if tx != nil {
		_, err = tx.Exec(sqlstr, user_id)
	} else {
		_, err = db.Exec(sqlstr, user_id)
	}
	//if err != nil {
	//	fmt.Printf("error:TxDeleteUserAuthorityGroup")
	//}
	return err
}

func TxDeleteUserSubAuthority(user_id, auth_id uint64, tx *sql.Tx) error {
	var auth_code string
	err := tx.QueryRow("SELECT code FROM tb_authority WHERE id=?", auth_id).Scan(&auth_code)
	if err != nil {
		return err
	}
	//if err != nil {
	//	fmt.Printf("error:TxDeleteUserSubAuthority")
	//}
	return TxDeleteUserAuthorityGroup(user_id, auth_code+".", tx)
}

func TxDeleteUserAuthority(user_id, auth_id uint64, tx *sql.Tx) error {
	var auth_code string
	err := tx.QueryRow("SELECT code FROM tb_authority WHERE id=?", auth_id).Scan(&auth_code)
	if err != nil {
		return err
	}
	//if err != nil {
	//	fmt.Printf("error:TxDeleteUserAuthority")
	//}
	return TxDeleteUserAuthorityGroup(user_id, auth_code, tx)
}