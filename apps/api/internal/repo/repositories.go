package repo

type Repositories struct {
	UserRepo            *UserRepo
	SchoolRepo          *SchoolRepo
	ClassRepo           *ClassRepo
	StudentRepo         *StudentRepo
	StudentParentRepo   *StudentParentRepo
	ParentRepo          *ParentRepo
	ParentCodeRepo      *ParentCodeRepo
	TeacherRepo         *TeacherRepo
	TeacherClassRepo    *TeacherClassRepo
	TeacherScopeRepo    *TeacherScopeRepo
	HealthLogRepo       *HealthLogRepo
	ParentScopeRepo     *ParentScopeRepo
	PostInteractionRepo *PostInteractionRepo
	SchoolAdminRepo     *SchoolAdminRepo
	ResetTokenRepo      *ResetTokenRepo
	ChatRepo            *ChatRepo
}
