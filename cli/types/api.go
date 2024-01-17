package types

import (
	"github.com/looplab/fsm"
	"github.com/plandex/plandex/shared"
)

type OnStreamPlanParams struct {
	Content string
	State   *fsm.FSM
	Err     error
}

type OnStreamPlan func(params OnStreamPlanParams)

type ApiClient interface {
	StartTrial() (*shared.StartTrialResponse, *shared.ApiError)
	ConvertTrial(req shared.ConvertTrialRequest) (*shared.SessionResponse, *shared.ApiError)

	CreateEmailVerification(email, customHost, userId string) (*shared.CreateEmailVerificationResponse, *shared.ApiError)

	CreateAccount(req shared.CreateAccountRequest, customHost string) (*shared.SessionResponse, *shared.ApiError)
	SignIn(req shared.SignInRequest, customHost string) (*shared.SessionResponse, *shared.ApiError)
	SignOut() *shared.ApiError

	ListOrgs() ([]*shared.Org, *shared.ApiError)
	CreateOrg(req shared.CreateOrgRequest) (*shared.CreateOrgResponse, *shared.ApiError)

	ListUsers() ([]*shared.User, *shared.ApiError)
	DeleteUser(userId string) *shared.ApiError

	InviteUser(req shared.InviteRequest) *shared.ApiError
	ListPendingInvites() ([]*shared.Invite, *shared.ApiError)
	ListAcceptedInvites() ([]*shared.Invite, *shared.ApiError)
	ListAllInvites() ([]*shared.Invite, *shared.ApiError)
	DeleteInvite(inviteId string) *shared.ApiError

	CreateProject(req shared.CreateProjectRequest) (*shared.CreateProjectResponse, *shared.ApiError)
	ListProjects() ([]*shared.Project, *shared.ApiError)
	SetProjectPlan(projectId string, req shared.SetProjectPlanRequest) *shared.ApiError
	RenameProject(projectId string, req shared.RenameProjectRequest) *shared.ApiError

	ListPlans(projectId string) ([]*shared.Plan, *shared.ApiError)
	ListArchivedPlans(projectId string) ([]*shared.Plan, *shared.ApiError)
	ListPlansRunning(projectId string) ([]*shared.Plan, *shared.ApiError)
	GetPlan(planId string) (*shared.Plan, *shared.ApiError)
	CreatePlan(projectId string, req shared.CreatePlanRequest) (*shared.CreatePlanResponse, *shared.ApiError)

	TellPlan(planId string, req shared.TellPlanRequest, onStreamPlan OnStreamPlan) *shared.ApiError
	DeletePlan(planId string) *shared.ApiError
	DeleteAllPlans(projectId string) *shared.ApiError
	ConnectPlan(planId string, onStreamPlan OnStreamPlan) *shared.ApiError
	StopPlan(planId string) *shared.ApiError

	ArchivePlan(planId string) *shared.ApiError

	GetCurrentPlanState(planId string) (*shared.CurrentPlanState, *shared.ApiError)
	ApplyPlan(planId string) *shared.ApiError
	RejectAllChanges(planId string) *shared.ApiError
	RejectResult(planId, resultId string) *shared.ApiError
	RejectReplacement(planId, resultId, replacementId string) *shared.ApiError

	LoadContext(planId string, req shared.LoadContextRequest) (*shared.LoadContextResponse, *shared.ApiError)
	UpdateContext(planId string, req shared.UpdateContextRequest) (*shared.UpdateContextResponse, *shared.ApiError)
	DeleteContext(planId string, req shared.DeleteContextRequest) (*shared.DeleteContextResponse, *shared.ApiError)
	ListContext(planId string) ([]*shared.Context, *shared.ApiError)

	ListConvo(planId string) ([]*shared.ConvoMessage, *shared.ApiError)
	ListLogs(planId string) (*shared.LogResponse, *shared.ApiError)
	RewindPlan(planId string, req shared.RewindPlanRequest) (*shared.RewindPlanResponse, *shared.ApiError)
}
