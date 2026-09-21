import { w2layout } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

helpers.w2init()
helpers.w2initLocale({
  defaultWeekStarts: 'M',
  defaultDateFormat: 'yyyy-MM-dd',
  defaultDatetimeFormat: 'yyyy-MM-dd|hh24:mi',
  defaultTimeFormat: 'hh24:mi',
})

export const dashboardLayout = new w2layout({
  name: 'dashboardLayout',
  box: '#dashboard-layout',
  panels: [
    {
      type: 'top',
      size: 40,
      toolbar: {
        items: [
          {
            type: 'button',
            id: 'tracks',
            text: 'Tracks',
            icon: 'fa fa-music',
            onClick: async function() {
              const module = await import('./widgets/track.js')
              tabManager.OpenTab({
                id: 'tracks-tab',
                name: 'Tracks',
                closable: false,
                createWidget: () => module.createTrackLayout(),
              })
            },
          },
          {
            type: 'button',
            id: 'playlists',
            text: 'Playlists',
            icon: 'fa fa-list',
            onClick: async function() {
              const module = await import('./widgets/playlist.js')
              tabManager.OpenTab({
                id: 'playlists-tab',
                name: 'Playlists',
                createWidget: () => module.createPlaylistLayout(),
              })
            },
          },
          {
            type: 'menu',
            id: 'track-properties',
            text: 'Properties',
            icon: 'fa fa-list-check',
            items: [
              {
                type: 'html',
                id: 'track-properties-item',
                text: 'Track Properties',
                disabled: true,
              },
              {
                type: 'button',
                id: 'genres',
                text: 'Genres',
                icon: 'fa fa-tags',
                onClick: async function() {
                  const module = await import('./widgets/genre.js')
                  tabManager.OpenTab({
                    id: 'genres-tab',
                    name: 'Genres',
                    createWidget: () => module.createGenreGrid(),
                  })
                },
              },
              {
                type: 'button',
                id: 'media-types',
                text: 'Media Types',
                icon: 'fa fa-file-audio',
                onClick: async function() {
                  const module = await import('./widgets/media_type.js')
                  tabManager.OpenTab({
                    id: 'media-types-tab',
                    name: 'Media Types',
                    createWidget: () => module.createMediaTypeGrid(),
                  })
                },
              },
            ],
          },
          { type: 'break' },
          {
            type: 'button',
            id: 'customers',
            text: 'Customers',
            icon: 'fa fa-user',
            onClick: async function() {
              const module = await import('./widgets/customer.js')
              tabManager.OpenTab({
                id: 'customers-tab',
                name: 'Customers',
                createWidget: () => module.createCustomerLayout(),
              })
            },
          },
          {
            type: 'button',
            id: 'employees',
            text: 'Employees',
            icon: 'fa fa-id-badge',
            onClick: async function() {
              const module = await import('./widgets/employee.js')
              tabManager.OpenTab({
                id: 'employees-tab',
                name: 'Employees',
                createWidget: () => module.createEmployeeLayout(),
              })
            },
          },
          { type: 'break' },
          {
            type: 'button',
            id: 'invoices',
            text: 'Invoices',
            icon: 'fa fa-file-invoice',
            onClick: async function() {
              const module = await import('./widgets/invoice.js')
              tabManager.OpenTab({
                id: 'invoices-tab',
                name: 'Invoices',
                createWidget: () => module.createInvoiceLayout(),
              })
            },
          },
          { type: 'spacer' },
          {
            id: 'sql-explorer',
            type: 'button',
            icon: 'fa fa-database',
            onClick: async function() {
              const module = await import('/lib/w2ui.widgets.js')
              tabManager.OpenTab({
                id: `sql-explorer-${Date.now()}-tab`,
                name: 'SQL Explorer',
                createWidget: () => module.createSqlExplorerLayout({ url: '/api/v1/sql' }),
              })
            },
          },
          {
            id: 'dark',
            type: 'check',
            icon: 'fa fa-moon',
            checked: helpers.isDarkTheme(),
            onClick: async function(event) {
              await event.complete
              helpers.setDarkTheme(event.detail.item.checked)
            },
          },
          { type: 'break' },
          {
            type: 'button',
            id: 'storefront',
            icon: 'fa fa-store',
            text: 'Storefront',
            onClick: function() {
              window.open('/', '_blank')
            },
          },
          {
            type: 'button',
            id: 'github',
            icon: 'fa-brands fa-github',
            text: 'GitHub',
            onClick: function() {
              window.open('https://github.com/dv1x3r/chinook-go', '_blank')
            },
          },
          {
            type: 'check',
            id: 'dark',
            icon: 'fa fa-moon',
            text: 'Dark Theme',
            checked: helpers.isDarkTheme(),
            onClick: async function(event) {
              await event.complete
              helpers.setDarkTheme(event.detail.item.checked)
            },
          },
        ],
        onClick: function(event) {
          // execute onClick function for menu items
          if (event.target.includes(':') && event.detail.subItem?.onClick) {
            event.preventDefault()
            event.detail.subItem.onClick(event)
          }
        },
      },
    },
    {
      type: 'main',
      tabs: {
        onClick: event => tabManager.OnClick(event),
        onClose: event => tabManager.OnClose(event),
      },
    },
  ],
  onRender: () => {
    document.addEventListener('DOMContentLoaded', () => {
      dashboardLayout.get('top').toolbar.items[0].onClick()
    })
  },
})

export const tabManager = new helpers.TabManager(dashboardLayout.get('main').tabs)

