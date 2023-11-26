package temporalbill

import (
	"testing"

	mockdb "encore.app/internal/db/bill"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	mockBillService BillService
	env             *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	// the db is mocked anyway
	s.mockBillService = mockdb.New()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_CreateBill_Activity() {
	activities := NewActivities(s.mockBillService)
	workflows := NewWorkFlows(s.mockBillService)
	const billID = "ABC"

	s.env.OnActivity(activities.CreateBill, mock.Anything, billID).Return(mockdb.Bill{}, nil)

	s.env.ExecuteWorkflow(workflows.CreateBill, billID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_CloseBill_Activity() {
	activities := NewActivities(s.mockBillService)
	workflows := NewWorkFlows(s.mockBillService)
	const billID = "ABC"

	s.env.OnActivity(activities.CloseBill, mock.Anything, billID).Return(mockdb.Bill{}, nil)

	s.env.ExecuteWorkflow(workflows.CloseBill, billID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
