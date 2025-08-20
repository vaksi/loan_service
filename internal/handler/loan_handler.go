package handler

import (
	"context"
	"net/http"
	"time"

	"loan_service/internal/domain"
	"loan_service/internal/repository"

	"github.com/gin-gonic/gin"
)

// LoanHandler defines HTTP handlers for loan‑related endpoints. It
// decouples the HTTP layer from the underlying services and focuses
// solely on request parsing, validation and response formatting.
type LoanHandler struct {
	svc LoanUsecase
}

// LoanUsecase abstracts service layer for handler
// to allow mocking in HTTP tests and to decouple layers.
type LoanUsecase interface {
	CreateLoan(ctx context.Context, input domain.Loan) (*domain.Loan, error)
	ApproveLoan(ctx context.Context, loanID, pictureURL, employeeID string, approvalDate time.Time) (*domain.Loan, error)
	InvestInLoan(ctx context.Context, loanID, investorID, investorName, investorEmail string, amount float64) (*domain.Loan, error)
	DisburseLoan(ctx context.Context, loanID, agreementURL, employeeID string, disbursementDate time.Time) (*domain.Loan, error)
	GetLoanByID(ctx context.Context, id string) (*domain.Loan, error)
	ListLoans(ctx context.Context) ([]domain.Loan, error)
}

// NewLoanHandler constructs a new LoanHandler.
func NewLoanHandler(svc LoanUsecase) *LoanHandler { return &LoanHandler{svc: svc} }

// RegisterRoutes registers the loan routes on the given Gin engine.
func (h *LoanHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/loans", h.createLoan)
	r.GET("/loans", h.listLoans)
	r.GET("/loans/:id", h.getLoan)
	r.POST("/loans/:id/approve", h.approveLoan)
	r.POST("/loans/:id/invest", h.investInLoan)
	r.POST("/loans/:id/disburse", h.disburseLoan)
}

// createLoan handles POST /loans. It expects a JSON payload
// containing borrower_id, principal, rate, roi and optionally
// agreement_letter_url.
func (h *LoanHandler) createLoan(c *gin.Context) {
	var req struct {
		BorrowerID         string  `json:"borrower_id" binding:"required"`
		Principal          float64 `json:"principal" binding:"required"`
		Rate               float64 `json:"rate" binding:"required"`
		ROI                float64 `json:"roi" binding:"required"`
		AgreementLetterURL string  `json:"agreement_letter_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan := domain.Loan{
		BorrowerID:         req.BorrowerID,
		Principal:          req.Principal,
		Rate:               req.Rate,
		ROI:                req.ROI,
		AgreementLetterURL: req.AgreementLetterURL,
	}
	created, err := h.svc.CreateLoan(context.Background(), loan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// listLoans handles GET /loans. It returns all loans without
// pagination. In a real system pagination parameters should be
// supported.
func (h *LoanHandler) listLoans(c *gin.Context) {
	loans, err := h.svc.ListLoans(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loans)
}

// getLoan handles GET /loans/:id. It returns a single loan with
// nested approval, investments and disbursement information.
func (h *LoanHandler) getLoan(c *gin.Context) {
	id := c.Param("id")
	loan, err := h.svc.GetLoanByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

// approveLoan handles POST /loans/:id/approve. It expects
// picture_url, employee_id and approval_date in the body. The
// approval_date must be a valid RFC3339 timestamp.
func (h *LoanHandler) approveLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PictureURL   string `json:"picture_url" binding:"required"`
		EmployeeID   string `json:"employee_id" binding:"required"`
		ApprovalDate string `json:"approval_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, err := time.Parse(time.RFC3339, req.ApprovalDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid approval_date; must be RFC3339"})
		return
	}
	loan, err := h.svc.ApproveLoan(context.Background(), id, req.PictureURL, req.EmployeeID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

// investInLoan handles POST /loans/:id/invest. It accepts optional
// investor_id or name/email to identify or create an investor.
func (h *LoanHandler) investInLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		InvestorID    string  `json:"investor_id"`
		InvestorName  string  `json:"investor_name"`
		InvestorEmail string  `json:"investor_email"`
		Amount        float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan, err := h.svc.InvestInLoan(context.Background(), id, req.InvestorID, req.InvestorName, req.InvestorEmail, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

// disburseLoan handles POST /loans/:id/disburse. It expects
// agreement_url, employee_id and disbursement_date in RFC3339
// format.
func (h *LoanHandler) disburseLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		AgreementURL     string `json:"agreement_url" binding:"required"`
		EmployeeID       string `json:"employee_id" binding:"required"`
		DisbursementDate string `json:"disbursement_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, err := time.Parse(time.RFC3339, req.DisbursementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid disbursement_date; must be RFC3339"})
		return
	}
	loan, err := h.svc.DisburseLoan(context.Background(), id, req.AgreementURL, req.EmployeeID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}
