<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { formatDateTime } from '../../helpers/dateFormatter'
import {
  AUDIT_ENTITY_OPTIONS,
  getAuditActorLabel,
  getAuditFriendlyAction,
  getAuditActionTone,
  getAuditEntitySummary,
  getAuditDetailLines,
} from '../../helpers/auditLogPresentation'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_PAGE_SIZE_OPTIONS,
  ADMIN_RETRY_BUTTON_TEXT,
} from '../../helpers/adminConfig'
import { useAuditLogs } from '../../composables/admin'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import '../../assets/css/view-switch.css'

const authStore = useAuthStore()
const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')

const {
  auditLogs,
  totalPages,
  currentPage,
  totalItems,
  pageSize,
  isLoading,
  errorMessage,
  searchQuery,
  actionQuery,
  selectedEntityType,
  fromDate,
  toDate,
  fetchAuditLogs,
  applyFilters,
  resetFilters,
  handlePageChange,
  handlePageSizeChange,
} = useAuditLogs(isSuperAdmin)

const viewMode = ref('friendly') // 'friendly' or 'raw'

onMounted(() => {
  fetchAuditLogs(1)
})
</script>

<template>
  <div class="admin-audit-logs page-stack">
    <div class="page-actions">
      <button
        class="btn btn--outline"
        type="button"
        :disabled="isLoading"
        @click="fetchAuditLogs(1)"
      >
        Làm mới
      </button>
    </div>

    <div v-if="!isSuperAdmin" class="card">
      <EmptyState
        title="Không có quyền truy cập"
        message="Chỉ Super Admin mới có thể xem nhật ký hệ thống."
      />
    </div>

    <template v-else>
      <div class="card filters-card">
        <div class="filters-grid">
          <div class="form-group mb-0">
            <label class="form-label" for="auditSearch">Tìm kiếm</label>
            <input
              id="auditSearch"
              v-model="searchQuery"
              type="search"
              class="form-input"
              placeholder="Tìm trong đường dẫn hoặc chi tiết..."
            />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="auditAction">Hành động</label>
            <input
              id="auditAction"
              v-model="actionQuery"
              type="search"
              class="form-input"
              placeholder="Ví dụ: appointments.book hoặc POST /admin/users"
            />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="entityType">Đối tượng</label>
            <select id="entityType" v-model="selectedEntityType" class="form-input">
              <option
                v-for="option in AUDIT_ENTITY_OPTIONS"
                :key="option.value || 'all'"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </select>
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="fromDate">Từ thời điểm</label>
            <input id="fromDate" v-model="fromDate" type="datetime-local" class="form-input" />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="toDate">Đến thời điểm</label>
            <input id="toDate" v-model="toDate" type="datetime-local" class="form-input" />
          </div>
        </div>

        <div class="filter-actions">
          <label class="page-size">
            <span class="page-size__label">Hiển thị</span>
            <select
              class="form-input page-size__select"
              :value="pageSize"
              @change="handlePageSizeChange(Number($event.target.value))"
            >
              <option v-for="size in ADMIN_PAGE_SIZE_OPTIONS" :key="size" :value="size">
                {{ size }}
              </option>
            </select>
          </label>

          <div class="filter-actions__buttons">
            <button class="btn btn--outline" type="button" @click="resetFilters">Xóa lọc</button>
            <button class="btn btn--primary" type="button" @click="applyFilters">Áp dụng</button>
          </div>
        </div>
      </div>

      <div v-if="errorMessage" class="alert alert--error">
        <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
        <p>{{ errorMessage }}</p>
        <button class="btn btn--outline mt-2" type="button" @click="fetchAuditLogs(currentPage)">
          {{ ADMIN_RETRY_BUTTON_TEXT }}
        </button>
      </div>

      <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu nhật ký..." />

      <div v-else-if="auditLogs.length === 0" class="card">
        <EmptyState
          title="Chưa có bản ghi nào"
          message="Hệ thống chưa ghi nhận hoạt động phù hợp với bộ lọc hiện tại."
        />
      </div>

      <template v-else>
        <div class="view-switch" style="margin-bottom: var(--spacing-4)">
          <button
            type="button"
            class="view-switch__btn"
            :class="{ 'view-switch__btn--active': viewMode === 'friendly' }"
            @click="viewMode = 'friendly'"
          >
            Xem dạng bảng
          </button>
          <button
            type="button"
            class="view-switch__btn"
            :class="{ 'view-switch__btn--active': viewMode === 'raw' }"
            @click="viewMode = 'raw'"
          >
            Xem dữ liệu gốc
          </button>
        </div>

        <div v-if="viewMode === 'friendly'" class="card desktop-table">
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>Thời gian</th>
                  <th>Người thực hiện</th>
                  <th>Hành động</th>
                  <th>Đối tượng</th>
                  <th>Chi tiết</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in auditLogs" :key="log.audit_log_id">
                  <td class="whitespace-nowrap">{{ formatDateTime(log.created_at) }}</td>
                  <td>{{ getAuditActorLabel(log) }}</td>
                  <td>
                    <span :class="getAuditActionTone(log)">
                      {{ getAuditFriendlyAction(log) }}
                    </span>
                  </td>
                  <td>{{ getAuditEntitySummary(log) }}</td>
                  <td>
                    <ul class="detail-list">
                      <li v-for="line in getAuditDetailLines(log)" :key="line">{{ line }}</li>
                    </ul>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="mobile-list">
          <article v-for="log in auditLogs" :key="log.audit_log_id" class="card audit-card">
            <div class="audit-card__head">
              <p class="audit-card__time">{{ formatDateTime(log.created_at) }}</p>
              <span :class="getAuditActionTone(log)">{{ getAuditFriendlyAction(log) }}</span>
            </div>

            <p class="audit-card__meta">Người thực hiện: {{ getAuditActorLabel(log) }}</p>
            <p class="audit-card__meta">Đối tượng: {{ getAuditEntitySummary(log) }}</p>

            <ul class="detail-list">
              <li v-for="line in getAuditDetailLines(log)" :key="line">{{ line }}</li>
            </ul>
          </article>
        </div>

        <div v-if="viewMode === 'raw'" class="card raw-view">
          <pre class="raw-json">{{ JSON.stringify(auditLogs, null, 2) }}</pre>
        </div>

        <PaginationBar
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-items="totalItems"
          :limit="pageSize"
          @page-change="handlePageChange"
        />
      </template>
    </template>
  </div>
</template>

<style scoped>
.page-stack,
.mobile-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

.filters-card,
.desktop-table,
.audit-card {
  padding: var(--spacing-4);
}

.filters-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.filter-actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  margin-top: var(--spacing-3);
}

.filter-actions__buttons {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.page-size {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.page-size__label {
  white-space: nowrap;
}

.page-size__select {
  width: 96px;
}

.table-responsive {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  padding: var(--spacing-3) var(--spacing-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
  vertical-align: top;
}

.table th {
  background-color: var(--color-background);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-transform: uppercase;
}

.table tbody tr:hover {
  background-color: var(--color-background);
}

.detail-list {
  margin: 0;
  padding-left: 1rem;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.detail-list li + li {
  margin-top: 0.3rem;
}

.whitespace-nowrap {
  white-space: nowrap;
}

.mobile-list {
  display: none;
}

.audit-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.audit-card__time,
.audit-card__meta {
  margin: 0;
}

.audit-card__time {
  font-weight: 700;
  color: var(--color-text);
}

.audit-card__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.mt-2 {
  margin-top: var(--spacing-2);
}

@media (min-width: 768px) {
  .filters-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .filter-actions {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

@media (min-width: 1200px) {
  .filters-grid {
    .mobile-list {
      display: none;
    }

    .raw-view {
      padding: 0;
      overflow: hidden;
    }

    .raw-json {
      margin: 0;
      padding: var(--spacing-4);
      background: var(--color-bg-code, #f6f8fa);
      color: var(--color-text);
      font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
      font-size: 0.875rem;
      line-height: 1.5;
      overflow-x: auto;
      white-space: pre;
    }

    grid-template-columns: repeat(5, minmax(0, 1fr));
  }
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }
}
</style>
