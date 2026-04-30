package return_ticket_service

import (
	"core-ticket/modules/return_ticket/return_ticket_dto"
	"core-ticket/modules/return_ticket/return_ticket_repository"
	"fmt"
	"sort"
	"strings"
)

type ReturnTicketServiceImpl struct {
	return_ticket_repository.ReturnTicketRepository
}

func NewReturnTicketService(
	returnTicketRepository return_ticket_repository.ReturnTicketRepository,
) ReturnTicketService {
	return &ReturnTicketServiceImpl{
		returnTicketRepository,
	}
}

// GetReturnTicketsWithAnalysis analyzes return tickets and groups them by issue type (devCheckNotes)
func (s *ReturnTicketServiceImpl) GetReturnTicketsWithAnalysis(productTypeID int) (*return_ticket_dto.AnalysisResult, error) {
	// Get return tickets from repository
	tickets, err := s.ReturnTicketRepository.GetReturnTicketsWithAnalysis(productTypeID)
	if err != nil {
		return nil, err
	}

	// Group tickets by devCheckNotes (issue type)
	issueMap := s.groupTicketsByIssueType(tickets)

	// Convert to analysis groups
	var issueGroups []return_ticket_dto.IssueAnalysis
	var mostCommon *return_ticket_dto.IssueAnalysis

	for issueType, groupedTickets := range issueMap {
		analysis := return_ticket_dto.IssueAnalysis{
			IssueType:   issueType,
			Count:       len(groupedTickets),
			Tickets:     groupedTickets,
			MainProblem: s.extractMainProblem(groupedTickets),
		}
		issueGroups = append(issueGroups, analysis)

		if mostCommon == nil || analysis.Count > mostCommon.Count {
			mostCommon = &analysis
		}
	}

	// Sort by count (descending)
	sort.Slice(issueGroups, func(i, j int) bool {
		return issueGroups[i].Count > issueGroups[j].Count
	})

	result := &return_ticket_dto.AnalysisResult{
		TotalTickets:    len(tickets),
		IssueGroups:     issueGroups,
		MostCommonIssue: mostCommon,
	}

	return result, nil
}

// GetReturnTicketsByIssueType retrieves return tickets for a specific issue type
func (s *ReturnTicketServiceImpl) GetReturnTicketsByIssueType(productTypeID int, issueType string) ([]return_ticket_dto.ReturnTicket, error) {
	return s.ReturnTicketRepository.GetReturnTicketsByIssueType(productTypeID, issueType)
}

// groupTicketsByIssueType groups tickets by intelligent keyword-based categorization
func (s *ReturnTicketServiceImpl) groupTicketsByIssueType(tickets []return_ticket_dto.ReturnTicket) map[string][]return_ticket_dto.ReturnTicket {
	issueMap := make(map[string][]return_ticket_dto.ReturnTicket)

	for _, ticket := range tickets {
		category := s.categorizeTicketByKeyword(ticket)
		issueMap[category] = append(issueMap[category], ticket)
	}

	return issueMap
}

// categorizeTicketByKeyword analyzes ticket content to assign it to a business domain category
func (s *ReturnTicketServiceImpl) categorizeTicketByKeyword(ticket return_ticket_dto.ReturnTicket) string {
	combinedText := strings.ToLower(strings.TrimSpace(ticket.CheckNotes + " " + ticket.DevCheckNotes))

	// Define category keywords with priority (more specific categories first)
	categories := []struct {
		name     string
		keywords []string
	}{
		// Highly specific categories first
		{
			"Stock Opname",
			[]string{"opname", "stock opname", "sp202"},
		},
		{
			"POS Data Upload",
			[]string{"pos data upload", "pos upload"},
		},
		{
			"Simple Transfer (STF)",
			[]string{"stf202", "simple transfer", "transferan"},
		},
		{
			"Production Result (PL)",
			[]string{"pl202", "production result", "hasil produksi"},
		},
		{
			"Production Return (PN)",
			[]string{"pn202", "production return", "retur produksi"},
		},
		{
			"Link Company (PO-SO)",
			[]string{"po-so", "po so", "link company", "supplier po", "po refer"},
		},
		{
			"Purchase Request (PR)",
			[]string{"pr202", "purchase request", "permintaan pembelian"},
		},
		{
			"Purchase Order (PO)",
			[]string{"po202", "purchase order", "pesanan pembelian"},
		},
		{
			"Charts of Account (COA)",
			[]string{"coa", "chart of account", "rekening"},
		},
		{
			"Data Erasure/Deletion",
			[]string{"penghapusan", "dihapus", "delete", "hapus", "menghapus"},
		},
		{
			"Advance Payment",
			[]string{"advance", "uang muka", "advance payment", "down payment"},
		},
		{
			"Printing",
			[]string{"print", "cetak", "printer", "printing"},
		},
		{
			"Asset Management",
			[]string{"asset", "aset", "fixed asset", "aktiva tetap"},
		},
		{
			"Aktiva Pasiva (Balance Sheet)",
			[]string{"aktiva", "pasiva", "balance sheet", "neraca", "aktiva pasiva"},
		},
		// General categories
		{
			"Goods Movement (GR/GD)",
			[]string{"gr202", "gd202", "goods receipt", "goods delivery", "barang masuk", "barang keluar"},
		},
		{
			"Stock & Inventory",
			[]string{"stock", "stok", "inventory", "stock card", "stock list", "minus stok", "variance"},
		},
		{
			"POS Transactions",
			[]string{"pos", "sbk", "pst202", "pos sales", "transaction", "receipt", "pointofsale"},
		},
		{
			"Product Sales (ERP)",
			[]string{"product sales", "sales actuation", "sales recap", "penjualan produk", "sales report"},
		},
		{
			"Payment & Banking",
			[]string{"payment", "bank", "cash", "bank reconcile", "bank ledger", "payment method", "bank fee", "settlement"},
		},
		{
			"Accounting & GL",
			[]string{"gl", "general ledger", "reconcil", "balance", "trial balance", "journal", "jurnal", "akun", "posting"},
		},
		{
			"Master Data",
			[]string{"master", "product", "bom", "menu", "category", "branch", "barcode", "sku", "menu code"},
		},
		{
			"Upload & Integration",
			[]string{"upload", "sync", "push", "pull", "payload", "import", "template", "format", "excel"},
		},
		{
			"Manufacturing & Production",
			[]string{"manufacturing", "manuf", "production", "batch", "simple manufacturing"},
		},
		{
			"Reports",
			[]string{"report", "recap", "actuation", "reporting", "display", "tampil"},
		},
		{
			"Transfer & Movement",
			[]string{"transfer", "tf202", "movement", "pergerakan"},
		},
		{
			"User & Permissions",
			[]string{"user", "role", "access", "login", "permission", "superadmin"},
		},
	}

	// Match category by keywords
	for _, cat := range categories {
		for _, keyword := range cat.keywords {
			if strings.Contains(combinedText, keyword) {
				return cat.name
			}
		}
	}

	// Determine status categories
	if strings.Contains(combinedText, "no feedback") || strings.Contains(combinedText, "belum feedback") ||
		strings.Contains(combinedText, "tidak ada feedback") {
		return "No Feedback"
	}

	if strings.Contains(combinedText, "no issue") || strings.Contains(combinedText, "bukan issue") ||
		strings.Contains(combinedText, "tidak issue") || strings.Contains(combinedText, "ga issue") {
		return "Non-Issue (Clarified)"
	}

	if strings.Contains(combinedText, "dev") || strings.Contains(combinedText, "development") ||
		strings.Contains(combinedText, "enhancement") || strings.Contains(combinedText, "sprint") {
		return "Development Required"
	}

	return "Other/Unclassified"
}

// extractMainProblem analyzes ticket descriptions (checkNotes, devCheckNotes, and description) to identify the main problem
// Supports both Indonesian and English text
func (s *ReturnTicketServiceImpl) extractMainProblem(tickets []return_ticket_dto.ReturnTicket) string {
	if len(tickets) == 0 {
		return "No description available"
	}

	// Combined stopwords: English + Indonesian
	stopwords := map[string]bool{
		// English
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "must": true, "can": true, "on": true,
		"in": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "from": true, "by": true, "as": true,
		// Indonesian
		"dan": true, "atau": true, "pada": true, "di": true, "ke": true,
		"yang": true, "ini": true, "itu": true, "ada": true, "tidak": true,
		"sudah": true, "telah": true, "akan": true, "adalah": true,
		"bisa": true, "dapat": true, "harus": true, "jika": true, "karena": true,
		"dari": true, "untuk": true, "dengan": true, "tanpa": true, "saat": true,
		"ketika": true, "setelah": true, "sebelum": true, "oleh": true, "melalui": true,
		// Common in both
		"return": true, "ticket": true, "issue": true, "problem": true,
		"error": true, "fail": true, "gagal": true, "berhasil": true,
	}

	// Analyze all relevant text fields
	wordCount := make(map[string]int)
	problemPhrases := make(map[string]int)

	for _, ticket := range tickets {
		// Combine all text fields for analysis
		combinedText := strings.ToLower(
			ticket.CheckNotes + " " + ticket.DevCheckNotes + " " + ticket.Description,
		)

		// Extract keywords (words longer than 3 chars, not stopwords)
		words := strings.FieldsFunc(combinedText, func(r rune) bool {
			return r == ' ' || r == ',' || r == '.' || r == '!' || r == '?' || r == ';' || r == ':' || r == '\n'
		})

		for _, word := range words {
			word = strings.Trim(word, ".,!?;:-()[]{}\"'")
			if len(word) > 3 && !stopwords[word] && word != "" {
				wordCount[word]++
			}
		}

		// Extract 2-3 word phrases related to issues
		if strings.Contains(combinedText, "error") || strings.Contains(combinedText, "gagal") {
			problemPhrases["Error/Failure"]++
		}
		if strings.Contains(combinedText, "tidak") && strings.Contains(combinedText, "bisa") {
			problemPhrases["Cannot/Unable"]++
		}
		if strings.Contains(combinedText, "tidak") && strings.Contains(combinedText, "terbentuk") {
			problemPhrases["Not Created"]++
		}
		if strings.Contains(combinedText, "selisih") {
			problemPhrases["Discrepancy/Mismatch"]++
		}
		if strings.Contains(combinedText, "minus") {
			problemPhrases["Negative Stock"]++
		}
		if strings.Contains(combinedText, "duplik") {
			problemPhrases["Duplicate"]++
		}
	}

	// Find most common significant words
	var topWords []struct {
		word  string
		count int
	}
	for word, count := range wordCount {
		if count > 1 {
			topWords = append(topWords, struct {
				word  string
				count int
			}{word, count})
		}
	}

	// Sort by frequency
	sort.Slice(topWords, func(i, j int) bool {
		return topWords[i].count > topWords[j].count
	})

	// Build result from top keywords and phrases
	var result strings.Builder
	result.WriteString("Primary issues: ")

	// Include top problem phrases
	if len(problemPhrases) > 0 {
		var phrases []struct {
			phrase string
			count  int
		}
		for phrase, count := range problemPhrases {
			if count > 0 {
				phrases = append(phrases, struct {
					phrase string
					count  int
				}{phrase, count})
			}
		}
		sort.Slice(phrases, func(i, j int) bool {
			return phrases[i].count > phrases[j].count
		})

		for i, p := range phrases {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(p.phrase)
		}
	}

	// Add top keywords
	if len(topWords) > 0 {
		result.WriteString(" | Top keywords: ")
		for i := 0; i < 3 && i < len(topWords); i++ {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(fmt.Sprintf("%s(%d)", topWords[i].word, topWords[i].count))
		}
	}

	if result.String() == "Primary issues: " {
		return "Mixed issues"
	}

	return result.String()
}
