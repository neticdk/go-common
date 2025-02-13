package errors

type GeneralError struct {
	Message string
	HelpMsg string
	CodeVal int
	Err     error
}

func (e *GeneralError) Error() string { return e.Message }
func (e *GeneralError) Help() string  { return e.HelpMsg }
func (e *GeneralError) Unwrap() error { return e.Err }
func (e *GeneralError) Code() int     { return e.CodeVal }
