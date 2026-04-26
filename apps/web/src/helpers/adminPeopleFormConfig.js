export function createAdminPersonEditFormConfig({ idField }) {
  const createInitialEditForm = () => ({
    [idField]: '',
    full_name: '',
    phone: '',
    school_id: '',
  })

  const toEditForm = (person, context) => ({
    [idField]: person[idField],
    full_name: person.full_name || '',
    phone: person.phone || '',
    school_id: person.school_id || context.selectedSchoolId || '',
  })

  const validateEditForm = (form) => {
    if (!form[idField] || !form.full_name.trim() || !form.school_id) {
      return 'Vui lòng nhập đầy đủ thông tin bắt buộc'
    }

    return ''
  }

  return {
    createInitialEditForm,
    toEditForm,
    validateEditForm,
  }
}
