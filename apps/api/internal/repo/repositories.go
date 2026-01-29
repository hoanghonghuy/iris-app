package repo

type Repositories struct {
	UserRepo          *UserRepo
	SchoolRepo        *SchoolRepo
	ClassRepo         *ClassRepo
	StudentRepo       *StudentRepo
	StudentParentRepo *StudentParentRepo
	ParentRepo        *ParentRepo
	TeacherRepo       *TeacherRepo
	TeacherClassRepo  *TeacherClassRepo
	TeacherScopeRepo  *TeacherScopeRepo
	ParentScopeRepo   *ParentScopeRepo
}
