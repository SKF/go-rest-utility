package problems

import "net/http"

type ValidationProblem struct {
	BasicProblem
	Reasons []ValidationReason `json:"reasons,omitempty"`
}

func Validation(reasons ...ValidationReason) ValidationProblem {
	return ValidationProblem{
		BasicProblem: BasicProblem{
			Type:   "/problems/invalid-request",
			Title:  "Your request parameters didn't validate.",
			Status: http.StatusBadRequest,
			Detail: "See the reasons field for more details.",
		},
		Reasons: reasons,
	}
}

func (problem ValidationProblem) TrimEmpty() error {
	if len(problem.Reasons) == 0 {
		return nil
	}

	return problem
}

func (problem *ValidationProblem) Errors() (errors []error) {
	for _, reason := range problem.Reasons {
		errors = append(errors, reason.Cause)
	}

	return errors
}

func (problem *ValidationProblem) Append(reasons ...ValidationReason) {
	problem.Reasons = append(problem.Reasons, reasons...)
}

func (problem *ValidationProblem) AppendWithPrefix(prefix string, reasons ...ValidationReason) {
	for _, subProblem := range reasons {
		if prefix != "" && subProblem.Name != "" {
			subProblem.Name = prefix + "." + subProblem.Name
		} else if prefix != "" {
			subProblem.Name = prefix
		}

		problem.Reasons = append(problem.Reasons, subProblem)
	}
}

func (problem *ValidationProblem) Merge(prefix string, err error) {
	if err == nil {
		return
	}

	if otherProblem, ok := err.(ValidationProblem); !ok {
		problem.AppendWithPrefix(prefix, ValidationReason{
			Reason: err.Error(),
			Cause:  err,
		})
	} else {
		problem.AppendWithPrefix(prefix, otherProblem.Reasons...)
	}
}

type ValidationReason struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Cause  error  `json:"-"`
}

func (err ValidationReason) Error() string {
	if err.Name != "" {
		return err.Name + " " + err.Reason
	}

	return err.Reason
}
