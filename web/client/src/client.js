import Htmx from 'htmx.org'
import Alpine from 'alpinejs'
import anchor from '@alpinejs/anchor'
import focus from '@alpinejs/focus'

import htmxRequestObserver from './alpine/htmx'
import tableComponent from './alpine/table_component'

window.htmx = Htmx

// the indicator styles are part of client.css
Htmx.config.includeIndicatorCSS = false

// htmx 4 copies the attributes of the old element onto the new one during the
// settle delay, nothing here animates a swap, so the copy is turned off
Htmx.config.defaultSettleDelay = 0

// full page reload on history navigation to prevent alpine state desync
Htmx.config.history = 'reload'

window.Alpine = Alpine
Alpine.plugin(anchor)
Alpine.plugin(focus)

Alpine.data('htmxRequestObserver', htmxRequestObserver)
Alpine.data('tableComponent', tableComponent)

Alpine.start()

// grammar utilities used by the table footer
window.makePlural = (word, count) => (count == -1 || count == 1) ? word : word + 's'
