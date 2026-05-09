<script setup>
import { AlertCircle, LoaderCircle, MessageSquare } from 'lucide-vue-next'
import PostCard from '../../components/PostCard.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import ParentFeedSummaryCard from './feed/ParentFeedSummaryCard.vue'
import { useParentFeedPage } from '../../composables/parent'

const {
  posts,
  children,
  selectedChildId,
  feedMode,
  loading,
  errorMessage,
  currentPage,
  pagination,
  totalPages,
  fetchFeed,
  setFeedMode,
  setSelectedChild,
  patchPostById,
  handlePageChange,
} = useParentFeedPage()
</script>

<template>
  <div class="parent-feed">
    <ParentFeedSummaryCard :total-posts="pagination.total || 0" />

    <div class="card filter-card">
      <div class="form-group mb-0">
        <label class="form-label" for="feedMode">Nguồn dữ liệu</label>
        <select id="feedMode" :value="feedMode" class="form-input" @change="setFeedMode($event.target.value)">
          <option value="feed">Feed tổng hợp</option>
          <option value="child_all">Theo từng con (tất cả bài)</option>
          <option value="child_class">Theo từng con (bài lớp)</option>
          <option value="child_student">Theo từng con (bài cá nhân)</option>
        </select>
      </div>

      <div v-if="feedMode !== 'feed'" class="form-group mb-0">
        <label class="form-label" for="childFilter">Chọn học sinh</label>
        <select
          id="childFilter"
          :value="selectedChildId"
          class="form-input"
          :disabled="children.length === 0"
          @change="setSelectedChild($event.target.value)"
        >
          <option v-if="children.length === 0" value="">Không có dữ liệu học sinh</option>
          <option
            v-for="child in children"
            :key="child.student_id"
            :value="child.student_id"
          >
            {{ child.full_name }}
          </option>
        </select>
      </div>
    </div>

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

.filter-card {
  display: grid;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
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

@media (min-width: 640px) {
  .filter-card {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
