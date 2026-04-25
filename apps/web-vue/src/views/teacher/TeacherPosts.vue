<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { AlertCircle, LoaderCircle, MessageSquare, Plus, X } from 'lucide-vue-next'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { POST_SCOPE_LABELS, POST_TYPE_OPTIONS } from '../../helpers/postConfig'
import PostCard from '../../components/PostCard.vue'
import PaginationBar from '../../components/PaginationBar.vue'

const classes = ref([])
const selectedClassId = ref('')
const students = ref([])
const posts = ref([])
const loading = ref(true)
const loadingPosts = ref(false)
const errorMessage = ref('')
const showForm = ref(false)
const scopeType = ref('class')
const formStudentId = ref('')
const postType = ref('announcement')
const content = ref('')
const submitting = ref(false)
const formError = ref('')
const pagination = ref({ total: 0, limit: 20, offset: 0, has_more: false })
const currentPage = ref(1)

const currentOffset = computed(() => (currentPage.value - 1) * pagination.value.limit)
const totalPages = computed(() => Math.max(1, Math.ceil((pagination.value.total || 0) / pagination.value.limit)))

async function fetchBootstrap() {
  loading.value = true
  errorMessage.value = ''
  try {
    const classResponse = await teacherService.getMyClasses()
    classes.value = classResponse?.data ?? []
    if (classes.value.length > 0) {
      selectedClassId.value = classes.value[0].class_id
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lớp học'
  } finally {
    loading.value = false
  }
}

async function fetchStudents() {
  if (!selectedClassId.value) {
    students.value = []
    return
  }

  try {
    const response = await teacherService.getStudentsInClass(selectedClassId.value)
    students.value = response?.data ?? []
    if (!formStudentId.value && students.value.length > 0) {
      formStudentId.value = students.value[0].student_id
    }
  } catch {
    students.value = []
  }
}

async function fetchPosts() {
  if (!selectedClassId.value) {
    posts.value = []
    return
  }

  loadingPosts.value = true
  errorMessage.value = ''
  try {
    const response = await teacherService.getClassPosts(selectedClassId.value, {
      limit: pagination.value.limit,
      offset: currentOffset.value,
    })
    posts.value = response?.data ?? []
    if (response?.pagination) {
      pagination.value = response.pagination
    } else {
      pagination.value = {
        total: posts.value.length,
        limit: pagination.value.limit,
        offset: currentOffset.value,
        has_more: false,
      }
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải bài đăng'
  } finally {
    loadingPosts.value = false
  }
}

async function handleCreatePost() {
  formError.value = ''

  if (!selectedClassId.value) {
    formError.value = 'Vui lòng chọn lớp học'
    return
  }

  if (!content.value.trim()) {
    formError.value = 'Vui lòng nhập nội dung bài đăng'
    return
  }

  if (scopeType.value === 'student' && !formStudentId.value) {
    formError.value = 'Vui lòng chọn học sinh'
    return
  }

  submitting.value = true
  try {
    await teacherService.createPost({
      scope_type: scopeType.value,
      class_id: selectedClassId.value,
      student_id: scopeType.value === 'student' ? formStudentId.value : undefined,
      type: postType.value,
      content: content.value.trim(),
    })
    content.value = ''
    showForm.value = false
    currentPage.value = 1
    await fetchPosts()
  } catch (error) {
    formError.value = extractErrorMessage(error) || 'Không thể tạo bài đăng'
  } finally {
    submitting.value = false
  }
}

function patchPostById(postId, patch) {
  posts.value = posts.value.map((post) => (post.post_id === postId ? { ...post, ...patch } : post))
}

function handlePageChange(page) {
  currentPage.value = page
}

watch(selectedClassId, async () => {
  currentPage.value = 1
  await fetchStudents()
  await fetchPosts()
})

watch(currentPage, fetchPosts)

onMounted(fetchBootstrap)
</script>

<template>
  <div class="teacher-posts">
    <div v-if="loading" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <template v-else>
      <div class="composer-shell card">
        <div class="composer-main">
          <div class="avatar">GV</div>
          <button type="button" class="composer-toggle" @click="showForm = !showForm">
            {{ showForm ? 'Đóng khung soạn bài' : 'Bạn muốn chia sẻ điều gì với lớp hôm nay?' }}
          </button>
        </div>

        <select v-if="classes.length > 0" v-model="selectedClassId" class="form-input class-select">
          <option v-for="classInfo in classes" :key="classInfo.class_id" :value="classInfo.class_id">
            {{ classInfo.name }}
          </option>
        </select>

        <div class="composer-stats">
          <span class="badge badge--outline">{{ pagination.total }} bài đăng</span>
          <span class="badge badge--info">{{ POST_SCOPE_LABELS[scopeType] }}</span>
          <button type="button" class="btn btn--sm btn--outline" @click="showForm = !showForm">
            <X v-if="showForm" :size="16" />
            <Plus v-else :size="16" />
            {{ showForm ? 'Đóng' : 'Tạo bài' }}
          </button>
        </div>
      </div>

      <div v-if="errorMessage" class="alert alert--error">
        <AlertCircle :size="16" />
        {{ errorMessage }}
      </div>

      <div v-if="showForm" class="card form-card">
        <div class="card__header">
          <h2 class="section-title">Tạo bài đăng mới</h2>
        </div>
        <div class="card__body">
          <form class="flex-col gap-4" @submit.prevent="handleCreatePost">
            <div v-if="formError" class="alert alert--error">{{ formError }}</div>

            <div class="form-grid">
              <div class="form-group mb-0">
                <label class="form-label">Phạm vi</label>
                <select v-model="scopeType" class="form-input">
                  <option value="class">Cả lớp</option>
                  <option value="student">Từng HS</option>
                </select>
              </div>

              <div v-if="scopeType === 'student'" class="form-group mb-0">
                <label class="form-label">Học sinh</label>
                <select v-model="formStudentId" class="form-input">
                  <option v-for="student in students" :key="student.student_id" :value="student.student_id">
                    {{ student.full_name }}
                  </option>
                </select>
              </div>

              <div class="form-group mb-0">
                <label class="form-label">Loại bài</label>
                <select v-model="postType" class="form-input">
                  <option v-for="option in POST_TYPE_OPTIONS" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>

            <div class="form-group mb-0">
              <label class="form-label" for="postContent">Nội dung</label>
              <textarea
                id="postContent"
                v-model="content"
                class="form-input"
                rows="4"
                placeholder="Nhập nội dung bài đăng..."
                required
              ></textarea>
            </div>

            <div class="form-actions">
              <button type="submit" class="btn btn--primary" :disabled="submitting">
                <LoaderCircle v-if="submitting" class="spin mr-2" :size="16" />
                Đăng
              </button>
            </div>
          </form>
        </div>
      </div>

      <div v-if="loadingPosts" class="loading-block">
        <LoaderCircle class="spin text-muted" :size="32" />
      </div>

      <div v-else-if="posts.length === 0" class="card empty-card">
        <MessageSquare :size="48" class="text-muted" />
        <h3>Chưa có bài đăng nào</h3>
        <p>Hãy tạo bài đầu tiên để cập nhật thông tin cho lớp học.</p>
        <button class="btn btn--primary" @click="showForm = true">
          <Plus :size="16" />
          Tạo bài đăng
        </button>
      </div>

      <div v-else class="posts-list">
        <PostCard
          v-for="post in posts"
          :key="post.post_id"
          :post="post"
          audience="teacher"
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
    </template>
  </div>
</template>

<style scoped>
.teacher-posts {
  max-width: 48rem;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.composer-shell {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
}

.composer-main,
.composer-stats,
.form-actions {
  display: flex;
  align-items: center;
}

.composer-main {
  gap: var(--spacing-3);
}

.composer-stats {
  flex-wrap: wrap;
  gap: var(--spacing-2);
  border-top: 1px solid var(--color-border);
  padding-top: var(--spacing-3);
}

.avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-background);
  color: var(--color-text-muted);
  font-weight: 700;
}

.composer-toggle {
  min-height: 2.5rem;
  flex: 1;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: transparent;
  color: var(--color-text-muted);
  padding: 0 var(--spacing-4);
  text-align: left;
}

.class-select {
  width: 100%;
}

.form-grid {
  display: grid;
  gap: var(--spacing-4);
}

.form-actions {
  justify-content: flex-end;
}

.section-title {
  margin: 0;
  font-size: var(--font-size-lg);
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.empty-card,
.loading-block {
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-block {
  padding: 3rem 0;
}

.empty-card {
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

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 640px) {
  .composer-shell {
    display: grid;
    grid-template-columns: 1fr 14rem;
  }

  .composer-stats {
    grid-column: 1 / -1;
  }

  .form-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
