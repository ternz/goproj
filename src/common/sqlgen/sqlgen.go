package sqlgen

type Model struct {
	_tb_name string `model`
	Id       uint64 `field:"id"`
	Name     string `field:"name"`
}

type SqlgenT struct {
	sqlstr string
}

var Sqlgen SqlgenT

func (s *SqlgenT) Select() {

}

func (s *SqlgenT) Insert() {

}

func (s *SqlgenT) Update() {

}

func (s *SqlgenT) Where() {

}

func (s *SqlgenT) Orderby() {

}

func (s *SqlgenT) Groupby() {

}
