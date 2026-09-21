import tableModel from './table_model'

export default () => ({
  _model: tableModel(),
  _observer: null,

  destroy() {
    if (this.body._observer) {
      this.body._observer.disconnect()
    }
  },

  body: {
    ['x-init']() {
      // reset table body state
      this.$nextTick(() => {
        this._model.selected = []
        this._model.count = this.$el.querySelectorAll('tr').length
      })

      // disconnect the old observer
      if (this._observer) {
        this._observer.disconnect()
      }

      // observe new rows
      this._observer = new MutationObserver(mutations => {
        mutations.forEach(mutation => {
          if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
            this._model.count = this.$el.querySelectorAll('tr').length
          }
        })
      })

      this._observer.observe(this.$el, { childList: true })
    },
  },

  footer: {
    ['x-init']() {
      this.$watch('_model', () => {
        if (this._model.selected.length > 0) {
          this.footer.text = `${this._model.selected.length} of ${this._model.count} ${makePlural('item', this._model.count)} selected`
        } else {
          this.footer.text = `Showing ${this._model.count} ${makePlural('item', this._model.count)}`
        }
      })
    },

    ['x-text']()  {
      return this.footer.text
    },
  },

  checkboxAll: {
    ['x-init']() {
      this.$watch('_model', () => {
        this.$el.checked = this._model.selected.length > 0 && this._model.selected.length === this._model.count
        this.$el.indeterminate = this._model.selected.length > 0 && this._model.selected.length !== this._model.count
      })
    },

    ['@click']() {
      const checkboxes = this.$root.querySelectorAll('input[type=checkbox][x-model="_model.selected"]')
      this._model.selected = this.$el.checked ? Array.from(checkboxes).map(x => x.value) : []
    },
  },

  checkboxFilter: {
    ['@click']() {
      if (this.$el.checked) {
        if (this.$root.dataset.field in this._model.filters) {
          this._model.filters[this.$root.dataset.field].push(this.$el.dataset.value)
        } else {
          this._model.filters[this.$root.dataset.field] = [this.$el.dataset.value]
        }
      } else {
        this._model.filters[this.$root.dataset.field] = this._model.filters[this.$root.dataset.field].filter(x => x != this.$el.dataset.value)
        if (this._model.filters[this.$root.dataset.field].length == 0) {
          delete this._model.filters[this.$root.dataset.field]
        }
      }
    },
  },
})

