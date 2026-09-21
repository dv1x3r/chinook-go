import { w2form, w2grid, w2layout, w2popup, w2utils } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createInvoiceLayout() {
  const lineGrid = createInvoiceLineGrid(reloadInvoices)
  const invoiceForm = createInvoiceForm(lineGrid, reloadInvoices)
  const invoiceGrid = createInvoiceGrid(lineGrid, invoiceForm)

  // the total lives on the invoice row, so the grid is re-read after every change
  function reloadInvoices() {
    return invoiceGrid.reload()
  }

  return new w2layout({
    name: 'invoiceLayout',
    padding: 5,
    panels: [
      { type: 'left', size: '50%', html: invoiceGrid, resizable: true },
      { type: 'main', size: '50%', html: invoiceForm },
    ],
    onDestroy: function() {
      invoiceGrid.destroy()
      invoiceForm.destroy()
    },
  })
}

function createInvoiceGrid(lineGrid, invoiceForm) {
  return new w2grid({
    name: 'invoiceGrid',
    url: {
      get: '/api/v1/invoice/grid',
      remove: '/api/v1/invoice/delete',
    },
    recid: 'id',
    recordHeight: 28,
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
        size: '70px',
        sortable: true,
      },
      {
        field: 'invoice_date',
        text: 'Date',
        size: '150px',
        render: 'datetime',
        sortable: true,
      },
      {
        field: 'customer',
        text: 'Customer',
        size: '160px',
        render: 'dropdown',
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/customer/dropdown') },
        sortable: true,
      },
      {
        field: 'billing_city',
        text: 'City',
        size: '130px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'billing_country',
        text: 'Country',
        size: '130px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'total',
        text: 'Total',
        size: '90px',
        render: 'money',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function(event) { openInvoicePopup(event.owner) },
    onDelete: function(event) {
      if (event.detail.data?.status == 'success') {
        invoiceForm.recid = null
        invoiceForm.clear()
        helpers.reloadChildGrid(lineGrid, 0)
      }
    },
    onSelect: async function(event) {
      await helpers.setFormRecordFromGridSelection(event, invoiceForm)
      await helpers.reloadChildGridFromSelection(lineGrid, event)
    },
  })
}

function createInvoiceLineGrid(onLineChange) {
  const grid = new w2grid({
    name: 'invoiceLineGrid',
    url: {
      get: '/api/v1/invoice/:invoiceID/item/grid',
      remove: '/api/v1/invoice/:invoiceID/item/delete',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: false,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: false,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        hidden: true,
      },
      {
        field: 'track',
        text: 'Track',
        size: '100%',
        min: '180px',
        render: 'dropdown',
      },
      {
        field: 'quantity',
        text: 'Qty',
        size: '60px',
        render: 'int',
      },
      {
        field: 'unit_price',
        text: 'Price',
        size: '80px',
        render: 'money',
      },
      {
        field: 'total',
        text: 'Total',
        size: '80px',
        render: 'money',
      },
    ],
    onAdd: function(event) { openLinePopup(event, onLineChange) },
    onEdit: function(event) { openLinePopup(event, onLineChange) },
    onDblClick: function(event) { openLinePopup(event, onLineChange) },
    onDelete: function(event) {
      if (event.detail.data?.status == 'success') {
        onLineChange()
      }
    },
  })
  helpers.bindChildGrid(grid, { route: 'invoiceID' })
  return grid
}

function createInvoiceForm(lineGrid, onSaved) {
  return new w2form({
    name: 'invoiceForm',
    url: '/api/v1/invoice/form',
    focus: -1,
    style: 'min-height: 100%; max-height: 100%;',
    fields: [
      {
        type: 'html',
        html: {
          html: '<div id="invoice-line-grid" style="height: 260px;"></div>',
          span: -1,
          column: 'before',
          group: 'Invoice Lines',
        },
      },
      {
        field: 'customer',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/customer/dropdown'),
        html: {
          label: 'Customer',
          attr: 'placeholder="Type to search..."',
          span: 4,
          group: 'Invoice',
        },
      },
      {
        field: 'invoice_date',
        type: 'text',
        required: true,
        html: {
          label: 'Date',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_address',
        type: 'text',
        html: {
          label: 'Address',
          attr: 'style="width: 100%;"',
          span: 4,
          group: 'Billing Address',
        },
      },
      {
        field: 'billing_city',
        type: 'text',
        html: {
          label: 'City',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_state',
        type: 'text',
        html: {
          label: 'State',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_country',
        type: 'text',
        html: {
          label: 'Country',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_postal_code',
        type: 'text',
        html: {
          label: 'Postal code',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
    ],
    actions: {
      async Save() {
        if (this.record.id) {
          const res = await this.save()
          if (res.status == 'success') {
            w2utils.notify('Form has been successfully saved!', { timeout: 4000 })
            await onSaved()
          }
        }
      },
    },
    onRender: async function(event) {
      await event.complete
      lineGrid.render('#invoice-line-grid')
    },
    onDestroy: function() {
      lineGrid.destroy()
    },
  })
}

function openLinePopup(event, onLineChange) {
  const grid = event.owner
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'invoiceLineForm',
    url: `/api/v1/invoice/${helpers.parentIDOf(grid)}/item/form`,
    focus: -1,
    record: record ?? { quantity: 1 },
    fields: [
      {
        field: 'id',
        type: 'text',
        html: {
          label: 'ID',
          attr: 'size="10" readonly',
          span: 4,
        },
      },
      {
        field: 'track',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/track/dropdown'),
        html: {
          label: 'Track',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'quantity',
        type: 'int',
        required: true,
        html: {
          label: 'Quantity',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'unit_price',
        type: 'money',
        required: true,
        html: {
          label: 'Price',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
    ],
    actions: {
      async Save() {
        const res = await this.save()
        if (res.status == 'success') {
          grid.reload()
          onLineChange()
          w2popup.close()
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: isEditMode ? 'Edit Invoice line' : 'Add track to the invoice',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 280, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

function openInvoicePopup(grid) {
  const form = new w2form({
    name: 'newInvoiceForm',
    url: '/api/v1/invoice/form',
    focus: -1,
    fields: [
      {
        field: 'customer',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/customer/dropdown'),
        html: {
          label: 'Customer',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'invoice_date',
        type: 'text',
        required: true,
        html: {
          label: 'Date',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_address',
        type: 'text',
        html: {
          label: 'Address',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_city',
        type: 'text',
        html: {
          label: 'City',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_state',
        type: 'text',
        html: {
          label: 'State',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_country',
        type: 'text',
        html: {
          label: 'Country',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'billing_postal_code',
        type: 'text',
        html: {
          label: 'Postal code',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
    ],
    record: { invoice_date: w2utils.formatDateTime(new Date(), 'yyyy-MM-dd|hh24:mi:ss') },
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
    title: 'New Invoice',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 400, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

