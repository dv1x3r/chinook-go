import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createPlaylistLayout() {
  const trackGrid = createPlaylistTrackGrid()
  const playlistGrid = createPlaylistGrid(trackGrid)
  return new w2layout({
    name: 'playlistLayout',
    padding: 5,
    panels: [
      { type: 'left', size: '40%', html: playlistGrid, resizable: true },
      { type: 'main', size: '60%', html: trackGrid },
    ],
    onDestroy: function() {
      playlistGrid.destroy()
      trackGrid.destroy()
    },
  })
}

function createPlaylistGrid(trackGrid) {
  return new w2grid({
    name: 'playlistGrid',
    url: {
      get: '/api/v1/playlist/grid',
      remove: '/api/v1/playlist/delete',
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
        text: 'Playlist',
        size: '100%',
        min: '160px',
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
    onSelect: function(event) { helpers.reloadChildGridFromSelection(trackGrid, event) },
  })
}

function createPlaylistTrackGrid() {
  const grid = new w2grid({
    name: 'playlistTrackGrid',
    url: {
      get: '/api/v1/playlist/:playlistID/track/grid',
      remove: '/api/v1/playlist/:playlistID/track/delete',
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
        size: '60px',
        hidden: true,
        searchable: 'int',
        sortable: true,
      },
      {
        field: 'name',
        text: 'Track',
        size: '100%',
        min: '280px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'artist',
        text: 'Artist',
        size: '160px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'album',
        text: 'Album',
        size: '160px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'genre',
        text: 'Genre',
        size: '160px',
        render: 'dropdown',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'name', direction: 'asc' },
    ],
    onAdd: function() { openAddTrackPopup(grid) },
  })
  helpers.bindChildGrid(grid, { route: 'playlistID' })
  return grid
}

function openAddTrackPopup(grid) {
  const form = new w2form({
    name: 'playlistTrackForm',
    url: `/api/v1/playlist/${helpers.parentIDOf(grid)}/track/add`,
    focus: -1,
    fields: [
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
    title: 'Add track to the playlist',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 180, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

function openFormPopup(event) {
  const grid = event.owner
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'playlistForm',
    url: '/api/v1/playlist/form',
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
    title: isEditMode ? 'Edit Playlist' : 'New Playlist',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 200, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

