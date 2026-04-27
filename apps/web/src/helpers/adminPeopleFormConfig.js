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

export function createAdminPersonRelationConfig({
  ownerIdField,
  relationIdField,
  relationNameField,
  unassignNameField,
  assignSelectionField,
  assignService,
  unassignService,
}) {
  const assignItem = ({ target, ...selection }) =>
    assignService(target[ownerIdField], selection[assignSelectionField])

  const toUnassignTarget = (owner, relation) => ({
    [ownerIdField]: owner[ownerIdField],
    [relationIdField]: relation[relationIdField],
    [unassignNameField]: relation[relationNameField],
  })

  const unassignItem = (target) => unassignService(target[ownerIdField], target[relationIdField])

  return {
    assignItem,
    toUnassignTarget,
    unassignItem,
  }
}
