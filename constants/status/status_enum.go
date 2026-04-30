package status

type Status int

const (
	New                Status = 1
	Rejected           Status = 2
	Authorized         Status = 3
	Receiving          Status = 4
	Received           Status = 5
	Delivering         Status = 6
	Delivered          Status = 7
	Finished           Status = 8
	HalfPaid           Status = 9
	FullPaid           Status = 10
	Invoice            Status = 11
	Cancelled          Status = 12
	Preparing          Status = 13
	Served             Status = 14
	Billing            Status = 15
	InProgress         Status = 16
	ResultSignOff      Status = 17
	ProductionFinished Status = 18
	PrintCancelled     Status = 19
	Active             Status = 20
	NotActive          Status = 21
	Redeem             Status = 22
	Expired            Status = 23
	Void               Status = 24
	Closed             Status = 25
	PendingRelease     Status = 26
	Released           Status = 27
	Available          Status = 28
	PartiallyCompleted Status = 30
	Completed          Status = 31
	PartiallyReleased  Status = 32
	Pending            Status = 33
	Prepared           Status = 34
	Processed          Status = 35
	ProcessRefund      Status = 36
	Refunded           Status = 37
	WaitingForApproval Status = 38
	FullReturned       Status = 39
	FinishByOutlet     Status = 40
	Draft              Status = 41
	ChangeBook         Status = 42
	Confirmed          Status = 43
	PartiallyReturned  Status = 45
	Hold               Status = 46
	Used               Status = 47
	NotUsed            Status = 48
	Sold               Status = 49
	Dispose            Status = 50
	Picking            Status = 51
	Picked             Status = 52
	Approved           Status = 53
)

func (s Status) String() string {
	return map[Status]string{
		New:                "New",
		Rejected:           "Rejected",
		Authorized:         "Authorized",
		Receiving:          "Receiving",
		Received:           "Received",
		Delivering:         "Delivering",
		Delivered:          "Delivered",
		Finished:           "Finished",
		HalfPaid:           "Half Paid",
		FullPaid:           "Full Paid",
		Invoice:            "Invoice",
		Cancelled:          "Cancelled",
		Preparing:          "Preparing",
		Served:             "Served",
		Billing:            "Billing",
		InProgress:         "In Progress",
		ResultSignOff:      "Result Sign Off",
		ProductionFinished: "Production Finished",
		PrintCancelled:     "Print Cancelled",
		Active:             "Active",
		NotActive:          "Not Active",
		Redeem:             "Redeem",
		Expired:            "Expired",
		Void:               "Void",
		Closed:             "Closed",
		PendingRelease:     "Pending Release",
		Released:           "Released",
		Available:          "Available",
		PartiallyCompleted: "Partially Completed",
		Completed:          "Completed",
		PartiallyReleased:  "Partially Released",
		Pending:            "Pending",
		Prepared:           "Prepared",
		Processed:          "Processed",
		ProcessRefund:      "Process Refund",
		Refunded:           "Refunded",
		WaitingForApproval: "Waiting For Approval",
		FullReturned:       "Full Returned",
		FinishByOutlet:     "Finish by Outlet",
		Draft:              "Draft",
		ChangeBook:         "Change Book",
		Confirmed:          "Confirmed",
		PartiallyReturned:  "Partially Returned",
		Hold:               "Hold",
		Used:               "Used",
		NotUsed:            "Not Used",
		Sold:               "Sold",
		Dispose:            "Dispose",
		Picking:            "Picking",
		Picked:             "Picked",
	}[s]
}

func (s Status) Int() int {
	return int(s)
}
