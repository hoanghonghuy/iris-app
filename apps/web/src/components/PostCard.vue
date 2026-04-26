<script setup>
import { ref, watch } from 'vue'
import { Heart, MessageCircle, SendHorizontal } from 'lucide-vue-next'
import { teacherService } from '../services/teacherService'
import { parentService } from '../services/parentService'
import { POST_SCOPE_LABELS, POST_TYPE_META } from '../helpers/postConfig'

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
})

const emit = defineEmits(['patch-post'])

const liked = ref(props.post.liked_by_me)
const likeCount = ref(props.post.like_count || 0)
const shareCount = ref(props.post.share_count || 0)
const commentCount = ref(props.post.comment_count || 0)
const showComments = ref(false)
const comments = ref([])
const commentsLoaded = ref(false)
const commentDraft = ref('')
const loadingComments = ref(false)
const processingLike = ref(false)
const processingShare = ref(false)
const submittingComment = ref(false)
const interactionError = ref('')

const services = {
  teacher: teacherService,
  parent: parentService,
}

watch(
  () => props.post,
  (post) => {
    liked.value = post.liked_by_me
    likeCount.value = post.like_count || 0
    shareCount.value = post.share_count || 0
    commentCount.value = post.comment_count || 0
  },
  { deep: true },
)

function patchPost(patch) {
  emit('patch-post', props.post.post_id, patch)
}

function formatDateTime(value) {
  if (!value) return ''
  return new Date(value).toLocaleString('vi-VN')
}

function initials(text) {
  return (text || props.authorLabel || 'GV').slice(0, 2).toUpperCase()
}

async function handleLikeToggle() {
  processingLike.value = true
  interactionError.value = ''
  try {
    const response = await services[props.audience].togglePostLike(props.post.post_id)
    const payload = response?.data ?? response
    liked.value = payload.liked_by_me
    likeCount.value = payload.like_count
    patchPost({ liked_by_me: payload.liked_by_me, like_count: payload.like_count })
  } catch {
    interactionError.value = 'Không thể cập nhật lượt thích. Vui lòng thử lại.'
  } finally {
    processingLike.value = false
  }
}

async function loadComments() {
  loadingComments.value = true
  interactionError.value = ''
  try {
    const response = await services[props.audience].getPostComments(props.post.post_id, {
      limit: 50,
      offset: 0,
    })
    comments.value = response?.data ?? []
    commentsLoaded.value = true
  } catch {
    interactionError.value = 'Không thể tải bình luận. Vui lòng thử lại.'
  } finally {
    loadingComments.value = false
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

  submittingComment.value = true
  interactionError.value = ''
  try {
    const response = await services[props.audience].createPostComment(props.post.post_id, {
      content,
    })
    const payload = response?.data ?? response
    if (payload.comment) {
      comments.value = [payload.comment, ...comments.value]
    }
    commentCount.value = payload.comment_count ?? commentCount.value + 1
    patchPost({ comment_count: commentCount.value })
    commentDraft.value = ''
  } catch {
    interactionError.value = 'Không thể gửi bình luận. Vui lòng thử lại.'
  } finally {
    submittingComment.value = false
  }
}

async function handleShare() {
  processingShare.value = true
  interactionError.value = ''
  try {
    const response = await services[props.audience].sharePost(props.post.post_id)
    const payload = response?.data ?? response
    shareCount.value = payload.share_count
    patchPost({ share_count: payload.share_count })
  } catch {
    interactionError.value = 'Không thể chia sẻ bài viết. Vui lòng thử lại.'
  } finally {
    processingShare.value = false
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
        <span class="badge" :class="POST_TYPE_META[post.type]?.badgeClass || 'badge--outline'">
          {{ POST_TYPE_META[post.type]?.label || post.type }}
        </span>
      </div>

      <p class="post-card__content">{{ post.content }}</p>

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
