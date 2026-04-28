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
  formData: {
    type: Object,
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
  fullNameInputId: {
    type: String,
    required: true,
  },
  phoneInputId: {
    type: String,
    required: true,
  },
  schoolInputId: {
    type: String,
    required: true,
  },
  fullNamePlaceholder: {
    type: String,
    default: 'Ví dụ: Nguyễn Văn A',
  },
  submitText: {
    type: String,
    default: 'Lưu lại',
  },
  submittingText: {
    type: String,
    default: 'Đang lưu...',
  },
})

const emit = defineEmits(['close', 'submit', 'update:formData'])

const fullName = computed({
  get: () => props.formData?.full_name || '',
  set: (value) => {
    emit('update:formData', {
      ...props.formData,
      full_name: value,
    })
  },
})

const phone = computed({
  get: () => props.formData?.phone || '',
  set: (value) => {
    emit('update:formData', {
      ...props.formData,
      phone: value,
    })
  },
})

const schoolId = computed({
  get: () => props.formData?.school_id || '',
  set: (value) => {
    emit('update:formData', {
      ...props.formData,
      school_id: value,
    })
  },
})
</script>

<template>
  <ActionModal :is-open="isOpen" :title="title" @close="emit('close')">
    <form class="modal-form" @submit.prevent="emit('submit')">
      <div v-if="errorMessage" class="alert alert--error">
        {{ errorMessage }}
      </div>

      <div class="form-group mb-0">
        <label class="form-label" :for="fullNameInputId"
          >Họ và tên <span class="text-danger">*</span></label
        >
        <input
          :id="fullNameInputId"
          v-model="fullName"
          type="text"
          class="form-input"
          :placeholder="fullNamePlaceholder"
          :disabled="isLoading"
          required
        />
      </div>

      <div class="form-group mb-0">
        <label class="form-label" :for="phoneInputId">Số điện thoại</label>
        <input
          :id="phoneInputId"
          v-model="phone"
          type="text"
          class="form-input"
          placeholder="Nhập số điện thoại"
          :disabled="isLoading"
        />
      </div>

      <div v-if="isSuperAdmin" class="form-group mb-0">
        <label class="form-label" :for="schoolInputId"
          >Trường học <span class="text-danger">*</span></label
        >
        <select
          :id="schoolInputId"
          v-model="schoolId"
          class="form-input"
          :disabled="isLoading"
          required
        >
          <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
            {{ school.name }}
          </option>
        </select>
      </div>

      <div class="modal-actions">
        <button class="btn btn--outline" type="button" :disabled="isLoading" @click="emit('close')">
          Hủy
        </button>
        <button class="btn btn--primary" type="submit" :disabled="isLoading">
          {{ isLoading ? submittingText : submitText }}
        </button>
      </div>
    </form>
  </ActionModal>
</template>
