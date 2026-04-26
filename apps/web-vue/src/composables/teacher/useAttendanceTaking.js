import { ref, computed } from 'vue'
import { teacherService } from '../../services/teacherService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { getDefaultAttendanceValue, isSameAttendance } from '../../helpers/attendanceConfig'
import { getTodayDateString } from '../../helpers/dateHelpers'

export function useAttendanceTaking() {
  const students = ref([])
  const attendanceData = ref({})
  const savedAttendance = ref({})
  const hasSavedForDate = ref({})

  const isLoadingStudents = ref(false)
  const savingRowId = ref('')
  const cancelingRowId = ref('')
  const isSavingAll = ref(false)
  const isSavingDisplayed = ref(false)
  const errorMessage = ref('')
  const successMessage = ref('')

  const studentSearch = ref('')
  const takeListFilter = ref('all')
  const listOrderMode = ref('prioritize')
  const showMobileTakeControls = ref(false)

  const today = getTodayDateString()

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

  async function fetchStudentsAndAttendance(classId) {
    if (!classId) {
      students.value = []
      attendanceData.value = {}
      savedAttendance.value = {}
      hasSavedForDate.value = {}
      return
    }

    isLoadingStudents.value = true
    errorMessage.value = ''
    successMessage.value = ''

    try {
      const studentList = normalizeListResponse(await teacherService.getStudentsInClass(classId))
      students.value = studentList

      const initialAttendance = {}
      const initialSavedAttendance = {}
      const initialHasSaved = {}

      await Promise.all(
        studentList.map(async (student) => {
          const fallback = getDefaultAttendanceValue()

          try {
            const records = normalizeListResponse(
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

  return {
    students,
    attendanceData,
    savedAttendance,
    hasSavedForDate,
    isLoadingStudents,
    savingRowId,
    cancelingRowId,
    isSavingAll,
    isSavingDisplayed,
    errorMessage,
    successMessage,
    studentSearch,
    takeListFilter,
    listOrderMode,
    showMobileTakeControls,
    displayedStudents,
    dirtyCount,
    displayedDirtyCount,
    displayedSavedCount,
    globalPendingCount,
    isRowDirty,
    fetchStudentsAndAttendance,
    handleAttendanceStatusChange,
    handleAttendanceNoteChange,
    handleMark,
    handleRevert,
    handleCancelSaved,
    handleSaveDisplayed,
    handleSaveAll,
    applyStatusToDisplayed,
  }
}
