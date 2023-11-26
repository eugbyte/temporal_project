package workflows

import (
	"errors"
	"testing"
	"time"

	db "encore.app/internal/db/bill"
	"github.com/bojanz/currency"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_CreateBill() {
	const billID = "ABC"

	s.env.OnActivity(CreateBillActivity, mock.Anything, billID).Return(db.Bill{}, nil)

	s.env.ExecuteWorkflow(CreateBill, billID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_CreateBill_FailedActivity() {
	const billID = "ABC"

	s.env.OnActivity(CreateBillActivity, mock.Anything, billID).Return(db.Bill{}, errors.New("mock_error"))

	s.env.ExecuteWorkflow(CreateBill, billID)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_CloseBill() {
	const billID = "ABC"
	mockBill := db.Bill{ID: billID}

	s.env.OnActivity(CloseBillActivity, mock.Anything, billID).Return(mockBill, nil)

	s.env.ExecuteWorkflow(CloseBill, billID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var actualBill db.Bill
	s.env.GetWorkflowResult(&actualBill)
	s.EqualValues(mockBill, actualBill)
}

func (s *UnitTestSuite) Test_ConfirmBillIncrease() {
	const billID = "ABC"
	usd, _ := currency.NewAmount("100", "USD")
	billDetail := db.TransactionDetail{
		Timestamp: 1000,
		ItemName:  "item",
		Amount:    usd,
	}

	s.env.OnActivity(IncreaseBillActivity, mock.Anything, billID, billDetail).Return(db.Bill{}, nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalChannel, true)
	}, time.Millisecond)

	s.env.ExecuteWorkflow(IncreaseBill, billID, billDetail)

	s.NoError(s.env.GetWorkflowError())
	s.True(s.env.IsWorkflowCompleted())
}

func (s *UnitTestSuite) Test_SanityCheck() {
	env := s.env

	// Mock activity implementation
	env.OnActivity(SanityCheckActivity, mock.Anything).Return(nil)
	env.ExecuteWorkflow(SanityCheck)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}
