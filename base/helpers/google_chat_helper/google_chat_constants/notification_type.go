package google_chat_constants

type NotificationType string

const (
	// TypeInternalServerError is a notification type when internal server error occurred.
	//
	// Available payload:
	//
	// - transactionDate: string
	//
	// - title: string
	//
	// - errorMessage: string
	TypeInternalServerError NotificationType = "Internal Server Error"

	// TypeFailedConsumeDisbursementData is a notification type when there is a failed
	// attempt to consume disbursement data.
	//
	// Available payload:
	//
	// - transactionDate: string
	//
	// - companyCode: string
	//
	// - disburseNum: string
	//
	// - title: string
	//
	// - description: string
	TypeFailedConsumeDisbursementData NotificationType = "Failed to Consume Disbursement Data"
)

func (nt NotificationType) String() string { return string(nt) }

func (nt NotificationType) IsValid() bool {
	switch nt {
	case TypeInternalServerError, TypeFailedConsumeDisbursementData:
		return true
	default:
		return false
	}
}
