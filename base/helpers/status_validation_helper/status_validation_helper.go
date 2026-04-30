package status_validation_helper

import (
	"core-ticket/base/helpers/array_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants/error_code"
	"core-ticket/constants/status"
	"fmt"
	"strings"
)

// ValidateStatusForOperation validates if the current status allows the specified operation.
//
// Parameters:
//   - currentStatus: The current status of the transaction
//   - allowedStatuses: Slice of status values permitted for the operation
//   - operation: String describing the operation (e.g., "updated", "deleted", "authorized")
//
// Returns error if validation fails, nil otherwise.
//
// Example:
//
//	err := ValidateStatusForOperation(model.StatusID, []status.Status{status.New, status.Rejected}, "deleted")
//	if err != nil {
//	    return err
//	}
func ValidateStatusForOperation(
	currentStatus status.Status,
	allowedStatuses []status.Status,
	operation string,
) error {
	if !array_helper.InArray(allowedStatuses, currentStatus) {
		statusNames := make([]string, len(allowedStatuses))
		for i, s := range allowedStatuses {
			statusNames[i] = s.String()
		}
		return error_helper.New(
			fmt.Errorf("only transaction with status %s can be %s", strings.Join(statusNames, " or "), operation),
			error_code.ValidationError,
		)
	}
	return nil
}
