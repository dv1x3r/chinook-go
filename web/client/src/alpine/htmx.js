export default () => ({
  isHtmxRequestActive: false,
  observer: null,

  init() {
    this.observer = new MutationObserver(mutations => {
      mutations.forEach(mutation => {
        if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
          this.isHtmxRequestActive = this.$el.classList.contains('htmx-request')
        }
      })
    })

    this.observer.observe(this.$el, { attributes: true })
  },

  destroy() {
    if (this.observer) {
      this.observer.disconnect()
    }
  },
})

