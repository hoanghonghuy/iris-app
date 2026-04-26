<script setup>
import { computed } from 'vue'
import ActionModal from '../ActionModal.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    required: true,
  },
  errorMessage: {
    type: String,
    default: '',
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
  isSuperAdmin: {
    type: Boolean,
    default: false,
  },
  schools: {
    type: Array,
    default: () => [],
  },
  classes: {
    type: Array,
    default: () => [],
  },
  students: {
    type: Array,
    default: () => [],
  },
  selectedSchoolId: {
    type: String,
    default: '',
  },
  selectedClassId: {
    type: String,
    default: '',
  },
  selectedStudentId: {
    type: String,
    default: '',
  },
  includeStudentSelector: {
    type: Boolean,
    default: false,
  },
  classLabel: {
    type: String,
    default: 'Chọn lớp',
  },
  classEmptyText: {
    type: String,
    default: 'Không có lớp',
  },
  studentLabel: {
    type: String,
    default: 'Chọn học sinh',
  },
  studentEmptyText: {
    type: String,
    default: 'Không có học sinh',
  },
  submitText: {
    type: String,
    required: true,
  },
  submittingText: {
    type: String,
    default: 'Đang gán...',
  },
})

const emit = defineEmits([
  'close',
  'submit',
  'update:selectedSchoolId',
  'update:selectedClassId',
  'update:selectedStudentId',
])

const canSubmit = computed(() => {
  if (!props.selectedClassId) {
    return false
  }

  if (props.includeStudentSelector && !props.selectedStudentId) {
    return false
  }

  return true
})
</script>

<template>
  <ActionModal :is-open="props.isOpen" :title="props.title" @close="emit('close')">
    <div class="modal-form">
      <div v-if="props.errorMessage" class="alert alert--error">
        {{ props.errorMessage }}
      </div>

      <div v-if="props.isSuperAdmin" class="form-group mb-0">
        <label class="form-label">Chọn trường</label>
        <select
          :value="props.selectedSchoolId"
          class="form-input"
          @change="(event) => emit('update:selectedSchoolId', event.target.value)"
        >
          <option
            v-for="school in props.schools"
            :key="school.school_id"
            :value="school.school_id"
          >
            {{ school.name }}
          </option>
        </select>
      </div>

      <div class="form-group mb-0">
        <label class="form-label">{{ props.classLabel }} <span class="text-danger">*</span></label>
        <select
          :value="props.selectedClassId"
          class="form-input"
          :disabled="props.classes.length === 0"
          @change="(event) => emit('update:selectedClassId', event.target.value)"
        >
          <option v-if="props.classes.length === 0" value="" disabled>
            {{ props.classEmptyText }}
          </option>
          <option
            v-for="classItem in props.classes"
            :key="classItem.class_id"
            :value="classItem.class_id"
          >
            {{ classItem.name }}
          </option>
        </select>
      </div>

      <div v-if="props.includeStudentSelector" class="form-group mb-0">
        <label class="form-label"
          >{{ props.studentLabel }} <span class="text-danger">*</span></label
        >
        <select
          :value="props.selectedStudentId"
          class="form-input"
          :disabled="props.students.length === 0"
          @change="(event) => emit('update:selectedStudentId', event.target.value)"
        >
          <option v-if="props.students.length === 0" value="" disabled>
            {{ props.studentEmptyText }}
          </option>
          <option
            v-for="student in props.students"
            :key="student.student_id"
            :value="student.student_id"
          >
            {{ student.full_name }}
          </option>
        </select>
      </div>

      <div class="modal-actions">
        <button
          class="btn btn--outline"
          type="button"
          :disabled="props.isLoading"
          @click="emit('close')"
        >
          Hủy
        </button>
        <button
          class="btn btn--primary"
          type="button"
          :disabled="props.isLoading || !canSubmit"
          @click="emit('submit')"
        >
          {{ props.isLoading ? props.submittingText : props.submitText }}
        </button>
      </div>
    </div>
  </ActionModal>
</template>

<style scoped>
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}
</style>
