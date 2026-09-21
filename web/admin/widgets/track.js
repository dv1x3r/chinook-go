import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createTrackLayout() {
  const trackGrid = createTrackGrid(albumID => albumGrid.get(albumID)?.title)
  const albumGrid = createAlbumGrid(trackGrid, artistID => artistGrid.get(artistID)?.name)
  const artistGrid = createArtistGrid(albumGrid, trackGrid)
  return new w2layout({
    name: 'trackLayout',
    padding: 5,
    panels: [
      { type: 'main', size: '40%', html: artistGrid },
      { type: 'preview', size: '50%', html: albumGrid, resizable: true },
      { type: 'right', size: '60%', html: trackGrid, resizable: true },
    ],
    onDestroy: function() {
      artistGrid.destroy()
      albumGrid.destroy()
      trackGrid.destroy()
    },
  })
}

function createArtistGrid(albumGrid, trackGrid) {
  return new w2grid({
    name: 'artistGrid',
    url: {
      get: '/api/v1/artist/grid',
      remove: '/api/v1/artist/delete',
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
        text: 'Artist',
        size: '100%',
        min: '160px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'albums',
        text: '# Albums',
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
    onAdd: function(event) { openArtistPopup(event) },
    onEdit: function(event) { openArtistPopup(event) },
    onDblClick: function(event) { openArtistPopup(event) },
    onSelect: async function(event) {
      await event.complete
      const selection = event.owner.getSelection()
      const artistID = selection.length == 1 ? selection[0] : 0
      albumGrid.routeData.artistID = artistID
      albumGrid.toolbar.set('w2ui-add', { disabled: !artistID })
      // clear drops the album picked before, so reload does not bring it back
      albumGrid.clear()
      albumGrid.reload()
      trackGrid.routeData = { artistID: artistID, albumID: 0 }
      trackGrid.toolbar.disable('w2ui-add')
      trackGrid.reload()
    },
  })
}

function createAlbumGrid(trackGrid, artistName) {
  const grid = new w2grid({
    name: 'albumGrid',
    url: {
      get: '/api/v1/artist/:artistID/album/grid',
      remove: '/api/v1/album/delete',
    },
    routeData: { artistID: 0 },
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
        field: 'title',
        text: 'Album',
        size: '100%',
        min: '160px',
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
      },
      {
        field: 'tracks',
        text: '# Tracks',
        size: '90px',
        render: 'int',
        sortable: true,
      },
    ],
    sortData: [
      { field: 'title', direction: 'asc' },
    ],
    onAdd: function(event) {
      const artistID = grid.routeData.artistID
      openAlbumPopup(event, { artist: { id: artistID, text: artistName(artistID) } })
    },
    onEdit: function(event) { openAlbumPopup(event) },
    onDblClick: function(event) { openAlbumPopup(event) },
    onSelect: async function(event) {
      await event.complete
      const selection = event.owner.getSelection()
      const albumID = selection.length == 1 ? selection[0] : 0
      trackGrid.routeData.albumID = albumID
      trackGrid.toolbar.set('w2ui-add', { disabled: !albumID })
      trackGrid.reload()
    },
  })
  grid.toolbar.disable('w2ui-add')
  return grid
}

function createTrackGrid(albumTitle) {
  const grid = new w2grid({
    name: 'trackGrid',
    url: {
      get: '/api/v1/artist/:artistID/album/:albumID/track/grid',
      save: '/api/v1/track/grid',
      remove: '/api/v1/track/delete',
    },
    routeData: { artistID: 0, albumID: 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: true,
      toolbarSearch: true,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        hidden: false,
        searchable: 'int',
        sortable: true,
      },
      {
        field: 'name',
        text: 'Track',
        size: '280px',
        render: 'text',
        searchable: 'text',
        searchAll: true,
        sortable: true,
      },
      {
        field: 'composer',
        text: 'Composer',
        size: '160px',
        render: 'text',
        editable: { type: 'text' },
        sortable: true,
        searchable: 'text',
        searchAll: true,
      },
      {
        field: 'album',
        text: 'Album',
        size: '160px',
        render: 'dropdown',
        editable: helpers.remoteListOptions('/api/v1/album/dropdown'),
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/album/dropdown') },
        sortable: true,
      },
      {
        field: 'genre',
        text: 'Genre',
        size: '160px',
        render: 'dropdown',
        editable: helpers.remoteListOptions('/api/v1/genre/dropdown'),
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/genre/dropdown') },
        sortable: true,
      },
      {
        field: 'media_type',
        text: 'Media',
        size: '160px',
        render: 'dropdown',
        editable: helpers.remoteListOptions('/api/v1/media_type/dropdown'),
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/media_type/dropdown') },
        sortable: true,
      },
      {
        field: 'milliseconds',
        text: 'Length',
        size: '80px',
        sortable: true,
        render: row => {
          const milliseconds = row.milliseconds
          if (milliseconds == null) {
            return ''
          }
          const total = Math.round(milliseconds / 1000)
          return `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`
        },
      },
      {
        field: 'unit_price',
        text: 'Price',
        size: '80px',
        render: 'money',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'asc' },
    ],
    onAdd: function(event) {
      const albumID = grid.routeData.albumID
      openTrackPopup(event, { album: { id: albumID, text: albumTitle(albumID) } })
    },
    onEdit: function(event) { openTrackPopup(event) },
    onDblClick: function(event) { helpers.doubleClickNonEditable(event, openTrackPopup) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onSearch: function(event) { helpers.searchAllFilter(event) },
  })
  grid.toolbar.disable('w2ui-add')
  return grid
}

function openArtistPopup(event) {
  const grid = event.owner
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'artistForm',
    url: '/api/v1/artist/form',
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
    title: isEditMode ? 'Edit Artist' : 'New Artist',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 200, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

function openAlbumPopup(event, defaults = null) {
  const grid = event.owner
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'albumForm',
    url: '/api/v1/album/form',
    focus: -1,
    record: record ?? defaults,
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
        field: 'title',
        type: 'text',
        required: true,
        html: {
          label: 'Title',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'artist',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/artist/dropdown'),
        html: {
          label: 'Artist',
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
    title: isEditMode ? 'Edit Album' : 'New Album',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 240, showMax: false, resizable: true,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

function openTrackPopup(event, defaults = null) {
  const grid = event.owner
  if (grid.getChanges().length > 0) {
    grid.message('You have unsaved changes. Save or discard them first.')
    return
  }
  const record = grid.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'trackForm',
    url: '/api/v1/track/form',
    focus: -1,
    record: record ?? defaults,
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
      {
        field: 'album',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/album/dropdown'),
        html: {
          label: 'Album',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'genre',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/genre/dropdown'),
        html: {
          label: 'Genre',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'media_type',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/media_type/dropdown'),
        html: {
          label: 'Media Type',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'composer',
        type: 'text',
        html: {
          label: 'Composer',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'milliseconds',
        type: 'int',
        required: true,
        html: {
          label: 'Length, ms',
          attr: 'style="width: 100%;"',
          span: 4,
        },
      },
      {
        field: 'bytes',
        type: 'int',
        html: {
          label: 'Size, bytes',
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
          w2popup.close()
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: isEditMode ? 'Edit Track' : 'New Track',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 680, height: 480, showMax: false, resizable: false,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

