<script setup>
import { ref, onMounted, watch } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDate } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'

const classes = ref([])
const selectedClassId = ref('')

const posts = ref([])
const isLoadingClasses = ref(true)
const isLoadingPosts = ref(false)
const isSubmitting = ref(false)
const errorMessage = ref('')
const postContent = ref('')

// Lấy danh sách lớp được phân công
const fetchClasses = async () => {
  try {
    const data = await teacherService.getMyClasses()
    classes.value = data.data || []
    if (classes.value.length > 0) {
      selectedClassId.value = classes.value[0].class_id
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoadingClasses.value = false
  }
}

// Lấy danh sách bài đăng
const fetchPosts = async () => {
  if (!selectedClassId.value) return
  
  isLoadingPosts.value = true
  errorMessage.value = ''
  
  try {
    const data = await teacherService.getClassPosts(selectedClassId.value, { limit: 20 })
    posts.value = data.data || []
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoadingPosts.value = false
  }
}

watch(selectedClassId, () => {
  if (selectedClassId.value) {
    fetchPosts()
  }
})

onMounted(() => {
  fetchClasses()
})

const handleCreatePost = async () => {
  if (!postContent.value.trim() || !selectedClassId.value) return
  
  isSubmitting.value = true
  try {
    await teacherService.createPost({
      scope_type: 'class',
      class_id: selectedClassId.value,
      type: 'activity',
      content: postContent.value
    })
    postContent.value = ''
    fetchPosts()
  } catch (error) {
    alert('Lỗi đăng bài: ' + extractErrorMessage(error))
  } finally {
    isSubmitting.value = false
  }
}

const toggleLike = async (post) => {
  try {
    const res = await teacherService.togglePostLike(post.post_id)
    post.liked_by_me = res.data.liked_by_me
    post.like_count = res.data.like_count
  } catch (error) {
    console.error(error)
  }
}
</script>

<template>
  <div class="teacher-activities">
    <div class="layout-grid">
      <!-- Cột trái: Lọc lớp & Đăng bài -->
      <div class="left-col">
        <div class="card mb-6 p-4">
          <div class="form-group mb-0">
            <label class="form-label" for="classFilter">Chọn lớp học</label>
            <select id="classFilter" v-model="selectedClassId" class="form-input" :disabled="isLoadingClasses">
              <option value="" disabled v-if="classes.length === 0">-- Không có lớp --</option>
              <option v-for="cls in classes" :key="cls.class_id" :value="cls.class_id">
                {{ cls.name }} ({{ cls.school_year }})
              </option>
            </select>
          </div>
        </div>

        <div class="card p-4" v-if="selectedClassId">
          <h3 class="font-bold mb-4">Tạo bài đăng mới</h3>
          <form @submit.prevent="handleCreatePost">
            <textarea 
              v-model="postContent" 
              class="form-input mb-4" 
              rows="4" 
              placeholder="Chia sẻ hoạt động hôm nay của lớp..."
              :disabled="isSubmitting"
              required
            ></textarea>
            
            <div class="flex justify-between items-center">
              <button type="button" class="btn btn--sm btn--outline" disabled title="Tính năng đăng ảnh sắp ra mắt">
                📷 Thêm ảnh
              </button>
              <button type="submit" class="btn btn--primary" :disabled="isSubmitting || !postContent.trim()">
                {{ isSubmitting ? 'Đang đăng...' : 'Đăng bài' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Cột phải: Timeline bài đăng -->
      <div class="right-col">
        <LoadingSpinner v-if="isLoadingClasses || isLoadingPosts" message="Đang tải bài đăng..." />
        
        <div v-else-if="errorMessage" class="p-4 bg-red-50 text-danger rounded border border-red-200">
          <p class="font-bold">Lỗi tải dữ liệu</p>
          <p>{{ errorMessage }}</p>
        </div>

        <div v-else-if="posts.length === 0">
          <EmptyState 
            title="Chưa có bài đăng nào" 
            message="Hãy là người đầu tiên chia sẻ hoạt động của lớp học này."
            icon="box"
          />
        </div>

        <div v-else class="posts-timeline flex-col gap-4">
          <div v-for="post in posts" :key="post.post_id" class="card p-4 post-card">
            <div class="flex items-center gap-3 mb-4">
              <div class="avatar-sm bg-primary text-white font-bold flex-center rounded-full">
                {{ post.author?.full_name?.charAt(0) || 'T' }}
              </div>
              <div>
                <p class="font-bold m-0">{{ post.author?.full_name || 'Giáo viên' }}</p>
                <p class="text-xs text-muted m-0">{{ formatDate(post.created_at) }}</p>
              </div>
            </div>
            
            <div class="post-content mb-4 text-sm whitespace-pre-line">
              {{ post.content }}
            </div>
            
            <div v-if="post.media_urls && post.media_urls.length > 0" class="post-media mb-4">
              <!-- Placeholder cho Media -->
              <div class="bg-gray-100 p-8 text-center text-muted rounded">
                [Đính kèm hình ảnh/video]
              </div>
            </div>
            
            <div class="post-stats flex justify-between text-xs text-muted pb-3 border-b border-gray-100 mb-3">
              <span>{{ post.like_count }} lượt thích</span>
              <span>{{ post.comment_count }} bình luận</span>
            </div>
            
            <div class="post-actions flex gap-2">
              <button 
                class="btn flex-1 btn--sm flex-center gap-2" 
                :class="post.liked_by_me ? 'text-primary bg-blue-50 border-none' : 'btn--outline border-none text-muted'"
                @click="toggleLike(post)"
              >
                <span v-if="post.liked_by_me">❤️ Đã thích</span>
                <span v-else>🤍 Thích</span>
              </button>
              <button class="btn btn--outline border-none flex-1 btn--sm flex-center gap-2 text-muted">
                💬 Bình luận
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Layout cục bộ */
.m-0 { margin-bottom: 0; margin-top: 0; }
.border-b { border-bottom: 1px solid var(--color-border); }

.layout-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-6);
  align-items: start;
}

@media (min-width: 1024px) {
  .layout-grid {
    grid-template-columns: 1fr 2fr;
  }
}

.avatar-sm {
  width: 40px;
  height: 40px;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.post-card {
  box-shadow: var(--shadow-sm);
}
</style>
