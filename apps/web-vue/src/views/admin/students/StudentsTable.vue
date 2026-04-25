<script setup>
import { Check, Copy, KeyRound, LoaderCircle, Pencil, Trash2, X } from 'lucide-vue-next'
import { formatDate } from '../../../helpers/dateFormatter'
import { getCodeExpiryText, getGenderLabel } from './studentPresentation'

defineProps({
  students: {
    type: Array,
    required: true,
  },
  generatingCodeStudentId: {
    type: String,
    default: '',
  },
  revokingCodeStudentId: {
    type: String,
    default: '',
  },
  copiedStudentId: {
    type: String,
    default: '',
  },
})

const emit = defineEmits([
  'edit',
  'delete',
  'generate-code',
  'copy-code',
  'revoke-code',
])
</script>

<template>
  <template v-if="students.length > 0">
    <div class="card desktop-table">
      <div class="table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Họ và tên</th>
              <th>Ngày sinh</th>
              <th>Giới tính</th>
              <th class="code-column">Mã PH</th>
              <th class="action-column">Hành động</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in students" :key="student.student_id">
              <td class="font-medium">
                <RouterLink :to="`/admin/students/${student.student_id}`" class="student-link">
                  {{ student.full_name }}
                </RouterLink>
              </td>
              <td>{{ formatDate(student.dob) }}</td>
              <td>
                <span class="badge badge--outline">{{ getGenderLabel(student.gender) }}</span>
              </td>
              <td class="code-column">
                <div v-if="student.active_parent_code" class="code-cell">
                  <div class="code-row">
                    <code class="code-pill">{{ student.active_parent_code }}</code>
                    <button
                      class="icon-button"
                      type="button"
                      :title="copiedStudentId === student.student_id ? 'Đã sao chép' : 'Sao chép mã'"
                      :aria-label="copiedStudentId === student.student_id ? 'Đã sao chép mã phụ huynh' : 'Sao chép mã phụ huynh'"
                      @click="emit('copy-code', student.active_parent_code, student.student_id)"
                    >
                      <Check v-if="copiedStudentId === student.student_id" :size="14" />
                      <Copy v-else :size="14" />
                    </button>
                    <button
                      class="icon-button icon-button--danger"
                      type="button"
                      title="Thu hồi mã"
                      aria-label="Thu hồi mã phụ huynh"
                      :disabled="revokingCodeStudentId === student.student_id"
                      @click="emit('revoke-code', student)"
                    >
                      <LoaderCircle v-if="revokingCodeStudentId === student.student_id" :size="14" class="spin" />
                      <X v-else :size="14" />
                    </button>
                    <span class="code-expiry" :class="{ 'code-expiry--expired': getCodeExpiryText(student.code_expires_at) === 'Hết hạn' }">
                      {{ getCodeExpiryText(student.code_expires_at) }}
                    </span>
                  </div>
                </div>

                <div v-else class="code-row code-row--create">
                  <button
                    class="btn btn--sm btn--outline compact-generate"
                    type="button"
                    :disabled="generatingCodeStudentId === student.student_id"
                    @click="emit('generate-code', student)"
                  >
                    <LoaderCircle v-if="generatingCodeStudentId === student.student_id" :size="14" class="spin" />
                    <KeyRound v-else :size="14" />
                    <span>{{ generatingCodeStudentId === student.student_id ? 'Đang tạo...' : 'Tạo mã' }}</span>
                  </button>
                </div>
              </td>
              <td class="action-column">
                <div class="table-actions">
                  <button
                    class="icon-button"
                    type="button"
                    title="Sửa học sinh"
                    aria-label="Sửa học sinh"
                    @click="emit('edit', student)"
                  >
                    <Pencil :size="14" />
                  </button>
                  <button
                    class="icon-button icon-button--danger"
                    type="button"
                    title="Xóa học sinh"
                    aria-label="Xóa học sinh"
                    @click="emit('delete', student)"
                  >
                    <Trash2 :size="14" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="mobile-list">
      <article v-for="student in students" :key="student.student_id" class="mobile-card">
        <div class="mobile-card__head">
          <div>
            <RouterLink :to="`/admin/students/${student.student_id}`" class="student-link mobile-card__title">
              {{ student.full_name }}
            </RouterLink>
            <p class="mobile-card__meta">
              {{ formatDate(student.dob) }} • {{ getGenderLabel(student.gender) }}
            </p>
          </div>
          <button class="btn btn--sm btn--outline" type="button" @click="emit('edit', student)">
            Sửa
          </button>
        </div>

        <div v-if="student.active_parent_code" class="mobile-code">
          <div class="mobile-code__row">
            <code class="code-pill">{{ student.active_parent_code }}</code>
            <span class="code-expiry" :class="{ 'code-expiry--expired': getCodeExpiryText(student.code_expires_at) === 'Hết hạn' }">
              {{ getCodeExpiryText(student.code_expires_at) }}
            </span>
          </div>

          <div class="mobile-card__actions">
            <button class="btn btn--sm btn--outline" type="button" @click="emit('copy-code', student.active_parent_code, student.student_id)">
              {{ copiedStudentId === student.student_id ? 'Đã chép' : 'Sao chép' }}
            </button>
            <button
              class="btn btn--sm btn--danger"
              type="button"
              :disabled="revokingCodeStudentId === student.student_id"
              @click="emit('revoke-code', student)"
            >
              {{ revokingCodeStudentId === student.student_id ? 'Đang thu hồi...' : 'Thu hồi' }}
            </button>
          </div>
        </div>

        <div v-else class="mobile-card__actions">
          <button
            class="btn btn--sm btn--outline compact-generate"
            type="button"
            :disabled="generatingCodeStudentId === student.student_id"
            @click="emit('generate-code', student)"
          >
            <LoaderCircle v-if="generatingCodeStudentId === student.student_id" :size="14" class="spin" />
            <KeyRound v-else :size="14" />
            <span>{{ generatingCodeStudentId === student.student_id ? 'Đang tạo...' : 'Tạo mã PH' }}</span>
          </button>
        </div>

        <div class="mobile-card__actions mobile-card__actions--end">
          <button class="btn btn--sm btn--danger" type="button" @click="emit('delete', student)">
            Xóa
          </button>
        </div>
      </article>
    </div>
  </template>
</template>

<style scoped>
.desktop-table,
.mobile-card {
  padding: var(--spacing-4);
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
  vertical-align: middle;
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

.student-link {
  color: var(--color-primary);
  text-decoration: none;
}

.student-link:hover {
  text-decoration: underline;
}

.code-column {
  width: 240px;
  min-width: 240px;
  text-align: right;
  white-space: nowrap;
}

.action-column {
  width: 96px;
  min-width: 96px;
  text-align: right;
  white-space: nowrap;
}

.table-actions,
.mobile-card__actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.table-actions,
.mobile-card__actions--end {
  justify-content: flex-end;
}

.code-cell {
  display: flex;
  justify-content: flex-end;
}

.code-row {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.code-row--create {
  width: 100%;
}

.code-pill {
  display: inline-flex;
  align-items: center;
  min-height: 2rem;
  padding: 0.25rem 0.6rem;
  border-radius: var(--radius-md);
  background: var(--color-background);
  border: 1px solid var(--color-border);
  font-size: var(--font-size-sm);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Courier New', monospace;
}

.code-expiry {
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.code-expiry--expired {
  color: var(--color-danger);
  font-weight: 600;
}

.compact-generate {
  gap: 0.4rem;
}

.icon-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease;
}

.icon-button:hover:not(:disabled) {
  background: var(--color-background);
  border-color: var(--color-text-muted);
  color: var(--color-text);
}

.icon-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon-button--danger {
  color: var(--color-danger);
}

.icon-button--danger:hover:not(:disabled) {
  background: var(--color-danger-soft-bg);
  border-color: var(--color-danger-soft-border);
  color: var(--color-danger-hover);
}

.spin {
  animation: spin 1s linear infinite;
}

.mobile-list {
  display: none;
}

.mobile-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
}

.mobile-card__head,
.mobile-code__row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.mobile-card__title,
.mobile-card__meta {
  margin: 0;
}

.mobile-card__title {
  color: var(--color-text);
  font-weight: 700;
}

.mobile-card__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.mobile-code {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  margin-top: var(--spacing-3);
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 1024px) {
  .code-column {
    width: 220px;
    min-width: 220px;
  }
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-4);
  }
}
</style>
