package myerror

type Error interface {
	getCode() int
	getMessage() string
	getClause() string
}

type DBError struct {
	Code int
	Message string
	Clause string
}

func (dbe *DBError) Error() string {
	return dbe.Message
}

func (dbe *DBError) getCode() int{
	return dbe.Code
}

func (dbe *DBError) getMessage() string{
	return dbe.Message
}

func (dbe *DBError) getClause() string {
	return dbe.Clause
}