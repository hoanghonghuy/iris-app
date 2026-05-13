<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Heart, MessageCircle, Pencil, SendHorizontal, Trash2, X } from 'lucide-vue-next'
import { POST_SCOPE_LABELS, POST_TYPE_META, POST_TYPE_OPTIONS } from '../helpers/postConfig'
import { formatDateTime } from '@/helpers/dateFormatter'
import { extractErrorMessage } from '@/helpers/errorHandler'
import { usePostInteractions } from '@/composables/usePostInteractions'
import { teacherPostService } from '@/services/postService'

const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
  authorLabel: {
    type: String,
    default: 'Giáo viên',
  },
  audience: {
    type: String,
    default: 'teacher',
  },
  enableInteractions: {
    type: Boolean,
    default: true,
  },
  enableShare: {
    type: Boolean,
    default: true,
  },
  enableTeacherManage: {
    type: Boolean,
    default: false,
  },
  editableClassId: {
    type: String,
    default: '',
  },
  editableStudents: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['patch-post', 'delete-post'])

// Post interaction composable
const {
  processing: processingLike,
  loadingComments,
  submittingComment,
  error: interactionError,
  toggleLike,
  loadComments: loadCommentsFromService,
  createComment: createCommentViaService,
  share: sharePost,
} = usePostInteractions(props.audience)

// Local state
const liked = ref(props.post.liked_by_me)
const likeCount = ref(props.post.like_count || 0)
const shareCount = ref(props.post.share_count || 0)
const commentCount = ref(props.post.comment_count || 0)
const showComments = ref(false)
const comments = ref([])
const commentsLoaded = ref(false)
const commentDraft = ref('')
const processingShare = ref(false)
const editing = ref(false)
const editDraft = ref('')
const editScopeType = ref('class')
const editPostType = ref('announcement')
const editStudentId = ref('')
const manageBusy = ref(false)
const manageError = ref('')
const EDIT_HINT = 'Ctrl+Enter để lưu, Esc để hủy'

const normalizedOriginalContent = computed(() => (props.post.content || '').trim())
const normalizedDraftContent = computed(() => editDraft.value.trim())
const normalizedOriginalScopeType = computed(() => props.post.scope_type || 'class')
const normalizedOriginalPostType = computed(() => props.post.type || 'announcement')
const normalizedOriginalStudentId = computed(() => props.post.student_id || '')
const hasDraftChanges = computed(
  () =>
    normalizedDraftContent.value !== normalizedOriginalContent.value ||
    editScopeType.value !== normalizedOriginalScopeType.value ||
    editPostType.value !== normalizedOriginalPostType.value ||
    (editScopeType.value === 'student' && editStudentId.value !== normalizedOriginalStudentId.value),
)
const draftLength = computed(() => normalizedDraftContent.value.length)

watch(
  () => props.post,
  (post) => {
    liked.value = post.liked_by_me
    likeCount.value = post.like_count || 0
    shareCount.value = post.share_count || 0
    commentCount.value = post.comment_count || 0
    if (!editing.value) {
      editDraft.value = post.content || ''
    }
  },
  { deep: true },
)

function patchPost(patch) {
  emit('patch-post', props.post.post_id, patch)
}

function startEdit() {
  manageError.value = ''
  editDraft.value = props.post.content || ''
  editScopeType.value = props.post.scope_type || 'class'
  editPostType.value = props.post.type || 'announcement'
  editStudentId.value = props.post.student_id || props.editableStudents[0]?.student_id || ''
  editing.value = true
}

function cancelEdit() {
  if (hasDraftChanges.value && !globalThis.confirm('Bạn có thay đổi chưa lưu. Hủy chỉnh sửa?')) {
    return
  }
  editing.value = false
  editDraft.value = props.post.content || ''
  editScopeType.value = props.post.scope_type || 'class'
  editPostType.value = props.post.type || 'announcement'
  editStudentId.value = props.post.student_id || ''
  manageError.value = ''
}

async function saveEdit() {
  if (!props.enableTeacherManage) return
  const content = normalizedDraftContent.value
  if (!content) {
    manageError.value = 'Nội dung không được để trống.'
    return
  }
  if (editScopeType.value !== 'class' && editScopeType.value !== 'student') {
    manageError.value = 'Phạm vi bài đăng không hợp lệ.'
    return
  }
  if (!editPostType.value) {
    manageError.value = 'Vui lòng chọn loại bài.'
    return
  }
  if (editScopeType.value === 'student' && !editStudentId.value) {
    manageError.value = 'Vui lòng chọn học sinh.'
    return
  }
  if (!props.editableClassId) {
    manageError.value = 'Không xác định được lớp đang chỉnh sửa.'
    return
  }
  if (!hasDraftChanges.value) {
    manageError.value = 'Bạn chưa thay đổi nội dung.'
    return
  }
  manageBusy.value = true
  manageError.value = ''
  try {
    const payload = {
      scope_type: editScopeType.value,
      class_id: props.editableClassId,
      student_id: editScopeType.value === 'student' ? editStudentId.value : undefined,
      type: editPostType.value,
      content,
    }
    await teacherPostService.updatePost(props.post.post_id, payload)
    patchPost({
      scope_type: payload.scope_type,
      class_id: payload.class_id,
      student_id: payload.student_id ?? null,
      type: payload.type,
      content,
      updated_at: new Date().toISOString(),
    })
    editing.value = false
  } catch (err) {
    manageError.value = extractErrorMessage(err) || 'Không thể lưu bài.'
  } finally {
    manageBusy.value = false
  }
}

async function requestDelete() {
  if (!props.enableTeacherManage) return
  if (!globalThis.confirm('Xóa bài đăng này?')) return
  manageBusy.value = true
  manageError.value = ''
  try {
    await teacherPostService.deletePost(props.post.post_id)
    emit('delete-post', props.post.post_id)
  } catch (err) {
    manageError.value = extractErrorMessage(err) || 'Không thể xóa bài.'
  } finally {
    manageBusy.value = false
  }
}

function initials(text) {
  return (text || props.authorLabel || 'GV').slice(0, 2).toUpperCase()
}

async function handleLikeToggle() {
  try {
    const payload = await toggleLike(props.post.post_id)
    liked.value = payload.liked_by_me
    likeCount.value = payload.like_count
    patchPost({ liked_by_me: payload.liked_by_me, like_count: payload.like_count })
  } catch {
    // Error already handled by composable
  }
}

async function loadComments() {
  try {
    const loadedComments = await loadCommentsFromService(props.post.post_id, {
      limit: 50,
      offset: 0,
    })
    comments.value = loadedComments
    commentsLoaded.value = true
  } catch {
    // Error already handled by composable
  }
}

async function toggleComments() {
  showComments.value = !showComments.value
  if (showComments.value && !commentsLoaded.value) {
    await loadComments()
  }
}

async function submitComment() {
  const content = commentDraft.value.trim()
  if (!content) return

  try {
    const payload = await createCommentViaService(props.post.post_id, content)
    if (payload && payload.comment) {
      comments.value = [payload.comment, ...comments.value]
    }
    commentCount.value = payload?.comment_count ?? commentCount.value + 1
    patchPost({ comment_count: commentCount.value })
    commentDraft.value = ''
  } catch {
    // Error already handled by composable
  }
}

async function handleShare() {
  processingShare.value = true
  try {
    const payload = await sharePost(props.post.post_id)
    shareCount.value = payload.share_count
    patchPost({ share_count: payload.share_count })
  } catch {
    // Error already handled by composable
  } finally {
    processingShare.value = false
  }
}

function handleEditKeydown(event) {
  if (event.key === 'Escape') {
    event.preventDefault()
    cancelEdit()
    return
  }
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
    void saveEdit()
  }
}
</script>

<template>
  <article class="card post-card">
    <div class="post-card__body">
      <div class="post-card__header">
        <div class="post-card__author">
          <div class="post-card__avatar">{{ initials(authorLabel) }}</div>
          <div class="post-card__author-copy">
            <p class="post-card__author-name">{{ authorLabel }}</p>
            <p class="post-card__meta">
              {{ formatDateTime(post.created_at) }} ·
              {{ POST_SCOPE_LABELS[post.scope_type] || post.scope_type }}
            </p>
          </div>
        </div>
        <div class="post-card__header-actions">
          <span class="badge" :class="POST_TYPE_META[post.type]?.badgeClass || 'badge--outline'">
            {{ POST_TYPE_META[post.type]?.label || post.type }}
          </span>
          <div v-if="enableTeacherManage" class="post-card__manage">
            <button
              v-if="!editing"
              type="button"
              class="post-card__manage-btn"
              :disabled="manageBusy"
              title="Sửa bài"
              @click="startEdit"
            >
              <Pencil :size="16" />
            </button>
            <button
              v-if="!editing"
              type="button"
              class="post-card__manage-btn post-card__manage-btn--danger"
              :disabled="manageBusy"
              title="Xóa bài"
              @click="requestDelete"
            >
              <Trash2 :size="16" />
            </button>
          </div>
        </div>
      </div>

      <template v-if="editing">
        <div class="post-card__edit-grid">
          <div class="form-group mb-0">
            <label class="form-label">Phạm vi</label>
            <select v-model="editScopeType" class="form-input">
              <option value="class">Cả lớp</option>
              <option value="student">Từng HS</option>
            </select>
          </div>
          <div class="form-group mb-0">
            <label class="form-label">Loại bài</label>
            <select v-model="editPostType" class="form-input">
              <option
                v-for="option in POST_TYPE_OPTIONS"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </select>
          </div>
          <div v-if="editScopeType === 'student'" class="form-group mb-0">
            <label class="form-label">Học sinh</label>
            <select v-model="editStudentId" class="form-input">
              <option
                v-for="student in editableStudents"
                :key="student.student_id"
                :value="student.student_id"
              >
                {{ student.full_name }}
              </option>
            </select>
          </div>
        </div>
        <textarea
          v-model="editDraft"
          class="form-input post-card__edit-area"
          rows="4"
          aria-label="Sửa nội dung bài đăng"
          @keydown="handleEditKeydown"
        />
        <div class="post-card__edit-meta">
          <p class="post-card__edit-hint">{{ EDIT_HINT }}</p>
          <p class="post-card__edit-count">{{ draftLength }} ký tự</p>
        </div>
        <div class="post-card__edit-actions">
          <button type="button" class="btn btn--sm btn--outline" :disabled="manageBusy" @click="cancelEdit">
            <X :size="16" />
            Hủy
          </button>
          <button
            type="button"
            class="btn btn--sm btn--primary"
            :disabled="manageBusy || !hasDraftChanges"
            @click="saveEdit"
          >
            Lưu
          </button>
        </div>
      </template>
      <p v-else class="post-card__content">{{ post.content }}</p>
      <p v-if="manageError" class="form-error">{{ manageError }}</p>

      <template v-if="enableInteractions">
        <div class="post-card__actions">
          <button
            type="button"
            class="post-card__action"
            :class="{ 'post-card__action--active': liked }"
            :disabled="processingLike"
            @click="handleLikeToggle"
          >
            <Heart :size="16" />
            Thích ({{ likeCount }})
          </button>

          <button type="button" class="post-card__action" @click="toggleComments">
            <MessageCircle :size="16" />
            Bình luận ({{ commentCount }})
          </button>

          <button
            v-if="enableShare"
            type="button"
            class="post-card__action"
            :disabled="processingShare"
            @click="handleShare"
          >
            <SendHorizontal :size="16" />
            Chia sẻ ({{ shareCount }})
          </button>
        </div>

        <p v-if="interactionError" class="form-error">{{ interactionError }}</p>

        <div v-if="showComments" class="post-card__comments">
          <form class="post-card__comment-form" @submit.prevent="submitComment">
            <input
              v-model="commentDraft"
              class="form-input"
              placeholder="Viết bình luận..."
              aria-label="Viết bình luận"
            />
            <button
              type="submit"
              class="btn btn--primary btn--sm"
              :disabled="!commentDraft.trim() || submittingComment"
            >
              Gửi
            </button>
          </form>

          <p v-if="loadingComments" class="text-xs text-muted">Đang tải bình luận...</p>
          <ul v-else-if="comments.length > 0" class="post-card__comment-list">
            <li v-for="comment in comments" :key="comment.comment_id" class="post-card__comment">
              <p class="post-card__comment-author">{{ comment.author_display }}</p>
              <p>{{ comment.content }}</p>
            </li>
          </ul>
        </div>
      </template>
    </div>
  </article>
</template>

<style scoped>
.post-card {
  overflow: hidden;
}

.post-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
}

.post-card__header,
.post-card__author,
.post-card__actions,
.post-card__comment-form {
  display: flex;
  align-items: center;
}

.post-card__header {
  justify-content: space-between;
  gap: var(--spacing-3);
}

.post-card__header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  flex-shrink: 0;
}

.post-card__manage {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-1);
}

.post-card__manage-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-2);
  border: 0;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-muted);
}

.post-card__manage-btn:hover:not(:disabled) {
  background: var(--color-background);
  color: var(--color-text);
}

.post-card__manage-btn--danger:hover:not(:disabled) {
  color: var(--color-danger, #b91c1c);
}

.post-card__edit-area {
  width: 100%;
  resize: vertical;
  min-height: 6rem;
}

.post-card__edit-grid {
  display: grid;
  gap: var(--spacing-3);
}

.post-card__edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.post-card__edit-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-2);
}

.post-card__edit-hint,
.post-card__edit-count {
  margin: 0;
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

@media (min-width: 640px) {
  .post-card__edit-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

.post-card__author {
  min-width: 0;
  gap: var(--spacing-3);
}

.post-card__avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary);
  color: var(--color-on-primary);
  font-weight: 700;
}

.post-card__author-copy {
  min-width: 0;
}

.post-card__author-name {
  margin: 0;
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.post-card__meta {
  margin: 0;
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.post-card__content {
  white-space: pre-line;
  font-size: var(--font-size-sm);
  line-height: 1.625;
}

.post-card__actions {
  flex-wrap: wrap;
  gap: var(--spacing-2);
  border-top: 1px solid var(--color-border);
  padding-top: var(--spacing-3);
}

.post-card__action {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-1);
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  border-radius: var(--radius-md);
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-sm);
}

.post-card__action:hover {
  background: var(--color-background);
  color: var(--color-text);
}

.post-card__action--active {
  color: var(--color-primary);
}

.post-card__comments {
  border-top: 1px solid var(--color-border);
  padding-top: var(--spacing-3);
}

.post-card__comment-form {
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-3);
}

.post-card__comment-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  list-style: none;
}

.post-card__comment {
  background: var(--color-background);
  border-radius: var(--radius-md);
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-sm);
}

.post-card__comment-author {
  margin: 0 0 var(--spacing-1);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
  font-weight: 700;
}
</style>
