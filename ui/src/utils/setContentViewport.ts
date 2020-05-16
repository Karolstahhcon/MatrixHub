const restorationIdKey = 'data-scroll-restoration-id'

const adminDataSelector = 'admin-content-scroll'
const contentDataSelector = 'content-scroll'

export const adminContentViewportSelector = `[${restorationIdKey}="${adminDataSelector}"]`
export const contentViewportSelector = `[${restorationIdKey}="${contentDataSelector}"]`
