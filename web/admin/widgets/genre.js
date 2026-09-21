import { w2form, w2grid, w2popup } from '/lib/w2ui.es6.min.js'

export function createGenreGrid() {
  return new w2grid({
    name: 'genreGrid',
    url: {
      get: '/api/v1/genre/grid',
      remove: '/api/v1/genre/delete',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
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
        field: 'name',
        text: 'Genre',
        size: '280px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'tracks',
        text: '# Tracks',
        size: '90px',
        render: 'int',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'name', direction: 'asc' },
    ],
    onAdd: function(event) { openFormPopup(event) },
    onEdit: function(event) { openFormPopup(event) },
    onDblClick: function(event) { openFormPopup(event) },
  })
}

function openFormPopup(event) {
  const grid = event.owner
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'genreForm',
    url: '/api/v1/genre/form',
    focus: -1,
    record: record,
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
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Name',
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
          w2popup.close()
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: isEditMode ? 'Edit Genre' : 'New Genre',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 200, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

