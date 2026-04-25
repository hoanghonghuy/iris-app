<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { AlertCircle, Check, History, LoaderCircle } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateTimeVN, formatDateVN } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'

const route = useRoute()

const ATTENDANCE_STATUS_OPTIONS = [
  { value: 'present', label: 'Có mặt', badge: 'badge badge--primary' },
  { value: 'absent', label: 'Vắng', badge: 'badge badge--danger' },
  { value: 'late', label: 'Muộn', badge: 'badge badge--warning' },
  { value: 'excused', label: 'Có phép', badge: 'badge badge--outline' },
]

const today = getDateInputValue(new Date())

const classes = ref([])
const selectedClassId = ref('')
const students = ref([])

const attendanceData = ref({})
const savedAttendance = ref({})
const hasSavedForDate = ref({})

const isLoadingClasses = ref(true)
const isLoadingStudents = ref(false)
const savingRowId = ref('')
const cancelingRowId = ref('')
const isSavingAll = ref(false)
const isSavingDisplayed = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const viewMode = ref('take')
const studentSearch = ref('')
const takeListFilter = ref('all')
const listOrderMode = ref('prioritize')
const showMobileTakeControls = ref(false)

const historyOpen = ref(new Set())
const historyLoading = ref(new Set())
const historyByStudent = ref({})

const historyFrom = ref(getDateInputValue(daysAgo(7)))
const historyTo = ref(today)
const historyStudentId = ref('all')
const historyStatus = ref('all')
const historyListLoading = ref(false)
const historyList = ref([])
const historyOffset = ref(0)
const historyLimit = ref(20)
const historyTotal = ref(0)
const historyHasMore = ref(false)

function daysAgo(count) {
  const date = new Date()
  date.setDate(date.getDate() - count)
  return date
}

function getDateInputValue(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

function unwrapList(value) {
  const data = value?.data ?? value
  return Array.isArray(data) ? data.filter(Boolean) : []
}

function getStatusOption(status) {
  return ATTENDANCE_STATUS_OPTIONS.find((option) => option.value === status)
}

function getStatusLabel(status) {
  return getStatusOption(status)?.label || status || '-'
}

function getDefaultAttendanceValue() {
  return { status: 'present', note: '' }
}

function isSameAttendance(left, right) {
  return left?.status === right?.status && (left?.note || '') === (right?.note || '')
}

function isRowDirty(studentId) {
  const currentValue = attendanceData.value[studentId]
  const savedValue = savedAttendance.value[studentId]

  if (!currentValue || !savedValue) {
    return !hasSavedForDate.value[studentId]
  }

  if (!hasSavedForDate.value[studentId]) {
    return true
  }

  return !isSameAttendance(currentValue, savedValue)
}

function resetStudentHistoryState() {
  historyOpen.value = new Set()
  historyLoading.value = new Set()
  historyByStudent.value = {}
}

const displayedStudents = computed(() => {
  const normalizedSearch = studentSearch.value.trim().toLowerCase()
  const searchedStudents = normalizedSearch
    ? students.value.filter((student) => student.full_name?.toLowerCase().includes(normalizedSearch))
    : students.value

  let filtered = searchedStudents
  if (takeListFilter.value === 'pending') {
    filtered = searchedStudents.filter((student) => isRowDirty(student.student_id))
  } else if (takeListFilter.value === 'saved') {
    filtered = searchedStudents.filter((student) => !isRowDirty(student.student_id))
  }

  if (listOrderMode.value === 'original') {
    return filtered
  }

  const unsavedStudents = filtered.filter((student) => isRowDirty(student.student_id))
  const savedStudents = filtered.filter((student) => !isRowDirty(student.student_id))
  return [...unsavedStudents, ...savedStudents]
})

const dirtyCount = computed(() => students.value.filter((student) => isRowDirty(student.student_id)).length)
const displayedDirtyCount = computed(() => displayedStudents.value.filter((student) => isRowDirty(student.student_id)).length)
const displayedSavedCount = computed(() => displayedStudents.value.length - displayedDirtyCount.value)
const globalPendingCount = computed(() => students.value.filter((student) => isRowDirty(student.student_id)).length)

async function fetchClasses() {
  isLoadingClasses.value = true
  errorMessage.value = ''

  try {
    classes.value = unwrapList(await teacherService.getMyClasses())

    const requestedClassId = typeof route.query.classId === 'string' ? route.query.classId : ''
    const nextClassId = classes.value.some((item) => item.class_id === requestedClassId)
      ? requestedClassId
      : classes.value[0]?.class_id || ''

    if (!classes.value.some((item) => item.class_id === selectedClassId.value)) {
      selectedClassId.value = nextClassId
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải lớp'
  } finally {
    isLoadingClasses.value = false
  }
}

async function fetchStudentsAndAttendance() {
  if (!selectedClassId.value) {
    students.value = []
    attendanceData.value = {}
    savedAttendance.value = {}
    hasSavedForDate.value = {}
    resetStudentHistoryState()
    return
  }

  isLoadingStudents.value = true
  errorMessage.value = ''
  successMessage.value = ''
  resetStudentHistoryState()

  try {
    const studentList = unwrapList(await teacherService.getStudentsInClass(selectedClassId.value))
    students.value = studentList

    const initialAttendance = {}
    const initialSavedAttendance = {}
    const initialHasSaved = {}

    await Promise.all(
      studentList.map(async (student) => {
        const fallback = getDefaultAttendanceValue()

        try {
          const records = unwrapList(
            await teacherService.getStudentAttendance(student.student_id, today, today),
          )

          const existingRecord = records.find((record) => String(record.date || '').slice(0, 10) === today)
          if (existingRecord) {
            const savedValue = {
              status: existingRecord.status || fallback.status,
              note: existingRecord.note || '',
            }
            initialAttendance[student.student_id] = savedValue
            initialSavedAttendance[student.student_id] = savedValue
            initialHasSaved[student.student_id] = true
            return
          }
        } catch {
          initialAttendance[student.student_id] = { ...fallback }
          initialSavedAttendance[student.student_id] = { ...fallback }
          initialHasSaved[student.student_id] = false
          return
        }

        initialAttendance[student.student_id] = { ...fallback }
        initialSavedAttendance[student.student_id] = { ...fallback }
        initialHasSaved[student.student_id] = false
      }),
    )

    attendanceData.value = initialAttendance
    savedAttendance.value = initialSavedAttendance
    hasSavedForDate.value = initialHasSaved
    historyStudentId.value = 'all'
  } catch (error) {
    students.value = []
    attendanceData.value = {}
    savedAttendance.value = {}
    hasSavedForDate.value = {}
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách học sinh'
  } finally {
    isLoadingStudents.value = false
  }
}

function handleAttendanceStatusChange(studentId, status) {
  const currentValue = attendanceData.value[studentId] || getDefaultAttendanceValue()
  attendanceData.value = {
    ...attendanceData.value,
    [studentId]: {
      ...currentValue,
      status,
    },
  }
}

function handleAttendanceNoteChange(studentId, note) {
  const currentValue = attendanceData.value[studentId] || getDefaultAttendanceValue()
  attendanceData.value = {
    ...attendanceData.value,
    [studentId]: {
      ...currentValue,
      note,
    },
  }
}

async function loadClassHistory(offset = 0) {
  if (!selectedClassId.value || students.value.length === 0) {
    historyList.value = []
    historyTotal.value = 0
    historyHasMore.value = false
    return
  }

  historyListLoading.value = true
  errorMessage.value = ''

  try {
    const response = await teacherService.getClassAttendanceChanges(selectedClassId.value, {
      from: historyFrom.value || undefined,
      to: historyTo.value || undefined,
      student_id: historyStudentId.value === 'all' ? undefined : historyStudentId.value,
      status: historyStatus.value === 'all' ? undefined : historyStatus.value,
      limit: historyLimit.value,
      offset,
    })

    const items = unwrapList(response).map((item) => ({
      ...item,
      student_name:
        item.student_name
        || students.value.find((student) => student.student_id === item.student_id)?.full_name
        || 'Không rõ',
    }))

    historyList.value = items
    historyTotal.value = response?.pagination?.total || items.length
    historyHasMore.value = Boolean(response?.pagination?.has_more)
    historyOffset.value = offset
  } catch (error) {
    historyList.value = []
    historyTotal.value = 0
    historyHasMore.value = false
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải lịch sử điểm danh'
  } finally {
    historyListLoading.value = false
  }
}

async function toggleHistory(studentId) {
  const nextOpen = new Set(historyOpen.value)
  const shouldOpen = !nextOpen.has(studentId)

  if (shouldOpen) {
    nextOpen.add(studentId)
  } else {
    nextOpen.delete(studentId)
  }
  historyOpen.value = nextOpen

  if (!shouldOpen || historyByStudent.value[studentId]) {
    return
  }

  const nextLoading = new Set(historyLoading.value)
  nextLoading.add(studentId)
  historyLoading.value = nextLoading

  try {
    const formattedFromDate = getDateInputValue(daysAgo(30))
    const records = unwrapList(
      await teacherService.getStudentAttendanceChanges(studentId, formattedFromDate, today),
    )
    historyByStudent.value = {
      ...historyByStudent.value,
      [studentId]: records,
    }
  } catch {
    historyByStudent.value = {
      ...historyByStudent.value,
      [studentId]: [],
    }
  } finally {
    const finishedLoading = new Set(historyLoading.value)
    finishedLoading.delete(studentId)
    historyLoading.value = finishedLoading
  }
}

async function handleMark(studentId) {
  const currentAttendance = attendanceData.value[studentId]
  if (!currentAttendance) return

  savingRowId.value = studentId
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await teacherService.markAttendance({
      student_id: studentId,
      date: today,
      status: currentAttendance.status,
      note: currentAttendance.note || '',
    })

    savedAttendance.value = {
      ...savedAttendance.value,
      [studentId]: { ...currentAttendance },
    }
    hasSavedForDate.value = {
      ...hasSavedForDate.value,
      [studentId]: true,
    }
    successMessage.value = 'Đã lưu điểm danh thành công.'
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể lưu điểm danh'
  } finally {
    savingRowId.value = ''
  }
}

function handleRevert(studentId) {
  const savedValue = savedAttendance.value[studentId]
  if (!savedValue) return

  attendanceData.value = {
    ...attendanceData.value,
    [studentId]: { ...savedValue },
  }
}

async function handleCancelSaved(studentId) {
  if (!hasSavedForDate.value[studentId]) return

  cancelingRowId.value = studentId
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await teacherService.cancelAttendance(studentId, today)
    const fallback = getDefaultAttendanceValue()

    attendanceData.value = {
      ...attendanceData.value,
      [studentId]: { ...fallback },
    }
    savedAttendance.value = {
      ...savedAttendance.value,
      [studentId]: { ...fallback },
    }
    hasSavedForDate.value = {
      ...hasSavedForDate.value,
      [studentId]: false,
    }

    const nextHistory = { ...historyByStudent.value }
    delete nextHistory[studentId]
    historyByStudent.value = nextHistory

    successMessage.value = 'Đã hủy điểm danh đã lưu.'
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể hủy điểm danh đã lưu'
  } finally {
    cancelingRowId.value = ''
  }
}

async function saveStudents(studentList, savingMode) {
  if (studentList.length === 0) return

  if (savingMode === 'displayed') {
    isSavingDisplayed.value = true
  } else {
    isSavingAll.value = true
  }

  errorMessage.value = ''
  successMessage.value = ''

  try {
    await Promise.all(
      studentList.map((student) => {
        const currentAttendance = attendanceData.value[student.student_id]
        if (!currentAttendance) {
          return Promise.resolve()
        }

        return teacherService.markAttendance({
          student_id: student.student_id,
          date: today,
          status: currentAttendance.status,
          note: currentAttendance.note || '',
        })
      }),
    )

    const nextSavedAttendance = { ...savedAttendance.value }
    const nextHasSaved = { ...hasSavedForDate.value }

    studentList.forEach((student) => {
      nextSavedAttendance[student.student_id] = { ...attendanceData.value[student.student_id] }
      nextHasSaved[student.student_id] = true
    })

    savedAttendance.value = nextSavedAttendance
    hasSavedForDate.value = nextHasSaved
    successMessage.value = savingMode === 'displayed'
      ? 'Đã lưu danh sách đang hiển thị.'
      : 'Đã lưu điểm danh cho toàn lớp.'
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể lưu điểm danh hàng loạt'
  } finally {
    isSavingDisplayed.value = false
    isSavingAll.value = false
  }
}

async function handleSaveDisplayed() {
  const dirtyStudents = displayedStudents.value.filter((student) => isRowDirty(student.student_id))
  await saveStudents(dirtyStudents, 'displayed')
}

async function handleSaveAll() {
  const dirtyStudents = students.value.filter((student) => isRowDirty(student.student_id))
  await saveStudents(dirtyStudents, 'all')
}

function applyStatusToDisplayed(status) {
  if (displayedStudents.value.length === 0) return

  const nextAttendance = { ...attendanceData.value }
  displayedStudents.value.forEach((student) => {
    const currentValue = nextAttendance[student.student_id] || getDefaultAttendanceValue()
    nextAttendance[student.student_id] = {
      ...currentValue,
      status,
    }
  })
  attendanceData.value = nextAttendance
}

function handleHistorySearch() {
  loadClassHistory(0)
}

function handleHistoryPrev() {
  if (historyOffset.value <= 0) return
  loadClassHistory(Math.max(0, historyOffset.value - historyLimit.value))
}

function handleHistoryNext() {
  if (!historyHasMore.value) return
  loadClassHistory(historyOffset.value + historyLimit.value)
}

watch(selectedClassId, async () => {
  showMobileTakeControls.value = false
  await fetchStudentsAndAttendance()
})

watch(viewMode, (mode) => {
  if (mode === 'history') {
    loadClassHistory(0)
  }
})

onMounted(async () => {
  await fetchClasses()
  if (selectedClassId.value) {
    await fetchStudentsAndAttendance()
  }
})
</script>

<template>
  <div class="teacher-attendance page-stack">
    <div class="view-switch">
      <button
        class="view-switch__btn"
        :class="{ 'view-switch__btn--active': viewMode === 'take' }"
        type="button"
        @click="viewMode = 'take'"
      >
        Điểm danh hôm nay
      </button>
      <button
        class="view-switch__btn"
        :class="{ 'view-switch__btn--active': viewMode === 'history' }"
        type="button"
        @click="viewMode = 'history'"
      >
        Lịch sử lớp
      </button>
    </div>

    <div class="card filters-card">
      <div class="filters-grid">
        <div class="form-group mb-0">
          <label class="form-label">Chọn lớp học</label>
          <select v-model="selectedClassId" class="form-input" :disabled="isLoadingClasses">
            <option v-if="classes.length === 0" value="" disabled>-- Không có lớp --</option>
            <option v-for="cls in classes" :key="cls.class_id" :value="cls.class_id">
              {{ cls.name }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error alert-row">
      <AlertCircle :size="16" />
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="alert alert--success alert-row">
      <Check :size="16" />
      {{ successMessage }}
    </div>

    <LoadingSpinner v-if="isLoadingClasses || isLoadingStudents" message="Đang tải dữ liệu..." />

    <div v-else-if="classes.length === 0" class="card">
      <EmptyState
        title="Chưa có lớp học"
        message="Bạn chưa được phân công phụ trách lớp học nào."
      />
    </div>

    <div v-else-if="students.length === 0" class="card">
      <EmptyState
        title="Không có học sinh"
        message="Lớp này hiện chưa có học sinh nào."
        icon="users"
      />
    </div>

    <template v-else-if="viewMode === 'take'">
      <div class="card take-controls">
        <div class="take-controls__mobile-head">
          <p class="take-controls__summary">
            Hiển thị {{ displayedStudents.length }}/{{ students.length }} • Chờ lưu {{ displayedDirtyCount }}
          </p>
          <button
            class="btn btn--outline btn--sm take-controls__toggle"
            type="button"
            :aria-expanded="showMobileTakeControls"
            @click="showMobileTakeControls = !showMobileTakeControls"
          >
            {{ showMobileTakeControls ? 'Ẩn bộ lọc' : 'Mở bộ lọc' }}
          </button>
        </div>

        <div
          class="take-controls__panel"
          :class="{ 'take-controls__panel--open': showMobileTakeControls }"
        >
          <div class="take-controls__grid">
            <input
              v-model="studentSearch"
              class="form-input"
              type="text"
              placeholder="Tìm học sinh theo tên..."
            />

            <select v-model="takeListFilter" class="form-input">
              <option value="all">Tất cả học sinh</option>
              <option value="pending">Chưa lưu / đang sửa</option>
              <option value="saved">Đã lưu</option>
            </select>

            <button
              class="btn btn--sm"
              :class="listOrderMode === 'prioritize' ? 'btn--primary' : 'btn--outline'"
              type="button"
              @click="listOrderMode = 'prioritize'"
            >
              Ưu tiên chưa lưu
            </button>

            <button
              class="btn btn--sm"
              :class="listOrderMode === 'original' ? 'btn--primary' : 'btn--outline'"
              type="button"
              @click="listOrderMode = 'original'"
            >
              Giữ nguyên thứ tự
            </button>
          </div>

          <div class="take-badges">
            <span class="badge badge--outline">Toàn lớp chờ lưu: {{ globalPendingCount }}</span>
            <span class="badge badge--outline">Đang hiển thị: {{ displayedStudents.length }}/{{ students.length }}</span>
            <span class="badge" :class="displayedDirtyCount > 0 ? 'badge--warning' : 'badge--outline'">
              Chờ lưu trong danh sách: {{ displayedDirtyCount }}
            </span>
            <span class="badge" :class="displayedSavedCount > 0 ? 'badge--primary' : 'badge--outline'">
              Đã lưu trong danh sách: {{ displayedSavedCount }}
            </span>
          </div>

          <div class="bulk-actions">
            <button
              v-for="option in ATTENDANCE_STATUS_OPTIONS"
              :key="option.value"
              class="btn btn--outline btn--sm"
              type="button"
              @click="applyStatusToDisplayed(option.value)"
            >
              Đặt tất cả hiển thị: {{ option.label }}
            </button>

            <button
              class="btn btn--primary btn--sm"
              type="button"
              :disabled="isSavingDisplayed || displayedDirtyCount === 0"
              @click="handleSaveDisplayed"
            >
              {{ isSavingDisplayed ? 'Đang lưu...' : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ''}` }}
            </button>
          </div>

          <p class="take-controls__hint">
            {{ listOrderMode === 'prioritize' ? 'Đang ưu tiên học sinh chưa lưu.' : 'Đang giữ nguyên thứ tự danh sách.' }}
          </p>
        </div>
      </div>

      <div v-if="displayedStudents.length === 0" class="card">
        <EmptyState
          title="Không có học sinh phù hợp"
          message="Hãy đổi từ khóa tìm kiếm hoặc bộ lọc danh sách."
        />
      </div>

      <div v-else class="attendance-list">
        <article
          v-for="student in displayedStudents"
          :key="student.student_id"
          class="card attendance-item"
          :class="{ 'attendance-item--saved': !isRowDirty(student.student_id) && hasSavedForDate[student.student_id] }"
        >
          <div class="attendance-item__head">
            <div class="attendance-item__identity">
              <p class="student-name">
                {{ student.full_name }}
                <span class="student-meta-inline">• {{ formatDateVN(student.dob) }}</span>
              </p>
              <p class="student-meta">
                {{ !hasSavedForDate[student.student_id] ? 'Chưa lưu' : isRowDirty(student.student_id) ? 'Đã chỉnh sửa, chưa lưu' : 'Đã lưu' }}
              </p>
            </div>

            <div class="attendance-item__mobile-status">
              <select
                class="form-input"
                :value="attendanceData[student.student_id]?.status || 'present'"
                @change="handleAttendanceStatusChange(student.student_id, $event.target.value)"
              >
                <option v-for="option in ATTENDANCE_STATUS_OPTIONS" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>

            <div class="attendance-item__status-chips">
              <button
                v-for="option in ATTENDANCE_STATUS_OPTIONS"
                :key="option.value"
                class="status-chip"
                :class="{ 'status-chip--active': attendanceData[student.student_id]?.status === option.value }"
                type="button"
                @click="handleAttendanceStatusChange(student.student_id, option.value)"
              >
                {{ option.label }}
              </button>
            </div>
          </div>

          <div class="attendance-item__actions-row">
            <input
              class="form-input attendance-note"
              type="text"
              placeholder="Ghi chú..."
              :value="attendanceData[student.student_id]?.note || ''"
              @input="handleAttendanceNoteChange(student.student_id, $event.target.value)"
            />

            <button
              class="btn btn--sm"
              :class="!isRowDirty(student.student_id) && hasSavedForDate[student.student_id] ? 'btn--outline' : 'btn--primary'"
              type="button"
              :disabled="savingRowId === student.student_id"
              @click="handleMark(student.student_id)"
            >
              {{
                savingRowId === student.student_id
                  ? 'Đang lưu...'
                  : !isRowDirty(student.student_id) && hasSavedForDate[student.student_id]
                    ? 'Đã lưu'
                    : hasSavedForDate[student.student_id]
                      ? 'Cập nhật'
                      : 'Lưu'
              }}
            </button>

            <button
              v-if="hasSavedForDate[student.student_id] && isRowDirty(student.student_id)"
              class="btn btn--outline btn--sm"
              type="button"
              @click="handleRevert(student.student_id)"
            >
              Hoàn tác
            </button>

            <button
              v-if="hasSavedForDate[student.student_id] && !isRowDirty(student.student_id)"
              class="btn btn--outline btn--sm attendance-item__cancel"
              type="button"
              :disabled="cancelingRowId === student.student_id"
              @click="handleCancelSaved(student.student_id)"
            >
              {{ cancelingRowId === student.student_id ? 'Đang hủy...' : 'Hủy lưu hôm nay' }}
            </button>
          </div>

          <div class="attendance-item__history-wrap">
            <button
              class="attendance-item__history-toggle"
              type="button"
              @click="toggleHistory(student.student_id)"
            >
              <History :size="14" />
              {{ historyOpen.has(student.student_id) ? 'Ẩn lịch sử' : 'Xem lịch sử 30 ngày' }}
            </button>

            <div v-if="historyOpen.has(student.student_id)" class="attendance-history">
              <div v-if="historyLoading.has(student.student_id)" class="loading-inline">
                <LoaderCircle class="spin text-muted" :size="16" />
                Đang tải lịch sử...
              </div>

              <p v-else-if="(historyByStudent[student.student_id] || []).length === 0" class="attendance-history__empty">
                Chưa có lịch sử điểm danh.
              </p>

              <div v-else class="attendance-history__list">
                <div
                  v-for="record in historyByStudent[student.student_id].slice(0, 8)"
                  :key="record.change_id"
                  class="attendance-history__item"
                >
                  <div class="attendance-history__row">
                    <span class="text-muted">{{ formatDateTimeVN(record.changed_at) }}</span>
                    <span class="attendance-history__type">
                      {{ record.change_type === 'create' ? 'Tạo mới' : record.change_type === 'delete' ? 'Hủy lưu' : 'Cập nhật' }}
                    </span>
                  </div>
                  <div class="attendance-history__row text-muted">
                    <span>
                      {{
                        record.change_type === 'create'
                          ? `Tạo: ${getStatusLabel(record.new_status)}`
                          : record.change_type === 'delete'
                            ? `${getStatusLabel(record.old_status)} → Đã hủy`
                            : `${getStatusLabel(record.old_status)} → ${getStatusLabel(record.new_status)}`
                      }}
                    </span>
                    <span class="attendance-history__note">
                      {{
                        record.change_type === 'delete'
                          ? `${record.old_note || '-'} → Đã xóa`
                          : `${record.old_note || '-'} → ${record.new_note || '-'}`
                      }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </article>
      </div>

      <div
        v-if="displayedDirtyCount > 0 || globalPendingCount > 0"
        class="take-summary-bar"
      >
        <p class="take-summary-bar__text">
          Còn {{ displayedDirtyCount }} học sinh chưa lưu trong danh sách hiển thị • Toàn lớp còn {{ globalPendingCount }} học sinh chưa lưu.
        </p>

        <div class="take-summary-bar__actions">
          <button
            class="btn btn--outline btn--sm"
            type="button"
            :disabled="isSavingDisplayed || displayedDirtyCount === 0"
            @click="handleSaveDisplayed"
          >
            {{ isSavingDisplayed ? 'Đang lưu...' : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ''}` }}
          </button>
          <button
            class="btn btn--primary btn--sm"
            type="button"
            :disabled="isSavingAll || dirtyCount === 0"
            @click="handleSaveAll"
          >
            {{ isSavingAll ? 'Đang lưu...' : `Lưu toàn lớp${dirtyCount > 0 ? ` (${dirtyCount})` : ''}` }}
          </button>
        </div>
      </div>
    </template>

    <div v-else class="history-stack">
      <div class="card history-filters-card">
        <div class="history-filters-grid">
          <input v-model="historyFrom" type="date" class="form-input" />
          <input v-model="historyTo" type="date" class="form-input" />

          <select v-model="historyStudentId" class="form-input">
            <option value="all">Tất cả học sinh</option>
            <option v-for="student in students" :key="student.student_id" :value="student.student_id">
              {{ student.full_name }}
            </option>
          </select>

          <select v-model="historyStatus" class="form-input">
            <option value="all">Tất cả trạng thái</option>
            <option v-for="option in ATTENDANCE_STATUS_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>

        <div class="history-actions-row">
          <button class="btn btn--primary btn--sm" type="button" :disabled="historyListLoading" @click="handleHistorySearch">
            {{ historyListLoading ? 'Đang tải...' : 'Xem lịch sử' }}
          </button>
          <p class="history-total">Tổng bản ghi: {{ historyTotal }}</p>
        </div>
      </div>

      <div class="card history-list-card">
        <div v-if="historyListLoading" class="loading-inline">
          <LoaderCircle class="spin text-muted" :size="20" />
          Đang tải lịch sử...
        </div>

        <EmptyState
          v-else-if="historyList.length === 0"
          title="Không có dữ liệu lịch sử"
          message="Không có bản ghi phù hợp với bộ lọc hiện tại."
          icon="users"
        />

        <div v-else class="history-list">
          <article v-for="record in historyList" :key="record.change_id" class="history-item">
            <div class="history-item__head">
              <div>
                <p class="student-name">{{ record.student_name }}</p>
                <p class="student-meta">{{ formatDateTimeVN(record.changed_at) }}</p>
              </div>
              <span class="badge badge--outline">
                {{ record.change_type === 'create' ? 'Tạo mới' : record.change_type === 'delete' ? 'Hủy lưu' : 'Cập nhật' }}
              </span>
            </div>

            <div class="history-item__body">
              <p>
                {{
                  record.change_type === 'create'
                    ? `Tạo: ${getStatusLabel(record.new_status)}`
                    : record.change_type === 'delete'
                      ? `${getStatusLabel(record.old_status)} → Đã hủy`
                      : `${getStatusLabel(record.old_status)} → ${getStatusLabel(record.new_status)}`
                }}
              </p>
              <p class="text-muted">
                {{
                  record.change_type === 'delete'
                    ? `${record.old_note || '-'} → Đã xóa`
                    : `${record.old_note || '-'} → ${record.new_note || '-'}`
                }}
              </p>
            </div>
          </article>
        </div>
      </div>

      <div class="history-pagination">
        <button class="btn btn--outline btn--sm" type="button" :disabled="historyOffset === 0" @click="handleHistoryPrev">
          Trang trước
        </button>
        <p class="history-total">
          {{ historyTotal === 0 ? '0-0' : `${historyOffset + 1}-${Math.min(historyOffset + historyLimit, historyTotal)}` }} / {{ historyTotal }}
        </p>
        <button class="btn btn--outline btn--sm" type="button" :disabled="!historyHasMore" @click="handleHistoryNext">
          Trang sau
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-stack,
.attendance-list,
.history-stack,
.history-list,
.attendance-history__list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.view-switch {
  display: inline-flex;
  gap: var(--spacing-1);
  width: fit-content;
  padding: 0.25rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}

.view-switch__btn {
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  border-radius: var(--radius-sm);
  padding: 0.5rem 0.9rem;
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.view-switch__btn--active {
  background: var(--color-background);
  color: var(--color-text);
  box-shadow: var(--shadow-sm);
}

.filters-card,
.take-controls,
.history-filters-card,
.history-list-card,
.attendance-item {
  padding: var(--spacing-4);
}

.filters-grid,
.take-controls__grid,
.history-filters-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: 1fr;
}

.take-controls__mobile-head,
.take-badges,
.bulk-actions,
.history-actions-row,
.alert-row,
.loading-inline,
.history-pagination,
.take-summary-bar__actions,
.attendance-item__actions-row,
.attendance-history__row {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  flex-wrap: wrap;
}

.take-controls__mobile-head,
.history-actions-row,
.history-pagination {
  justify-content: space-between;
}

.take-controls__mobile-head {
  display: none;
}

.take-controls__summary,
.take-controls__hint,
.student-name,
.student-meta,
.history-total,
.history-item__body p,
.attendance-history__empty {
  margin: 0;
}

.take-controls__panel {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.take-controls__hint,
.student-meta,
.student-meta-inline,
.history-total,
.text-muted,
.attendance-history__empty,
.attendance-history__note {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.student-name {
  font-weight: 700;
  color: var(--color-text);
}

.student-meta-inline {
  font-weight: 400;
}

.attendance-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.attendance-item--saved {
  border-color: color-mix(in srgb, var(--color-success) 30%, var(--color-border));
  background: color-mix(in srgb, var(--color-success) 7%, transparent);
}

.attendance-item__head,
.history-item__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-3);
}

.attendance-item__identity {
  min-width: 0;
  flex: 1;
}

.attendance-item__mobile-status {
  width: 100%;
}

.attendance-item__status-chips {
  display: none;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.status-chip {
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  border-radius: var(--radius-full);
  padding: 0.25rem 0.65rem;
  font-size: var(--font-size-xs);
  font-weight: 700;
}

.status-chip--active {
  border-color: color-mix(in srgb, var(--color-primary) 30%, var(--color-border));
  background: color-mix(in srgb, var(--color-primary) 12%, transparent);
  color: var(--color-primary);
}

.attendance-item__actions-row {
  align-items: stretch;
}

.attendance-note {
  flex: 1;
  min-width: 15rem;
}

.attendance-item__cancel {
  color: var(--color-danger);
}

.attendance-item__history-wrap {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.attendance-item__history-toggle {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 0;
}

.attendance-history {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--color-background) 70%, transparent);
  padding: var(--spacing-3);
}

.attendance-history__item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  font-size: var(--font-size-xs);
}

.attendance-history__row {
  justify-content: space-between;
  align-items: flex-start;
}

.attendance-history__type {
  font-weight: 700;
  color: var(--color-text);
}

.history-item {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--spacing-3);
}

.history-item__body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  font-size: var(--font-size-sm);
}

.take-summary-bar {
  position: sticky;
  bottom: 0.75rem;
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  padding: var(--spacing-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--color-surface) 95%, transparent);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(8px);
}

.take-summary-bar__text {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 768px) {
  .filters-grid {
    grid-template-columns: minmax(0, 22rem);
  }

  .take-controls__grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .history-filters-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .attendance-item__mobile-status {
    display: none;
  }

  .attendance-item__status-chips {
    display: flex;
    justify-content: flex-end;
    max-width: 50%;
  }

  .take-summary-bar {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

@media (max-width: 767px) {
  .view-switch {
    width: 100%;
  }

  .view-switch__btn {
    flex: 1;
  }

  .take-controls__mobile-head {
    display: flex;
  }

  .take-controls__panel {
    display: none;
  }

  .take-controls__panel--open {
    display: flex;
  }

  .attendance-item__head {
    flex-direction: column;
  }

  .attendance-note {
    min-width: 100%;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
