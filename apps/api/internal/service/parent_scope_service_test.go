package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeParentScopeRepo struct {
	isParentCalls        int
	isParentParentUserID uuid.UUID
	isParentStudentID    uuid.UUID
	isParentResult       bool
	isParentErr          error

	listMyChildrenCalls int
	listMyChildrenUser  uuid.UUID
	listMyChildrenRes   []model.Student
	listMyChildrenErr   error

	listClassPostsCalls int
	listClassParentID   uuid.UUID
	listClassStudentID  uuid.UUID
	listClassLimit      int
	listClassOffset     int
	listClassRes        []model.Post
	listClassTotal      int
	listClassErr        error

	listStudentPostsCalls int
	listStudentParentID   uuid.UUID
	listStudentStudentID  uuid.UUID
	listStudentLimit      int
	listStudentOffset     int
	listStudentRes        []model.Post
	listStudentTotal      int
	listStudentErr        error

	listAllPostsCalls int
	listAllParentID   uuid.UUID
	listAllStudentID  uuid.UUID
	listAllLimit      int
	listAllOffset     int
	listAllRes        []model.Post
	listAllTotal      int
	listAllErr        error

	feedCalls    int
	feedParentID uuid.UUID
	feedLimit    int
	feedOffset   int
	feedRes      []model.Post
	feedTotal    int
	feedErr      error

	countChildrenCalls int
	countChildrenUser  uuid.UUID
	countChildrenRes   int
	countChildrenErr   error

	countRecentPostsCalls int
	countRecentPostsUser  uuid.UUID
	countRecentPostsSince time.Time
	countRecentPostsRes   int
	countRecentPostsErr   error

	countRecentHealthCalls int
	countRecentHealthUser  uuid.UUID
	countRecentHealthSince time.Time
	countRecentHealthRes   int
	countRecentHealthErr   error
}

func (f *fakeParentScopeRepo) IsParentOfStudent(_ context.Context, parentUserID, studentID uuid.UUID) (bool, error) {
	f.isParentCalls++
	f.isParentParentUserID = parentUserID
	f.isParentStudentID = studentID
	return f.isParentResult, f.isParentErr
}

func (f *fakeParentScopeRepo) ListMyChildren(_ context.Context, parentUserID uuid.UUID) ([]model.Student, error) {
	f.listMyChildrenCalls++
	f.listMyChildrenUser = parentUserID
	return f.listMyChildrenRes, f.listMyChildrenErr
}

func (f *fakeParentScopeRepo) ListMyChildClassPosts(_ context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.listClassPostsCalls++
	f.listClassParentID = parentUserID
	f.listClassStudentID = studentID
	f.listClassLimit = limit
	f.listClassOffset = offset
	return f.listClassRes, f.listClassTotal, f.listClassErr
}

func (f *fakeParentScopeRepo) ListMyChildStudentPosts(_ context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.listStudentPostsCalls++
	f.listStudentParentID = parentUserID
	f.listStudentStudentID = studentID
	f.listStudentLimit = limit
	f.listStudentOffset = offset
	return f.listStudentRes, f.listStudentTotal, f.listStudentErr
}

func (f *fakeParentScopeRepo) ListAllMyChildPosts(_ context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.listAllPostsCalls++
	f.listAllParentID = parentUserID
	f.listAllStudentID = studentID
	f.listAllLimit = limit
	f.listAllOffset = offset
	return f.listAllRes, f.listAllTotal, f.listAllErr
}

func (f *fakeParentScopeRepo) GetMyFeed(_ context.Context, parentUserID uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.feedCalls++
	f.feedParentID = parentUserID
	f.feedLimit = limit
	f.feedOffset = offset
	return f.feedRes, f.feedTotal, f.feedErr
}

func (f *fakeParentScopeRepo) CountMyChildren(_ context.Context, parentUserID uuid.UUID) (int, error) {
	f.countChildrenCalls++
	f.countChildrenUser = parentUserID
	return f.countChildrenRes, f.countChildrenErr
}

func (f *fakeParentScopeRepo) CountMyRecentPosts(_ context.Context, parentUserID uuid.UUID, since time.Time) (int, error) {
	f.countRecentPostsCalls++
	f.countRecentPostsUser = parentUserID
	f.countRecentPostsSince = since
	return f.countRecentPostsRes, f.countRecentPostsErr
}

func (f *fakeParentScopeRepo) CountMyRecentHealthAlerts(_ context.Context, parentUserID uuid.UUID, since time.Time) (int, error) {
	f.countRecentHealthCalls++
	f.countRecentHealthUser = parentUserID
	f.countRecentHealthSince = since
	return f.countRecentHealthRes, f.countRecentHealthErr
}

type fakePostInteractionRepo struct {
	canAccessCalls    int
	canAccessParentID uuid.UUID
	canAccessPostID   uuid.UUID
	canAccessRes      bool
	canAccessErr      error

	toggleCalls  int
	toggleUserID uuid.UUID
	togglePostID uuid.UUID
	toggleLiked  bool
	toggleCount  int
	toggleErr    error

	addCommentCalls   int
	addCommentUserID  uuid.UUID
	addCommentPostID  uuid.UUID
	addCommentContent string
	addCommentRes     model.PostComment
	addCommentErr     error

	countCommentsCalls int
	countCommentsPost  uuid.UUID
	countCommentsRes   int
	countCommentsErr   error

	listCommentsCalls  int
	listCommentsPostID uuid.UUID
	listCommentsLimit  int
	listCommentsOffset int
	listCommentsRes    []model.PostComment
	listCommentsTotal  int
	listCommentsErr    error

	shareCalls  int
	shareUserID uuid.UUID
	sharePostID uuid.UUID
	shareCount  int
	shareErr    error
}

func (f *fakePostInteractionRepo) ParentCanAccessPost(_ context.Context, parentUserID, postID uuid.UUID) (bool, error) {
	f.canAccessCalls++
	f.canAccessParentID = parentUserID
	f.canAccessPostID = postID
	return f.canAccessRes, f.canAccessErr
}

func (f *fakePostInteractionRepo) ToggleLike(_ context.Context, userID, postID uuid.UUID) (bool, int, error) {
	f.toggleCalls++
	f.toggleUserID = userID
	f.togglePostID = postID
	return f.toggleLiked, f.toggleCount, f.toggleErr
}

func (f *fakePostInteractionRepo) AddComment(_ context.Context, userID, postID uuid.UUID, content string) (model.PostComment, error) {
	f.addCommentCalls++
	f.addCommentUserID = userID
	f.addCommentPostID = postID
	f.addCommentContent = content
	return f.addCommentRes, f.addCommentErr
}

func (f *fakePostInteractionRepo) CountComments(_ context.Context, postID uuid.UUID) (int, error) {
	f.countCommentsCalls++
	f.countCommentsPost = postID
	return f.countCommentsRes, f.countCommentsErr
}

func (f *fakePostInteractionRepo) ListComments(_ context.Context, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error) {
	f.listCommentsCalls++
	f.listCommentsPostID = postID
	f.listCommentsLimit = limit
	f.listCommentsOffset = offset
	return f.listCommentsRes, f.listCommentsTotal, f.listCommentsErr
}

func (f *fakePostInteractionRepo) AddShare(_ context.Context, userID, postID uuid.UUID) (int, error) {
	f.shareCalls++
	f.shareUserID = userID
	f.sharePostID = postID
	return f.shareCount, f.shareErr
}

type fakeParentScopeAppointmentRepo struct {
	countUpcomingCalls int
	countUpcomingUser  uuid.UUID
	countUpcomingRes   int
	countUpcomingErr   error
}

func (f *fakeParentScopeAppointmentRepo) CountParentUpcomingAppointments(_ context.Context, parentUserID uuid.UUID) (int, error) {
	f.countUpcomingCalls++
	f.countUpcomingUser = parentUserID
	return f.countUpcomingRes, f.countUpcomingErr
}

func TestNormalizeParentScopeLimit(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "default when non-positive", in: 0, want: 20},
		{name: "clamp max", in: 1000, want: 100},
		{name: "keep valid", in: 40, want: 40},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeParentScopeLimit(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeParentScopeLimit(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestEnsureParentStudentAccess(t *testing.T) {
	parentUserID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("db error")

	tests := []struct {
		name         string
		parentUserID uuid.UUID
		studentID    uuid.UUID
		repoResult   bool
		repoErr      error
		wantErr      error
		wantCalls    int
	}{
		{name: "invalid parent user id", parentUserID: uuid.Nil, studentID: studentID, wantErr: ErrInvalidUserID, wantCalls: 0},
		{name: "invalid student id", parentUserID: parentUserID, studentID: uuid.Nil, wantErr: ErrInvalidUserID, wantCalls: 0},
		{name: "repo error", parentUserID: parentUserID, studentID: studentID, repoErr: sentinelErr, wantErr: sentinelErr, wantCalls: 1},
		{name: "not parent of student", parentUserID: parentUserID, studentID: studentID, repoResult: false, wantErr: ErrForbidden, wantCalls: 1},
		{name: "success", parentUserID: parentUserID, studentID: studentID, repoResult: true, wantCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeParentScopeRepo{isParentResult: tc.repoResult, isParentErr: tc.repoErr}
			svc := &ParentScopeService{parentScopeRepo: repo}

			err := svc.ensureParentStudentAccess(context.Background(), tc.parentUserID, tc.studentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ensureParentStudentAccess() error = %v", err)
			}

			if repo.isParentCalls != tc.wantCalls {
				t.Fatalf("isParent calls = %d, want %d", repo.isParentCalls, tc.wantCalls)
			}
		})
	}
}

func TestListMyChildren(t *testing.T) {
	parentUserID := uuid.New()
	sentinelErr := errors.New("repo failed")
	students := []model.Student{{StudentID: uuid.New(), FullName: "Child"}}

	t.Run("invalid parent user id", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{}}
		_, err := svc.ListMyChildren(context.Background(), uuid.Nil)
		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &fakeParentScopeRepo{listMyChildrenErr: sentinelErr}
		svc := &ParentScopeService{parentScopeRepo: repo}
		_, err := svc.ListMyChildren(context.Background(), parentUserID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeParentScopeRepo{listMyChildrenRes: students}
		svc := &ParentScopeService{parentScopeRepo: repo}
		got, err := svc.ListMyChildren(context.Background(), parentUserID)
		if err != nil {
			t.Fatalf("ListMyChildren() error = %v", err)
		}
		if len(got) != 1 || got[0].FullName != "Child" {
			t.Fatalf("unexpected students: %#v", got)
		}
		if repo.listMyChildrenUser != parentUserID {
			t.Fatalf("parent user id not forwarded")
		}
	})
}

func TestListMyChildPostsAndFeed(t *testing.T) {
	parentUserID := uuid.New()
	studentID := uuid.New()
	posts := []model.Post{{PostID: uuid.New(), Content: "hello"}}
	sentinelErr := errors.New("repo failed")

	t.Run("class posts normalizes limit and checks access", func(t *testing.T) {
		repo := &fakeParentScopeRepo{isParentResult: true, listClassRes: posts, listClassTotal: 1}
		svc := &ParentScopeService{parentScopeRepo: repo}
		got, total, err := svc.ListMyChildClassPosts(context.Background(), parentUserID, studentID, 0, 3)
		if err != nil {
			t.Fatalf("ListMyChildClassPosts() error = %v", err)
		}
		if len(got) != 1 || total != 1 {
			t.Fatalf("unexpected result: len=%d total=%d", len(got), total)
		}
		if repo.listClassLimit != 20 || repo.listClassOffset != 3 {
			t.Fatalf("limit/offset forwarded = %d/%d", repo.listClassLimit, repo.listClassOffset)
		}
	})

	t.Run("student posts wraps repo error", func(t *testing.T) {
		repo := &fakeParentScopeRepo{isParentResult: true, listStudentErr: sentinelErr}
		svc := &ParentScopeService{parentScopeRepo: repo}
		_, _, err := svc.ListMyChildStudentPosts(context.Background(), parentUserID, studentID, 25, 0)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("all posts denied by access check", func(t *testing.T) {
		repo := &fakeParentScopeRepo{isParentResult: false}
		svc := &ParentScopeService{parentScopeRepo: repo}
		_, _, err := svc.ListAllMyChildPosts(context.Background(), parentUserID, studentID, 25, 0)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
		if repo.listAllPostsCalls != 0 {
			t.Fatalf("list all posts should not be called when access denied")
		}
	})

	t.Run("feed invalid parent user", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{}}
		_, _, err := svc.GetMyFeed(context.Background(), uuid.Nil, 20, 0)
		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
		}
	})

	t.Run("feed normalizes limit and wraps repo error", func(t *testing.T) {
		repo := &fakeParentScopeRepo{feedErr: sentinelErr}
		svc := &ParentScopeService{parentScopeRepo: repo}
		_, _, err := svc.GetMyFeed(context.Background(), parentUserID, 999, 9)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
		if repo.feedLimit != 100 || repo.feedOffset != 9 {
			t.Fatalf("feed limit/offset forwarded = %d/%d", repo.feedLimit, repo.feedOffset)
		}
	})
}

func TestTogglePostLike(t *testing.T) {
	parentUserID := uuid.New()
	postID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name         string
		parentUserID uuid.UUID
		postID       uuid.UUID
		accessRes    bool
		accessErr    error
		toggleLiked  bool
		toggleCount  int
		toggleErr    error
		wantErr      error
		wantLiked    bool
		wantCount    int
		wantCalls    int
	}{
		{name: "invalid parent user id", parentUserID: uuid.Nil, postID: postID, wantErr: ErrInvalidUserID, wantCalls: 0},
		{name: "invalid post id", parentUserID: parentUserID, postID: uuid.Nil, wantErr: ErrInvalidValue, wantCalls: 0},
		{name: "access check error", parentUserID: parentUserID, postID: postID, accessErr: sentinelErr, wantErr: sentinelErr, wantCalls: 1},
		{name: "forbidden", parentUserID: parentUserID, postID: postID, accessRes: false, wantErr: ErrForbidden, wantCalls: 1},
		{name: "toggle like error", parentUserID: parentUserID, postID: postID, accessRes: true, toggleErr: sentinelErr, wantErr: sentinelErr, wantCalls: 1},
		{name: "success", parentUserID: parentUserID, postID: postID, accessRes: true, toggleLiked: true, toggleCount: 3, wantLiked: true, wantCount: 3, wantCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			postRepo := &fakePostInteractionRepo{canAccessRes: tc.accessRes, canAccessErr: tc.accessErr, toggleLiked: tc.toggleLiked, toggleCount: tc.toggleCount, toggleErr: tc.toggleErr}
			svc := &ParentScopeService{postInteractRepo: postRepo}

			liked, count, err := svc.TogglePostLike(context.Background(), tc.parentUserID, tc.postID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("TogglePostLike() error = %v", err)
			}
			if liked != tc.wantLiked || count != tc.wantCount {
				t.Fatalf("liked/count = %v/%d, want %v/%d", liked, count, tc.wantLiked, tc.wantCount)
			}
			if postRepo.canAccessCalls != tc.wantCalls {
				t.Fatalf("access calls = %d, want %d", postRepo.canAccessCalls, tc.wantCalls)
			}
		})
	}
}

func TestAddPostComment(t *testing.T) {
	parentUserID := uuid.New()
	postID := uuid.New()
	sentinelErr := errors.New("repo failed")
	comment := model.PostComment{CommentID: uuid.New(), PostID: postID, AuthorUserID: parentUserID, Content: "ok", CreatedAt: time.Now().UTC()}

	t.Run("reject empty content", func(t *testing.T) {
		svc := &ParentScopeService{postInteractRepo: &fakePostInteractionRepo{}}
		_, _, err := svc.AddPostComment(context.Background(), parentUserID, postID, "   ")
		if !errors.Is(err, ErrInvalidValue) {
			t.Fatalf("error = %v, want wrapped %v", err, ErrInvalidValue)
		}
	})

	t.Run("access error", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessErr: sentinelErr}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		_, _, err := svc.AddPostComment(context.Background(), parentUserID, postID, "hello")
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: false}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		_, _, err := svc.AddPostComment(context.Background(), parentUserID, postID, "hello")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("count comment error", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: true, addCommentRes: comment, countCommentsErr: sentinelErr}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		_, _, err := svc.AddPostComment(context.Background(), parentUserID, postID, " hello ")
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
		if postRepo.addCommentContent != "hello" {
			t.Fatalf("content should be trimmed, got %q", postRepo.addCommentContent)
		}
	})

	t.Run("success", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: true, addCommentRes: comment, countCommentsRes: 4}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		gotComment, gotCount, err := svc.AddPostComment(context.Background(), parentUserID, postID, "hello")
		if err != nil {
			t.Fatalf("AddPostComment() error = %v", err)
		}
		if gotComment.CommentID != comment.CommentID || gotCount != 4 {
			t.Fatalf("unexpected result: %#v count=%d", gotComment, gotCount)
		}
	})
}

func TestListPostCommentsAndSharePost(t *testing.T) {
	parentUserID := uuid.New()
	postID := uuid.New()
	comments := []model.PostComment{{CommentID: uuid.New(), PostID: postID, AuthorUserID: uuid.New(), Content: "x", CreatedAt: time.Now().UTC()}}
	sentinelErr := errors.New("repo failed")

	t.Run("list comments normalizes limit and wraps list error", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: true, listCommentsErr: sentinelErr}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		_, _, err := svc.ListPostComments(context.Background(), parentUserID, postID, 0, 5)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
		if postRepo.listCommentsLimit != 20 || postRepo.listCommentsOffset != 5 {
			t.Fatalf("limit/offset forwarded = %d/%d", postRepo.listCommentsLimit, postRepo.listCommentsOffset)
		}
	})

	t.Run("list comments success", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: true, listCommentsRes: comments, listCommentsTotal: 1}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		items, total, err := svc.ListPostComments(context.Background(), parentUserID, postID, 10, 1)
		if err != nil {
			t.Fatalf("ListPostComments() error = %v", err)
		}
		if len(items) != 1 || total != 1 {
			t.Fatalf("unexpected result: len=%d total=%d", len(items), total)
		}
	})

	t.Run("share post forbidden", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: false}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		_, err := svc.SharePost(context.Background(), parentUserID, postID)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("share post success", func(t *testing.T) {
		postRepo := &fakePostInteractionRepo{canAccessRes: true, shareCount: 7}
		svc := &ParentScopeService{postInteractRepo: postRepo}
		count, err := svc.SharePost(context.Background(), parentUserID, postID)
		if err != nil {
			t.Fatalf("SharePost() error = %v", err)
		}
		if count != 7 {
			t.Fatalf("share count = %d, want 7", count)
		}
	})
}

func TestGetMyAnalytics(t *testing.T) {
	parentUserID := uuid.New()
	sentinelErr := errors.New("repo failed")

	t.Run("invalid parent user", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{}, appointmentRepo: &fakeParentScopeAppointmentRepo{}}
		_, err := svc.GetMyAnalytics(context.Background(), uuid.Nil)
		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
		}
	})

	t.Run("count children error", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{countChildrenErr: sentinelErr}, appointmentRepo: &fakeParentScopeAppointmentRepo{}}
		_, err := svc.GetMyAnalytics(context.Background(), parentUserID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("count upcoming error", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{countChildrenRes: 1}, appointmentRepo: &fakeParentScopeAppointmentRepo{countUpcomingErr: sentinelErr}}
		_, err := svc.GetMyAnalytics(context.Background(), parentUserID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("count recent posts error", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{countChildrenRes: 1, countRecentPostsErr: sentinelErr}, appointmentRepo: &fakeParentScopeAppointmentRepo{countUpcomingRes: 2}}
		_, err := svc.GetMyAnalytics(context.Background(), parentUserID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("count health alerts error", func(t *testing.T) {
		svc := &ParentScopeService{parentScopeRepo: &fakeParentScopeRepo{countChildrenRes: 1, countRecentPostsRes: 3, countRecentHealthErr: sentinelErr}, appointmentRepo: &fakeParentScopeAppointmentRepo{countUpcomingRes: 2}}
		_, err := svc.GetMyAnalytics(context.Background(), parentUserID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		parentRepo := &fakeParentScopeRepo{countChildrenRes: 2, countRecentPostsRes: 5, countRecentHealthRes: 1}
		appointmentRepo := &fakeParentScopeAppointmentRepo{countUpcomingRes: 4}
		svc := &ParentScopeService{parentScopeRepo: parentRepo, appointmentRepo: appointmentRepo}

		nowBefore := time.Now()
		got, err := svc.GetMyAnalytics(context.Background(), parentUserID)
		nowAfter := time.Now()
		if err != nil {
			t.Fatalf("GetMyAnalytics() error = %v", err)
		}
		if got.TotalChildren != 2 || got.UpcomingAppointments != 4 || got.RecentPosts7d != 5 || got.RecentHealthAlerts7d != 1 {
			t.Fatalf("unexpected analytics: %#v", got)
		}
		if appointmentRepo.countUpcomingUser != parentUserID {
			t.Fatalf("upcoming count user not forwarded")
		}
		sinceLowerBound := nowBefore.AddDate(0, 0, -7).Add(-2 * time.Second)
		sinceUpperBound := nowAfter.AddDate(0, 0, -7).Add(2 * time.Second)
		if parentRepo.countRecentPostsSince.Before(sinceLowerBound) || parentRepo.countRecentPostsSince.After(sinceUpperBound) {
			t.Fatalf("recent posts since out of expected range: %v", parentRepo.countRecentPostsSince)
		}
		if parentRepo.countRecentHealthSince.Before(sinceLowerBound) || parentRepo.countRecentHealthSince.After(sinceUpperBound) {
			t.Fatalf("recent health since out of expected range: %v", parentRepo.countRecentHealthSince)
		}
	})
}
