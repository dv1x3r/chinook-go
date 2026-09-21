import { w2form, w2grid, w2layout, w2popup, w2utils } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createEmployeeLayout() {
  const grid = new w2grid({
    name: 'employeeGrid',
    url: {
      get: '/api/v1/employee/grid',
      remove: '/api/v1/employee/delete',
    },
    recid: 'id',
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: true,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        hidden: true,
        sortable: true,
      },
      {
        field: 'full_name',
        text: 'Full Name',
        size: '180px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'title',
        text: 'Title',
        size: '160px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'birth_date',
        text: 'Birth Date',
        size: '110px',
        render: 'date',
        sortable: true,
      },
      {
        field: 'hire_date',
        text: 'Hire Date',
        size: '110px',
        render: 'date',
        sortable: true,
      },
      {
        field: 'reports_to',
        text: 'Reports To',
        size: '160px',
        render: 'dropdown',
        sortable: true,
      },
      {
        field: 'customers',
        text: '# Customers',
        size: '100px',
        render: 'int',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'full_name', direction: 'asc' },
    ],
    onAdd: function() { openFormPopup(grid) },
    onDelete: function(event) {
      if (event.detail.data?.status == 'success') {
        form.recid = null
        form.clear()
      }
    },
    onSelect: function(event) { helpers.setFormRecordFromGridSelection(event, form) },
  })

  const form = new w2form({
    name: 'employeeForm',
    url: '/api/v1/employee/form',
    focus: -1,
    style: 'min-height: 100%; max-height: 100%;',
    fields: employeeFields(),
    actions: {
      async Save() {
        if (!this.recid) {
          return
        }
        const res = await this.save()
        if (res.status == 'success') {
          w2utils.notify('Form has been successfully saved!', { timeout: 4000 })
          await grid.reload()
        }
      },
    },
  })

  return new w2layout({
    name: 'employeeLayout',
    padding: 5,
    panels: [
      { type: 'left', size: '70%', html: grid, resizable: true },
      { type: 'main', size: '30%', html: form },
    ],
    onDestroy: function() {
      grid.destroy()
      form.destroy()
    },
  })
}

function openFormPopup(grid) {
  const form = new w2form({
    name: 'newEmployeeForm',
    url: '/api/v1/employee/form',
    focus: -1,
    fields: employeeFields(),
    actions: {
      async Save() {
        const res = await this.save()
        if (res.status == 'success') {
          grid.reload()
          w2popup.close()
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: 'New Employee',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 400, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

function employeeFields() {
  return [
    {
      field: 'id',
      type: 'text',
      html: {
        label: 'ID',
        attr: 'size="10" readonly',
        span: 4,
        group: 'Employee',
      },
    },
    {
      field: 'first_name',
      type: 'text',
      required: true,
      html: {
        label: 'First name',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'last_name',
      type: 'text',
      required: true,
      html: {
        label: 'Last name',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'title',
      type: 'text',
      html: {
        label: 'Title',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'reports_to',
      type: 'list',
      options: helpers.remoteListOptions('/api/v1/employee/dropdown'),
      html: {
        label: 'Reports to',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'birth_date',
      type: 'text',
      html: {
        label: 'Birth date',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'hire_date',
      type: 'text',
      html: {
        label: 'Hire date',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'email',
      type: 'text',
      html: {
        label: 'Email',
        attr: 'style="width: 100%;"',
        span: 4,
        group: 'Contact',
      },
    },
    {
      field: 'phone',
      type: 'text',
      html: {
        label: 'Phone',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'fax',
      type: 'text',
      html: {
        label: 'Fax',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'address',
      type: 'text',
      html: {
        label: 'Address',
        attr: 'style="width: 100%;"',
        span: 4,
        group: 'Address',
      },
    },
    {
      field: 'city',
      type: 'text',
      html: {
        label: 'City',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'state',
      type: 'text',
      html: {
        label: 'State',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'country',
      type: 'text',
      html: {
        label: 'Country',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
    {
      field: 'postal_code',
      type: 'text',
      html: {
        label: 'Postal code',
        attr: 'style="width: 100%;"',
        span: 4,
      },
    },
  ]
}

