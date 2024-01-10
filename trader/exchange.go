package trader

type Exchange interface {
	GetName() (name string)
	Subscribe(params []map[string]string) (err error)
}
