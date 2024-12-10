package spinner

import "github.com/chelnak/ysmrr"

func WithSpinner(s ysmrr.SpinnerManager, message string, f func() error) error {
	if s == nil {
		return f()
	}
	spinner := s.AddSpinner(message)
	spinner.UpdatePrefix("  ")
	err := f()
	if err != nil {
		spinner.Error()
		return err
	}
	spinner.Complete()
	return nil
}
