<script setup>
import ActionModal from '../../../components/ActionModal.vue'
import { GENDER_OPTIONS } from './studentPresentation'

defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
  formMode: {
    type: String,
    required: true,
  },
  formData: {
    type: Object,
    required: true,
  },
  selectedClassName: {
    type: String,
    default: '',
  },
  isSubmitting: {
    type: Boolean,
    default: false,
  },
  errorMessage: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['close', 'submit', 'update:formData'])
</script>

<template>
  <ActionModal
    :is-open="isOpen"
    :title="formMode === 'add' ? `Thêm học sinh — ${selectedClassName}` : 'Sửa thông tin học sinh'"
    @close="emit('close')"
  >
    <form class="modal-form" @submit.prevent="emit('submit')">
      <div v-if="errorMessage" class="alert alert--error">
        {{ errorMessage }}
      </div>

      <div class="form-group mb-0">
        <label class="form-label" for="studentName">Họ và tên</label>
        <input
          id="studentName"
          :value="formData.full_name"
          type="text"
          class="form-input"
          placeholder="Nhập họ và tên học sinh"
          :disabled="isSubmitting"
          required
          @input="emit('update:formData', { ...formData, full_name: $event.target.value })"
        />
      </div>

      <div class="modal-grid">
        <div class="form-group mb-0">
          <label class="form-label" for="studentDob">Ngày sinh</label>
          <input
            id="studentDob"
            :value="formData.dob"
            type="date"
            class="form-input"
            :disabled="isSubmitting"
            required
            @input="emit('update:formData', { ...formData, dob: $event.target.value })"
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="studentGender">Giới tính</label>
          <select
            id="studentGender"
            :value="formData.gender"
            class="form-input"
            :disabled="isSubmitting"
            @change="emit('update:formData', { ...formData, gender: $event.target.value })"
          >
            <option v-for="option in GENDER_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>

      <div class="modal-actions">
        <button type="button" class="btn btn--outline" :disabled="isSubmitting" @click="emit('close')">
          Hủy
        </button>
        <button type="submit" class="btn btn--primary" :disabled="isSubmitting">
          {{ isSubmitting ? 'Đang lưu...' : 'Lưu lại' }}
        </button>
      </div>
    </form>
  </ActionModal>
</template>

<style scoped>
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.modal-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

@media (min-width: 768px) {
  .modal-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
