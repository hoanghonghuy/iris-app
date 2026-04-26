<script setup>
import { AlertCircle, LoaderCircle, MessageSquare } from 'lucide-vue-next'
import PostCard from '../../components/PostCard.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import ParentFeedSummaryCard from './feed/ParentFeedSummaryCard.vue'
import { useParentFeedPage } from './feed/useParentFeedPage'

const {
  posts,
  loading,
  errorMessage,
  currentPage,
  pagination,
  totalPages,
  fetchFeed,
  patchPostById,
  handlePageChange,
} = useParentFeedPage()
</script>

<template>
  <div class="parent-feed">
    <ParentFeedSummaryCard :total-posts="pagination.total || 0" />

    <div v-if="errorMessage" class="alert alert--error alert-row">
      <AlertCircle :size="16" />
      <span>{{ errorMessage }}</span>
      <button class="btn btn--outline btn--sm" type="button" @click="fetchFeed">Thử lại</button>
    </div>

    <div v-if="loading" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <div v-else-if="posts.length === 0 && !errorMessage" class="card empty-card">
      <MessageSquare :size="48" class="text-muted" />
      <h3>Chưa có bài đăng nào</h3>
      <p>Bảng tin sẽ hiển thị thông báo và cập nhật từ giáo viên.</p>
    </div>

    <div v-else class="posts-list">
      <PostCard
        v-for="post in posts"
        :key="post.post_id"
        :post="post"
        author-label="Giáo viên"
        audience="parent"
        :enable-share="false"
        @patch-post="patchPostById"
      />
    </div>

    <PaginationBar
      v-if="posts.length > 0"
      :current-page="currentPage"
      :total-pages="totalPages"
      :total-items="pagination.total"
      :limit="pagination.limit"
      @page-change="handlePageChange"
    />
  </div>
</template>

<style scoped>
.parent-feed {
  max-width: 48rem;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.alert-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  flex-wrap: wrap;
}

.loading-block {
  display: flex;
  justify-content: center;
  padding: 3rem 0;
}

.empty-card {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: var(--spacing-3);
  padding: 3rem var(--spacing-4);
  text-align: center;
  color: var(--color-text-muted);
}

.empty-card h3,
.empty-card p {
  margin: 0;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
