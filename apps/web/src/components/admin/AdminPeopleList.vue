<script setup>
import { computed } from 'vue'
import { Phone, X } from 'lucide-vue-next'
import EmptyState from '../common/EmptyState.vue'
import PaginationBar from '../common/PaginationBar.vue'

const props = defineProps({
  items: {
    type: Array,
    required: true,
  },
  filteredItems: {
    type: Array,
    required: true,
  },
  searchQuery: {
    type: String,
    default: '',
  },
  searchPlaceholder: {
    type: String,
    default: 'Tìm theo tên, email, SĐT...',
  },
  emptyTitle: {
    type: String,
    required: true,
  },
  emptyMessage: {
    type: String,
    required: true,
  },
  emptySearchLabel: {
    type: String,
    required: true,
  },
  itemKeyField: {
    type: String,
    required: true,
  },
  nameField: {
    type: String,
    default: 'full_name',
  },
  emailField: {
    type: String,
    default: 'email',
  },
  phoneField: {
    type: String,
    default: 'phone',
  },
  relationField: {
    type: String,
    required: true,
  },
  relationKeyField: {
    type: String,
    required: true,
  },
  relationNameField: {
    type: String,
    required: true,
  },
  relationColumnTitle: {
    type: String,
    required: true,
  },
  actionColumnTitle: {
    type: String,
    required: true,
  },
  actionColumnWidth: {
    type: Number,
    default: 220,
  },
  noRelationText: {
    type: String,
    required: true,
  },
  removeRelationTitle: {
    type: String,
    required: true,
  },
  currentPage: {
    type: Number,
    required: true,
  },
  totalPages: {
    type: Number,
    required: true,
  },
  totalItems: {
    type: Number,
    required: true,
  },
  limit: {
    type: Number,
    required: true,
  },
})

const emit = defineEmits(['update:searchQuery', 'edit', 'assign', 'remove-relation', 'page-change'])

const localSearchQuery = computed({
  get: () => props.searchQuery,
  set: (value) => emit('update:searchQuery', value),
})

function getItemValue(item, field, fallback = '-') {
  return item?.[field] || fallback
}

function getRelationList(item) {
  const relations = item?.[props.relationField]
  return Array.isArray(relations) ? relations : []
}

function getRelationKey(relation) {
  return relation?.[props.relationKeyField] || relation?.[props.relationNameField]
}
</script>

<template>
  <div class="page-stack">
    <div v-if="items.length > 0" class="card toolbar-card">
      <div class="toolbar-grid">
        <input
          v-model="localSearchQuery"
          type="search"
          class="form-input"
          :placeholder="searchPlaceholder"
        />
      </div>
    </div>

    <div v-if="items.length === 0" class="card">
      <EmptyState :title="emptyTitle" :message="emptyMessage" />
    </div>

    <div v-else-if="filteredItems.length === 0" class="card empty-search">
      Không tìm thấy {{ emptySearchLabel }} nào phù hợp với "{{ searchQuery }}"
    </div>

    <template v-else>
      <div class="card desktop-table">
        <div class="table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th>Họ tên</th>
                <th>Email</th>
                <th>{{ relationColumnTitle }}</th>
                <th
                  class="action-column"
                  :style="{ width: `${actionColumnWidth}px`, minWidth: `${actionColumnWidth}px` }"
                >
                  {{ actionColumnTitle }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredItems" :key="item[itemKeyField]">
                <td>
                  <div class="font-medium">{{ getItemValue(item, nameField, 'Chưa cập nhật') }}</div>
                  <div class="text-xs text-muted mt-1 phone-line">
                    <Phone :size="12" />
                    <span>{{ getItemValue(item, phoneField) }}</span>
                  </div>
                </td>
                <td class="text-muted">{{ getItemValue(item, emailField) }}</td>
                <td>
                  <template v-if="getRelationList(item).length > 0">
                    <div class="flex flex-wrap gap-1">
                      <span
                        v-for="relation in getRelationList(item)"
                        :key="getRelationKey(relation)"
                        class="badge badge--outline badge--sm flex items-center gap-1"
                      >
                        {{ relation[relationNameField] }}
                        <button
                          class="badge-remove-btn"
                          type="button"
                          :title="removeRelationTitle"
                          @click="emit('remove-relation', { item, relation })"
                        >
                          <X :size="11" />
                        </button>
                      </span>
                    </div>
                  </template>
                  <span v-else class="text-muted text-sm italic">{{ noRelationText }}</span>
                </td>
                <td
                  class="action-column"
                  :style="{ width: `${actionColumnWidth}px`, minWidth: `${actionColumnWidth}px` }"
                >
                  <div class="table-action-buttons">
                    <slot name="desktop-actions" :item="item" />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="mobile-list">
        <article v-for="item in filteredItems" :key="item[itemKeyField]" class="card mobile-card">
          <div class="mobile-card__head">
            <p class="mobile-card__title">{{ getItemValue(item, nameField, 'Chưa cập nhật') }}</p>
            <slot name="mobile-head-extra" :item="item" />
          </div>

          <p class="mobile-card__meta">{{ getItemValue(item, emailField) }}</p>
          <p class="mobile-card__meta mobile-card__phone">
            <Phone :size="12" />
            <span>{{ getItemValue(item, phoneField) }}</span>
          </p>

          <div v-if="getRelationList(item).length > 0" class="mobile-card__chips">
            <span
              v-for="relation in getRelationList(item)"
              :key="getRelationKey(relation)"
              class="badge badge--outline badge--sm flex items-center gap-1"
            >
              {{ relation[relationNameField] }}
              <button
                class="badge-remove-btn"
                type="button"
                :title="removeRelationTitle"
                @click="emit('remove-relation', { item, relation })"
              >
                <X :size="11" />
              </button>
            </span>
          </div>
          <p v-else class="mobile-card__meta italic">{{ noRelationText }}</p>

          <div class="mobile-card__actions">
            <slot name="mobile-actions" :item="item" />
          </div>
        </article>
      </div>

      <PaginationBar
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :limit="limit"
        @page-change="(page) => emit('page-change', page)"
      />
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

.toolbar-card,
.desktop-table,
.mobile-card {
  padding: var(--spacing-4);
}

.toolbar-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: minmax(0, 1fr);
}

.empty-search {
  padding: var(--spacing-6);
  text-align: center;
  color: var(--color-text-muted);
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
  font-weight: 600;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  text-transform: uppercase;
}

.table tbody tr:hover {
  background-color: var(--color-background);
}

.action-column {
  text-align: right !important;
  white-space: nowrap;
}

.table-action-buttons {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.table-action-buttons .btn,
.mobile-card__head .btn,
.mobile-card__actions .btn {
  gap: 0.35rem;
}

.phone-line,
.mobile-card__phone {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

.badge--sm {
  font-size: 0.7rem;
  padding: 2px 6px;
}

.badge-remove-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-muted);
  padding: 0;
  line-height: 1;
}

.badge-remove-btn:hover {
  color: var(--color-danger);
}

.mobile-list {
  display: none;
}

.mobile-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
}

.mobile-card__head {
  display: flex;
  align-items: center;
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

.mobile-card__chips,
.mobile-card__actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.mobile-card__actions {
  justify-content: flex-end;
}

.mt-1 {
  margin-top: var(--spacing-1);
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