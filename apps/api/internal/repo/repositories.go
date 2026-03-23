package repo

type Repositories struct {
	UserRepo          *UserRepo
	SchoolRepo        *SchoolRepo
	ClassRepo         *ClassRepo
	StudentRepo       *StudentRepo
	StudentParentRepo *StudentParentRepo
	ParentRepo        *ParentRepo
	ParentCodeRepo    *ParentCodeRepo
	TeacherRepo       *TeacherRepo
	TeacherClassRepo  *TeacherClassRepo
	TeacherScopeRepo  *TeacherScopeRepo
	ParentScopeRepo   *ParentScopeRepo
	SchoolAdminRepo   *SchoolAdminRepo
	ResetTokenRepo    *ResetTokenRepo
	ChatRepo          *ChatRepo
}
