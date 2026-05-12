interface EditFormConfig<T = any> {
  createInitialEditForm: () => T
  toEditForm: (person: any, context: any) => T
  validateEditForm: (form: T) => string
}

export function createAdminPersonEditFormConfig({ idField }: { idField: string }): EditFormConfig {
  const createInitialEditForm = () => ({
    [idField]: '',
    full_name: '',
    phone: '',
    school_id: '',
  })

  const toEditForm = (person: any, context: any) => ({
    [idField]: person[idField],
    full_name: person.full_name || '',
    phone: person.phone || '',
    school_id: person.school_id || context.selectedSchoolId || '',
  })

  const validateEditForm = (form: any) => {
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

interface RelationConfig {
  ownerIdField: string
  relationIdField: string
  relationNameField: string
  unassignNameField: string
  assignSelectionField: string
  assignService: (ownerId: string, relationId: string) => Promise<any>
  unassignService: (ownerId: string, relationId: string) => Promise<any>
}

export function createAdminPersonRelationConfig(config: RelationConfig) {
  const {
    ownerIdField,
    relationIdField,
    relationNameField,
    unassignNameField,
    assignSelectionField,
    assignService,
    unassignService,
  } = config

  const assignItem = ({ target, ...selection }: any) =>
    assignService(target[ownerIdField], selection[assignSelectionField])

  const toUnassignTarget = (owner: any, relation: any) => ({
    [ownerIdField]: owner[ownerIdField],
    [relationIdField]: relation[relationIdField],
    [unassignNameField]: relation[relationNameField],
  })

  const unassignItem = (target: any) => unassignService(target[ownerIdField], target[relationIdField])

  return {
    assignItem,
    toUnassignTarget,
    unassignItem,
  }
}
